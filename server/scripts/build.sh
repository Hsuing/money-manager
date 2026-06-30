#!/bin/bash

# 解决 Git Bash 环境下 Windows 路径和 Linux 路径的转换问题
export MSYS_NO_PATHCONV=1

# 获取绝对路径（优先尝试 pwd -W 获取 Windows 格式路径，以防挂载失败）
CURRENT_DIR=$(cd "$(dirname "$0")/.." && pwd -W 2>/dev/null || pwd)

echo "开始在 Docker 中编译服务端代码..."
echo "挂载目录: ${CURRENT_DIR}"

docker run --rm \
  -v "${CURRENT_DIR}:/app" \
  -w /app \
  golang:1.25 \
  sh -c "go env -w GOPROXY=https://goproxy.cn,direct && go build -buildvcs=false -o bin/server ./cmd/server"

echo "编译完成！生成的可执行文件在 bin/server"
