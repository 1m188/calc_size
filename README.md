# calc_size

calculate the file or folder size

## Usage

```shell
calc_size(.exe) --help
```

## Example

### Normal
```shell
calc_size(.exe) -p a.txt,.,..,D:/Code/,/path/not/exist
```

```shell
a.txt            :   1.00  B
.                :  12.00 MB
..               :  10.00 GB
D:/Code/         :  12.00 GB
/path/not/exist  :  -1.00  B
```

### JSON
```shell
calc_size(.exe) -p a.txt,.,..,D:/Code/,/path/not/exist -j
```

```json
[{"path":"a.txt","size":"1"},{"path":".","size":"12582912"},{"path":"..","size":"10737418240"},{"path":"D:/Code/","size":"12884901888"},{"path":"/path/not/exist","size":"-1"}]
```

### CSV
```shell
calc_size(.exe) -p a.txt,.,..,D:/Code/,/path/not/exist -c
```

```csv
a.txt,1
.,12582912
..,10737418240
D:/Code/,12884901888
/path/not/exist,-1
```

### 设置递归最大深度
```shell
calc_size(.exe) -p D:/Code/ -d 2
```

## TODO

- [ √ ] 用 pflag/cobra 等替代手写 os.Args，支持长选项、子命令、自动 help 等功能。
- [ √ ] 默认 human-readable（1.2 G、345 K）；加 `-j/--json` 输出 JSON，`-c/--csv` 输出 CSV，方便被脚本调用。
- [ √ ] 实时显示已处理文件/字节数。
- [ √ ] 单元测试：造几个固定大小的文件/目录，断言总大小；以及一些边界测试。
- [ √ ] 基准测试：对比单线程vs并发，不同文件规模的性能差异。
- [ √ ] 加 .github/workflows 做 CI 自动发布 Release。
- [ ] 排除/包含规则：`-e node_modules -e *.log` 支持通配符/正则忽略；`-i "*.go"` 只统计 Go 源文件等。
- [ √ ] 深度限制：`-d 3` 只统计 3 层子目录，快速看顶层分布。
- [ ] 大文件 Top-N：`-t 20` 列出占用最大的 20 个文件（类似 ncdu）。
- [ ] 重复文件检测：扫完大小后，对相同大小的文件计算哈希，输出重复列表（可删除）。
- [ ] 可视化：终端直方图，打印大小分布条形图；或生成 `tree-map` HTML，浏览器里拖拽看块大小。
- [ ] 压缩比预估：边扫边用 gzip/zstd 压缩采样 1 MB，估算目录“可压缩空间”，提示用户是否值得打包。
- [ ] 软链/硬链特殊处理：默认不跟随软链；加 `-L` 跟随；加 `--count-hard-links` 把硬链只算一次或多次。
- [ ] 权限报表：扫完输出“无读权限文件列表”，方便管理员审计。
