import test from 'node:test';
import assert from 'node:assert/strict';

import { formatSize, formatSizeText } from '../js/format.js';

test('0 显示为 0.00 B', () => {
  assert.deepEqual(formatSize(0), { value: '0.00', unit: 'B' });
});

test('1 字节显示为 1.00 B（README 示例）', () => {
  assert.deepEqual(formatSize(1), { value: '1.00', unit: 'B' });
});

test('1023 字节仍是 B', () => {
  assert.deepEqual(formatSize(1023), { value: '1023.00', unit: 'B' });
});

test('恰好 1024 字节不进位（与 Go 版“严格大于”语义一致）', () => {
  assert.deepEqual(formatSize(1024), { value: '1024.00', unit: 'B' });
});

test('1025 字节进位为 1.00 KB', () => {
  assert.deepEqual(formatSize(1025), { value: '1.00', unit: 'KB' });
});

test('12 MiB 显示为 12.00 MB（README 示例）', () => {
  assert.deepEqual(formatSize(12582912), { value: '12.00', unit: 'MB' });
});

test('10 GiB 显示为 10.00 GB（README 示例）', () => {
  assert.deepEqual(formatSize(10737418240), { value: '10.00', unit: 'GB' });
});

test('恰好 1 TiB 显示为 1024.00 GB', () => {
  assert.deepEqual(formatSize(2 ** 40), { value: '1024.00', unit: 'GB' });
});

test('1 TiB 多一点进位为 1.00 TB', () => {
  assert.deepEqual(formatSize(2 ** 40 + 2 ** 30), { value: '1.00', unit: 'TB' });
});

test('恰好 1 PiB 显示为 1024.00 TB', () => {
  assert.deepEqual(formatSize(2 ** 50), { value: '1024.00', unit: 'TB' });
});

test('1 PiB 多一点进位为 1.00 PB', () => {
  assert.deepEqual(formatSize(2 ** 50 + 2 ** 40), { value: '1.00', unit: 'PB' });
});

test('PB 是最大档位：2^53-1 显示为 8.00 PB', () => {
  assert.deepEqual(formatSize(Number.MAX_SAFE_INTEGER), { value: '8.00', unit: 'PB' });
});

test('负数走防御分支', () => {
  assert.deepEqual(formatSize(-1), { value: '-', unit: '' });
});

test('非有限值走防御分支', () => {
  assert.deepEqual(formatSize(Infinity), { value: '-', unit: '' });
  assert.deepEqual(formatSize(NaN), { value: '-', unit: '' });
});

test('formatSizeText 拼接数字与单位', () => {
  assert.equal(formatSizeText(1), '1.00 B');
  assert.equal(formatSizeText(12582912), '12.00 MB');
});

test('formatSizeText 对防御分支只显示占位符', () => {
  assert.equal(formatSizeText(-1), '-');
});
