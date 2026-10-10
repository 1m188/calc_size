// 遍历与聚合，对应 Go 版 get_file_size 包的职责。
// 浏览器的 File API 已直接给出每个文件的字节数（File.size，来自元数据，无需读文件内容），
// 这里只做两件事：把 FileList 归一成 { path, size } 列表；按顶层子项聚合求和。

/**
 * 把 FileList（或任何可迭代的 { name, size, webkitRelativePath? }）归一为条目列表。
 * 选择文件夹时路径为 webkitRelativePath（含根段）；选择文件时退化为文件名。
 * @param {Iterable<{ name: string, size: number, webkitRelativePath?: string }>} fileList
 * @returns {Array<{ path: string, size: number }>}
 */
export function fileListToEntries(fileList) {
  const entries = [];
  for (const file of fileList) {
    entries.push({ path: file.webkitRelativePath || file.name, size: file.size });
  }
  return entries;
}

/**
 * 聚合条目：
 * - 文件选择：每个文件一行；
 * - 文件夹选择：所有条目共享同一个根段（webkitdirectory 的行为），
 *   按顶层子项聚合，行路径去掉根前缀；
 * - 根段不一致（异常情况）时退化为按完整路径逐行。
 * 返回 { rows: [{ path, size }], total }，rows 保持首次出现顺序。
 * @param {Array<{ path: string, size: number }>} entries
 * @returns {{ rows: Array<{ path: string, size: number }>, total: number }}
 */
export function aggregateEntries(entries) {
  if (entries.length === 0) {
    return { rows: [], total: 0 };
  }

  const hasFolderPaths = entries.every((e) => e.path.includes('/'));
  const root = hasFolderPaths ? entries[0].path.split('/')[0] : null;
  const isFolderPick =
    hasFolderPaths && entries.every((e) => e.path.split('/')[0] === root);

  if (!isFolderPick) {
    return sumBy(entries, (e) => e.path);
  }

  const { rows, total } = sumBy(entries, (e) => {
    const segments = e.path.split('/');
    return segments.slice(0, 2).join('/');
  });
  const prefix = `${root}/`;
  return {
    rows: rows.map((row) => ({
      path: row.path.startsWith(prefix) ? row.path.slice(prefix.length) : row.path,
      size: row.size,
    })),
    total,
  };
}

/**
 * 按 keyOf 分组求和，保持首次出现顺序。
 */
function sumBy(entries, keyOf) {
  const sizes = new Map();
  let total = 0;
  for (const entry of entries) {
    const key = keyOf(entry);
    sizes.set(key, (sizes.get(key) ?? 0) + entry.size);
    total += entry.size;
  }
  return { rows: [...sizes].map(([path, size]) => ({ path, size })), total };
}
