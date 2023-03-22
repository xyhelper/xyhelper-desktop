#!/bin/bash
set -e

# 构建wails项目
rm -rf build/bin
wails build -platform=windows/amd64 -webview2=embed -o xyhelper-win64.exe
wails build -platform=darwin/amd64,darwin/arm64 -webview2=embed  -o xyhelper-mac
