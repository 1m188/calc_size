import test from 'node:test';
import assert from 'node:assert/strict';

import { fileListToEntries, aggregateEntries } from '../js/scan.js';

test('fileListToEntries：无 webkitRelativePath 时使用文件名', () => {
  const entries = fileListToEntries([{ name: 'a.txt', size: 1 }]);
  assert.deepEqual(entries, [{ path: 'a.txt', size: 1 }]);
});

test('fileListToEntries：文件夹选择使用 webkitRelativePath', () => {
  const entries = fileListToEntries([
    { name: 'f.txt', size: 2, webkitRelativePath: 'dir/f.txt' },
  ]);
  assert.deepEqual(entries, [{ path: 'dir/f.txt', size: 2 }]);
});

test('聚合：空输入返回空行与 0 总计', () => {
  assert.deepEqual(aggregateEntries([]), { rows: [], total: 0 });
});

test('聚合：多文件选择逐行显示并求总计', () => {
  const result = aggregateEntries([
    { path: 'a.txt', size: 1 },
    { path: 'b.bin', size: 5 },
  ]);
  assert.deepEqual(result, {
    rows: [
      { path: 'a.txt', size: 1 },
      { path: 'b.bin', size: 5 },
    ],
    total: 6,
  });
});

test('聚合：文件夹选择按顶层子项聚合，显示时去掉根前缀', () => {
  const result = aggregateEntries([
    { path: 'root/sub1/f1.txt', size: 100 },
    { path: 'root/sub1/nested/f2.bin', size: 50 },
    { path: 'root/sub2/f3.log', size: 30 },
    { path: 'root/file.txt', size: 20 },
  ]);
  assert.deepEqual(result, {
    rows: [
      { path: 'sub1', size: 150 },
      { path: 'sub2', size: 30 },
      { path: 'file.txt', size: 20 },
    ],
    total: 200,
  });
});

test('聚合：顶层子项按首次出现顺序排列', () => {
  const result = aggregateEntries([
    { path: 'root/z/1', size: 1 },
    { path: 'root/a/2', size: 2 },
  ]);
  assert.deepEqual(result.rows.map((r) => r.path), ['z', 'a']);
});

test('聚合：根目录下的散文件单独成行', () => {
  const result = aggregateEntries([{ path: 'root/f.txt', size: 7 }]);
  assert.deepEqual(result, { rows: [{ path: 'f.txt', size: 7 }], total: 7 });
});

test('聚合：根段不一致时退化为按完整路径逐行', () => {
  const result = aggregateEntries([
    { path: 'a/x', size: 1 },
    { path: 'b/y', size: 2 },
  ]);
  assert.deepEqual(result, {
    rows: [
      { path: 'a/x', size: 1 },
      { path: 'b/y', size: 2 },
    ],
    total: 3,
  });
});
