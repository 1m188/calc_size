#!/usr/bin/env bash
# 用法：bash build-all.sh

set -euo pipefail

# ====== 临时环境变量（仅本次终端）======
export CGO_ENABLED=0                 # 纯静态
export GOOS="${GOOS:-}"              # 允许外部传入，下同
export GOARCH="${GOARCH:-}"

# ====== 要构建的平台列表 ======
platforms=(
  "windows/amd64"
  "windows/386"
  "windows/arm64"
  "linux/amd64"
  "linux/386"
  "linux/arm64"
  "linux/arm"
  "linux/mips64le"
  "linux/mipsle"
  "linux/s390x"
  "linux/ppc64le"
  "darwin/amd64"
  "darwin/arm64"
  "freebsd/amd64"
  "freebsd/arm64"
  "openbsd/amd64"
  "netbsd/amd64"
  "android/arm64"
  # "android/arm"            # require cgo
  # "android/amd64"          # require cgo 
  # "ios/arm64"              # require cgo
  # "ios/amd64"              # require cgo
  "js/wasm"
  "aix/ppc64"
  "solaris/amd64"
)

# ====== 输出目录 ======
mkdir -p dist

# ====== 循环交叉编译 ======
for pl in "${platforms[@]}"; do
  IFS='/' read -r os arch <<< "$pl"
  export GOOS="$os"
  export GOARCH="$arch"

  out="dist/calc_${os}_${arch}"
  [ "$os" = "windows" ] && out+='.exe'

  echo "Building $out ..."
  go build -ldflags="-s -w" -o "$out" .
done

echo "✅ 全部完成，产物在 dist/ 目录"