# calc_size web

`calc_size` 的网页版：选择本地文件或文件夹，在浏览器中计算并列出大小。

- 纯静态、零依赖：HTML + CSS + 原生 JavaScript（ES Modules），无构建步骤。
- 所有计算均在浏览器本地完成，不会上传任何数据。
- 选择文件夹时按顶层子项聚合显示（子目录显示内部所有文件之和），与命令行版语义一致。

## 本地预览

ES Modules 必须通过 HTTP 访问（直接双击 `index.html` 会被浏览器的 CORS 规则拦截），
任选其一：

```shell
python -m http.server 8000   # 打开 http://localhost:8000
# 或
npx serve .
```

## 运行测试

核心逻辑（单位换算、聚合求和）带单元测试，需要 Node ≥ 20（无需安装任何依赖）：

```shell
npm test        # 等价于 node --test
```

## 部署到 GitHub Pages

本目录是纯静态文件，无需构建：在仓库 Settings → Pages 中将其指向 `web/` 目录
（或用 Actions 上传该目录作为 artifact）即可。

## 与命令行版的差异

- 浏览器枚举不到的文件（无权限、系统隐藏文件）与空目录不会出现在结果中，
  数字可能因此略有差异；网页版不存在"-1"错误值，读不到的直接缺失。
- 单位换算与 Go 版一致（1024 进制、严格大于 1024 才进位、两位小数），
  但补上了 Go 版缺失的 TB 档（B → KB → MB → GB → TB → PB）。
