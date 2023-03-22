#!/bin/bash
set -e

# 构建wails项目
wails build -platform=windows/amd64,darwin/amd64,darwin/arm64 -clean -webview2=embed