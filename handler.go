package main

import (
	"io"
	"net/http"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

// AssetHandler is a custom asset handler
type AssetHandler struct {
	http.Handler
}

// NewAssetHandler creates a new asset handler
func NewAssetHandler() *AssetHandler {
	return &AssetHandler{}
}

// ServeHTTP serves the api
func (a *AssetHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	// 获取请求的路径
	path := req.URL.Path
	// 获取请求的方法
	method := req.Method
	// 打印日志
	g.Log().Infof(ctx, "path: %s, method: %s", path, method)
	mux := http.NewServeMux()
	// 注册会话处理函数
	mux.HandleFunc("/api/session", ApiSession)
	// 注册配置处理函数
	mux.HandleFunc("/api/config", ApiConfig)
	// 注册上传图片处理函数
	mux.HandleFunc("/api/upload-image", UploadImage)
	// 分发请求到不同的处理函数
	mux.ServeHTTP(res, req)
}

// ApiSession 处理会话
func ApiSession(res http.ResponseWriter, req *http.Request) {
	// 获取请求的方法
	method := req.Method
	// 如果不是POST请求，则返回错误
	if method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// 构造要返回的 JSON 数据
	responseData := map[string]interface{}{
		"status":  "Success",
		"message": "",
		"data": map[string]interface{}{
			"auth":  false,
			"model": "ChatGPTUnofficialProxyAPI",
		},
	}
	responseJson := gjson.New(responseData)
	// 将 JSON 数据写入响应
	res.Header().Set("Content-Type", "application/json")
	res.Write(responseJson.MustToJson())
}

// ApiConfig 配置
func ApiConfig(res http.ResponseWriter, req *http.Request) {
	// 获取请求的方法
	// method := req.Method
	// 如果不是POST请求，则返回错误
	// if method != http.MethodPost {
	// 	res.WriteHeader(http.StatusMethodNotAllowed)
	// 	return
	// }
	// 构造要返回的 JSON 数据
	responseData := map[string]interface{}{
		"status":  "Success",
		"message": "",

		"data": map[string]interface{}{
			"reverseProxy": BaseURI + "/backend-api/conversation",
			"accessToken":  AccessToken,
			"version":      Version,
		},
	}
	responseJson := gjson.New(responseData)
	// 将 JSON 数据写入响应
	res.Header().Set("Content-Type", "application/json")
	res.Write(responseJson.MustToJson())
}

// UploadImage 上传图片
func UploadImage(res http.ResponseWriter, req *http.Request) {
	g.Log().Infof(req.Context(), "UploadImage~~~~~~~~")

	// 获取请求的方法
	method := req.Method
	// 如果不是POST请求，则返回错误
	if method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// 解析上传的附件
	file, header, err := req.FormFile("image")
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	// 打开保存文件对话框
	filepath, err := runtime.SaveFileDialog(app.ctx, runtime.SaveDialogOptions{
		DefaultFilename: header.Filename,
		Title:           "保存图片",
	})
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	g.Log().Infof(req.Context(), "filepath: %s", filepath)
	// 保存文件
	if filepath != "" {
		out, err := os.Create(filepath)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

	} else {
		http.Error(res, "filepath is empty", http.StatusBadRequest)
		return
	}
	// responseData := map[string]interface{}{
	// 	"status":  "Success",
	// 	"message": "",
	// 	"data": map[string]interface{}{
	// 		"image": header.Filename,
	// 	},
	// }
	// responseJson := gjson.New(responseData)
	// // 将 JSON 数据写入响应
	// res.Header().Set("Content-Type", "application/json")
	// res.Write(responseJson.MustToJson())
}
