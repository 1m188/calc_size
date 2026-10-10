import test from 'node:test';
import assert from 'node:assert/strict';

import { fileListToEntries, entriesToTree, compareNodes } from '../js/scan.js';

// —— fileListToEntries（行为保持不变）——

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

// —— entriesToTree：把条目构建为可展开的树 ——

test('建树：文件选择产生叶子根节点', () => {
  const roots = entriesToTree([
    { path: 'a.txt', size: 1 },
    { path: 'b.bin', size: 5 },
  ]);
  assert.deepEqual(roots, [
    { name: 'a.txt', path: 'a.txt', size: 1, isDir: false, children: [] },
    { name: 'b.bin', path: 'b.bin', size: 5, isDir: false, children: [] },
  ]);
});

test('建树：文件夹选择构建嵌套树并聚合目录大小', () => {
  const roots = entriesToTree([
    { path: 'root/sub1/f1.txt', size: 100 },
    { path: 'root/sub1/nested/f2.bin', size: 50 },
    { path: 'root/sub2/f3.log', size: 30 },
    { path: 'root/file.txt', size: 20 },
  ]);
  assert.deepEqual(roots, [
    {
      name: 'root',
      path: 'root',
      size: 200,
      isDir: true,
      children: [
        {
          name: 'sub1',
          path: 'root/sub1',
          size: 150,
          isDir: true,
          children: [
            { name: 'f1.txt', path: 'root/sub1/f1.txt', size: 100, isDir: false, children: [] },
            {
              name: 'nested',
              path: 'root/sub1/nested',
              size: 50,
              isDir: true,
              children: [
                { name: 'f2.bin', path: 'root/sub1/nested/f2.bin', size: 50, isDir: false, children: [] },
              ],
            },
          ],
        },
        {
          name: 'sub2',
          path: 'root/sub2',
          size: 30,
          isDir: true,
          children: [
            { name: 'f3.log', path: 'root/sub2/f3.log', size: 30, isDir: false, children: [] },
          ],
        },
        { name: 'file.txt', path: 'root/file.txt', size: 20, isDir: false, children: [] },
      ],
    },
  ]);
});

test('建树：深层嵌套逐级聚合', () => {
  const roots = entriesToTree([{ path: 'root/a/b/c.txt', size: 7 }]);
  assert.deepEqual(roots, [
    {
      name: 'root',
      path: 'root',
      size: 7,
      isDir: true,
      children: [
        {
          name: 'a',
          path: 'root/a',
          size: 7,
          isDir: true,
          children: [
            {
              name: 'b',
              path: 'root/a/b',
              size: 7,
              isDir: true,
              children: [
                { name: 'c.txt', path: 'root/a/b/c.txt', size: 7, isDir: false, children: [] },
              ],
            },
          ],
        },
      ],
    },
  ]);
});

test('建树：同一次调用中同名顶层合并为同一棵树', () => {
  const roots = entriesToTree([
    { path: 'root/a.txt', size: 10 },
    { path: 'root/b.txt', size: 20 },
  ]);
  assert.equal(roots.length, 1);
  assert.equal(roots[0].size, 30);
  assert.deepEqual(roots[0].children.map((c) => c.name), ['a.txt', 'b.txt']);
});

test('建树：空输入返回空数组', () => {
  assert.deepEqual(entriesToTree([]), []);
});

// —— compareNodes：资源管理器式排序 ——

const file = (name, size) => ({ name, path: name, size, isDir: false, children: [] });
const dir = (name, size, children = []) => ({ name, path: name, size, isDir: true, children });

test('排序：目录永远排在文件前面（即使按大小降序）', () => {
  const nodes = [file('big.txt', 100), dir('small-dir', 1)];
  const sorted = [...nodes].sort(compareNodes('size', 'desc'));
  assert.deepEqual(sorted.map((n) => n.name), ['small-dir', 'big.txt']);
});

test('排序：按大小升序/降序', () => {
  const nodes = [file('a', 5), file('b', 2), file('c', 8)];
  assert.deepEqual([...nodes].sort(compareNodes('size', 'asc')).map((n) => n.size), [2, 5, 8]);
  assert.deepEqual([...nodes].sort(compareNodes('size', 'desc')).map((n) => n.size), [8, 5, 2]);
});

test('排序：名称比较数字感知（f2 排在 f10 前）', () => {
  const nodes = [file('f10', 1), file('f2', 1)];
  assert.deepEqual([...nodes].sort(compareNodes('name', 'asc')).map((n) => n.name), ['f2', 'f10']);
});

test('排序：大小相同时按名称升序兜底（不随方向翻转）', () => {
  const nodes = [file('b', 5), file('a', 5)];
  assert.deepEqual([...nodes].sort(compareNodes('size', 'desc')).map((n) => n.name), ['a', 'b']);
});
