// 单位换算，语义与 Go 版 calc_size 的 data/file_size.go 对齐：
// 1024 进制、严格大于 1024 才进位、保留两位小数；单位序列补上了 Go 版缺失的 TB 档。

const UNITS = ['B', 'KB', 'MB', 'GB', 'TB', 'PB'];

/**
 * 把字节数换算为人类可读的 { value, unit }。
 * 非有限值或负数返回 { value: '-', unit: '' }
 * （网页版不会出现 -1，此分支纯属防御）。
 * @param {number} bytes
 * @returns {{ value: string, unit: string }}
 */
export function formatSize(bytes) {
  if (!Number.isFinite(bytes) || bytes < 0) {
    return { value: '-', unit: '' };
  }
  let size = bytes;
  let unit = 0;
  while (size > 1024 && unit < UNITS.length - 1) {
    size /= 1024;
    unit += 1;
  }
  return { value: size.toFixed(2), unit: UNITS[unit] };
}

/**
 * 直接返回 "12.00 GB" 形式的文本；防御分支只显示 "-"。
 * @param {number} bytes
 */
export function formatSizeText(bytes) {
  const { value, unit } = formatSize(bytes);
  return unit ? `${value} ${unit}` : value;
}
