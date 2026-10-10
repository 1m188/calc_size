import { formatSizeText } from './format.js';
import { fileListToEntries, entriesToTree, compareNodes } from './scan.js';

const mainEl = document.querySelector('main');
const pickMenu = document.getElementById('pick-menu');
const fileInput = document.getElementById('pick-files');
const folderInput = document.getElementById('pick-folder');
const clearBtn = document.getElementById('clear');
const statusEl = document.getElementById('status');
const resultEl = document.getElementById('result');
const rowsEl = document.getElementById('rows');
const thName = document.getElementById('th-name');
const thSize = document.getElementById('th-size');
const totalLabelEl = document.getElementById('total-label');
const totalSizeEl = document.getElementById('total-size');
const dropOverlay = document.getElementById('drop-overlay');

// ---- 状态 ----
const selections = new Map(); // 顶层名 -> 根节点；后选择的同名项覆盖旧结果
const expanded = new Set(); // 已展开目录的 path
let sortKey = 'name';
let sortDir = 'asc'; // 大小列首次点击为降序（找大文件更实用）
let lastToggled = null; // { path, depth }：最近一次展开/折叠的目录，用于子树入场动画

// ---- 选择入口：菜单 ----
fileInput.addEventListener('change', () => {
  addEntries(fileListToEntries(fileInput.files));
  fileInput.value = ''; // 清空以允许再次选择同一路径时仍触发 change
  pickMenu.open = false;
});

folderInput.addEventListener('change', () => {
  addEntries(fileListToEntries(folderInput.files));
  folderInput.value = '';
  pickMenu.open = false;
});

clearBtn.addEventListener('click', () => {
  selections.clear();
  expanded.clear();
  render(false);
  statusEl.textContent = '已清空。';
});

// ---- 选择入口：拖拽（唯一支持"文件+文件夹混合多选"的浏览器机制）----
mainEl.addEventListener('dragenter', onDragOver);
mainEl.addEventListener('dragover', onDragOver);
mainEl.addEventListener('dragleave', (e) => {
  if (!mainEl.contains(e.relatedTarget)) {
    dropOverlay.classList.remove('show');
  }
});
mainEl.addEventListener('drop', async (e) => {
  e.preventDefault();
  dropOverlay.classList.remove('show');

  const entries = [];
  const walkers = [];
  for (const item of e.dataTransfer?.items ?? []) {
    const entry = item.webkitGetAsEntry?.();
    if (entry) {
      walkers.push(walkEntry(entry, entries));
    } else {
      const file = item.getAsFile?.();
      if (file) {
        entries.push({ path: file.name, size: file.size });
      }
    }
  }
  await Promise.all(walkers);
  addEntries(entries);
});

function onDragOver(e) {
  e.preventDefault(); // 允许 drop
  dropOverlay.classList.add('show');
}

// 递归枚举拖入的文件/目录；readEntries 每次最多返回 100 条，必须循环读到空数组为止
async function walkEntry(entry, out, prefix = '') {
  if (entry.isFile) {
    const file = await new Promise((resolve, reject) => entry.file(resolve, reject));
    out.push({ path: prefix ? `${prefix}/${entry.name}` : entry.name, size: file.size });
  } else if (entry.isDirectory) {
    const dirPath = prefix ? `${prefix}/${entry.name}` : entry.name;
    const reader = entry.createReader();
    for (;;) {
      const batch = await new Promise((resolve, reject) => reader.readEntries(resolve, reject));
      if (batch.length === 0) {
        break;
      }
      for (const child of batch) {
        await walkEntry(child, out, dirPath);
      }
    }
  }
}

// ---- 数据合并 ----
function addEntries(entries) {
  if (entries.length === 0) {
    statusEl.textContent = '没有可统计的文件：所选项为空，或浏览器未能枚举到任何文件。';
    return;
  }
  statusEl.textContent = '';
  for (const root of entriesToTree(entries)) {
    selections.set(root.name, root); // 同名顶层项：后选择的覆盖
    if (root.isDir) {
      expanded.add(root.path); // 顶层目录默认展开一层
    }
  }
  render(true);
  flashTotal();
}

// ---- 渲染 ----
// animateAll：新数据到达时让可见行整体级联入场；
// 否则只有 lastToggled 的子树行入场（展开/折叠的局部动画）。
function render(animateAll = false) {
  const roots = [...selections.values()];

  // 先收集当前可见的节点（含排序），再统一渲染
  const visible = [];
  const collect = (nodes, depth) => {
    for (const node of [...nodes].sort(compareNodes(sortKey, sortDir))) {
      visible.push({ node, depth });
      if (node.isDir && node.children.length > 0 && expanded.has(node.path)) {
        collect(node.children, depth + 1);
      }
    }
  };
  collect(roots, 0);

  // 比例条基准：当前可见节点中的最大值
  const maxNode = visible.reduce((max, v) => Math.max(max, v.node.size), 0);

  rowsEl.replaceChildren();
  const fragment = document.createDocumentFragment();
  let animatedCount = 0; // 级联动画上限，超过后直接显示，避免大列表卡顿
  for (const { node, depth } of visible) {
    const inSubtree = lastToggled !== null && node.path.startsWith(`${lastToggled.path}/`);
    let anim = null;
    if ((animateAll || inSubtree) && animatedCount < 40) {
      anim = { delay: animatedCount * 18 };
      animatedCount += 1;
    }
    fragment.append(buildRow(node, depth, anim, maxNode));
  }
  rowsEl.append(fragment);

  let total = 0;
  let fileCount = 0;
  for (const root of roots) {
    total += root.size;
    fileCount += countFiles(root);
  }
  totalLabelEl.textContent = `总计（${fileCount} 个文件）`;
  totalSizeEl.textContent = formatSizeText(total);

  resultEl.hidden = roots.length === 0;
  clearBtn.hidden = roots.length === 0;
  syncSortHeader(thName, 'name');
  syncSortHeader(thSize, 'size');
  lastToggled = null; // 本次渲染已消费
}

// 只渲染可见行：折叠的目录不渲染其子树
function buildRow(node, depth, anim, maxNode) {
  const tr = document.createElement('tr');
  if (anim) {
    tr.classList.add('enter');
    tr.style.setProperty('--d', `${anim.delay}ms`);
  }
  if (lastToggled !== null && node.path === lastToggled.path) {
    tr.classList.add('just-toggled'); // 箭头旋转补间
  }
  if (maxNode > 0 && node.size > 0) {
    tr.style.setProperty('--w', `${Math.min(100, (node.size / maxNode) * 100)}%`);
  }

  const pathTd = document.createElement('td');
  pathTd.style.paddingLeft = `${0.5 + depth * 1.1}rem`;
  pathTd.title = node.path;

  if (node.isDir && node.children.length > 0) {
    const isOpen = expanded.has(node.path);
    const toggle = document.createElement('button');
    toggle.type = 'button';
    toggle.className = isOpen ? 'toggle open' : 'toggle';
    toggle.textContent = '▸';
    toggle.setAttribute('aria-expanded', String(isOpen));
    toggle.addEventListener('click', () => toggleNode(node, depth));
    const name = document.createElement('span');
    name.className = 'dir-name';
    name.textContent = node.name;
    name.addEventListener('click', () => toggleNode(node, depth));
    pathTd.append(toggle, name);
  } else {
    pathTd.textContent = node.name;
  }

  const sizeTd = document.createElement('td');
  sizeTd.textContent = formatSizeText(node.size);
  tr.append(pathTd, sizeTd);
  return tr;
}

function toggleNode(node, depth) {
  if (expanded.has(node.path)) {
    expanded.delete(node.path);
  } else {
    expanded.add(node.path);
  }
  lastToggled = { path: node.path, depth };
  render(false);
}

function countFiles(node) {
  if (!node.isDir) {
    return 1;
  }
  let count = 0;
  for (const child of node.children) {
    count += countFiles(child);
  }
  return count;
}

function flashTotal() {
  totalSizeEl.classList.remove('flash');
  void totalSizeEl.offsetWidth; // 强制回流以重启动画
  totalSizeEl.classList.add('flash');
}

// ---- 排序 ----
for (const [buttonId, key, firstDir] of [
  ['sort-name', 'name', 'asc'],
  ['sort-size', 'size', 'desc'],
]) {
  document.getElementById(buttonId).addEventListener('click', () => {
    if (sortKey === key) {
      sortDir = sortDir === 'asc' ? 'desc' : 'asc';
    } else {
      sortKey = key;
      sortDir = firstDir;
    }
    render(false);
  });
}

function syncSortHeader(th, key) {
  const active = sortKey === key;
  const value = active ? (sortDir === 'asc' ? 'ascending' : 'descending') : 'none';
  th.setAttribute('aria-sort', value);
}

render(); // 初始：隐藏表格、同步表头状态
