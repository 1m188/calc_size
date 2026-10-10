import { formatSizeText } from './format.js';
import { fileListToEntries, aggregateEntries } from './scan.js';

const fileInput = document.getElementById('pick-files');
const folderInput = document.getElementById('pick-folder');
const statusEl = document.getElementById('status');
const resultEl = document.getElementById('result');
const rowsEl = document.getElementById('rows');
const totalLabelEl = document.getElementById('total-label');
const totalSizeEl = document.getElementById('total-size');

function handlePick(input) {
  const entries = fileListToEntries(input.files);
  input.value = ''; // 清空以允许再次选择同一路径时仍触发 change 事件

  if (entries.length === 0) {
    resultEl.hidden = true;
    statusEl.textContent = '没有可统计的文件：文件夹为空，或浏览器未能枚举到任何文件。';
    return;
  }
  statusEl.textContent = '';

  const { rows, total } = aggregateEntries(entries);

  rowsEl.replaceChildren();
  const fragment = document.createDocumentFragment();
  for (const row of rows) {
    const tr = document.createElement('tr');
    const path = document.createElement('td');
    path.textContent = row.path;
    const size = document.createElement('td');
    size.textContent = formatSizeText(row.size);
    tr.append(path, size);
    fragment.append(tr);
  }
  rowsEl.append(fragment);

  totalLabelEl.textContent = `总计（${entries.length} 个文件）`;
  totalSizeEl.textContent = formatSizeText(total);
  resultEl.hidden = false;
}

fileInput.addEventListener('change', () => handlePick(fileInput));
folderInput.addEventListener('change', () => handlePick(folderInput));
