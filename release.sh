#!/bin/bash
# Useage: ./release.sh version
set -e

# 生成版本号
if [ -z "$1" ]; then
    echo "请输入版本号"
    exit 1
fi
version=$1


# 构建wails项目
rm -rf build/bin
# 构建windows
wails build -platform=windows/amd64 -webview2=embed -o xyhelper-windows-amd64.exe -ldflags "-X 'main.Version=$version'" -obfuscated -garbleargs "-literals -tiny -seed=myrandomseed"
# 构建 linux amd64
# wails build --platform=linux/amd64 -webview2=embed -ldflags "-X 'main.Version=$version'" -o xyhelper-linux-amd64

# 构建 darwin amd64
wails build -platform=darwin/amd64 -webview2=embed -ldflags "-X 'main.Version=$version'"
mv build/bin/xyhelper.app build/bin/xyhelper-darwin-amd64.app
# 给应用签名
# codesign --deep --force --verbose --sign "Developer ID Application: xyhelper" xyhelper-darwin-amd64.app
# 构建 darwin arm64
wails build -platform=darwin/arm64 -webview2=embed -ldflags "-X 'main.Version=$version'"
mv build/bin/xyhelper.app build/bin/xyhelper-darwin-arm64.app
# 给应用签名
# codesign --deep --force --verbose --sign "Developer ID Application: xyhelper" xyhelper-darwin-arm64.app
# 压缩为zip
cd build/bin
zip xyhelper-windows-amd64.zip xyhelper-windows-amd64.exe
zip -r xyhelper-darwin-arm64.zip xyhelper-darwin-arm64.app
zip -r xyhelper-darwin-amd64.zip xyhelper-darwin-amd64.app



