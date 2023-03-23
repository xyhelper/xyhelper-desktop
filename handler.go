package main

import (
	"net/http"

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
			"reverseProxy": "https://freechat.xyhelper.cn/backend-api/conversation",
			"accessToken":  AccessToken,
		},
	}
	responseJson := gjson.New(responseData)
	// 将 JSON 数据写入响应
	res.Header().Set("Content-Type", "application/json")
	res.Write(responseJson.MustToJson())
}
