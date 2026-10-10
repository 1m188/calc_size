// 遍历与建树，对应 Go 版 get_file_size 包的职责。
// 浏览器的 File API 已直接给出每个文件的字节数（File.size，来自元数据，无需读文件内容），
// 这里只做三件事：FileList 归一化、按用户所选顶层项构建可展开的树、提供排序比较器。

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
 * 把条目构建为一组根节点（用户所选的每个顶层项是一棵子树）。
 * 条目路径相对于"用户所选"：文件夹选择含根段，文件选择为文件名本身。
 * 同一次调用中出现相同顶层名时合并进同一棵树；跨调用由调用方决定覆盖策略。
 *
 * 节点结构：{ name, path, size, isDir, children }
 *   - name:    当前层级的显示名
 *   - path:    从顶层项算起的相对路径（顶层项 path === name）
 *   - size:    文件为自身大小；目录为所有后代文件之和
 *   - children 保持首次出现顺序（排序是视图层的事，见 compareNodes）
 *
 * @param {Array<{ path: string, size: number }>} entries
 * @returns {Array<{ name: string, path: string, size: number, isDir: boolean, children: Array }>}
 */
export function entriesToTree(entries) {
  const tops = new Map(); // 顶层名 -> { root, dirs: Map<相对路径, 目录节点> }
  for (const entry of entries) {
    const segments = entry.path.split('/');
    const topName = segments[0];
    let top = tops.get(topName);
    if (!top) {
      const root = {
        name: topName,
        path: topName,
        size: segments.length === 1 ? entry.size : 0,
        isDir: segments.length > 1,
        children: [],
      };
      top = { root, dirs: new Map([['', root]]) };
      tops.set(topName, top);
    }

    let parentPath = '';
    for (let i = 1; i < segments.length; i += 1) {
      const name = segments[i];
      const path = parentPath ? `${parentPath}/${name}` : name;
      const isLeaf = i === segments.length - 1;
      let node = top.dirs.get(path);
      if (!node) {
        node = {
          name,
          path: `${topName}/${path}`,
          size: isLeaf ? entry.size : 0,
          isDir: !isLeaf,
          children: [],
        };
        top.dirs.get(parentPath).children.push(node);
        if (!isLeaf) {
          top.dirs.set(path, node);
        }
      }
      parentPath = path;
    }
  }

  const roots = [...tops.values()].map((t) => t.root);
  for (const root of roots) {
    computeSize(root);
  }
  return roots;
}

// 自底向上聚合目录大小（文件节点已在建树时写入 size）
function computeSize(node) {
  if (!node.isDir) {
    return node.size;
  }
  let sum = 0;
  for (const child of node.children) {
    sum += computeSize(child);
  }
  node.size = sum;
  return sum;
}

/**
 * 节点排序比较器（资源管理器习惯）：
 * 目录始终排在文件前面；再按 key 与 dir 排序；
 * 大小相同时按名称升序兜底（不随方向翻转）。
 * @param {'name'|'size'} key
 * @param {'asc'|'desc'} dir
 */
export function compareNodes(key, dir) {
  const sign = dir === 'desc' ? -1 : 1;
  const byName = (a, b) => a.name.localeCompare(b.name, 'zh-Hans-CN', { numeric: true });
  return (a, b) => {
    if (a.isDir !== b.isDir) {
      return a.isDir ? -1 : 1;
    }
    if (key === 'size') {
      return (a.size - b.size) * sign || byName(a, b);
    }
    return byName(a, b) * sign;
  };
}
