package main

import (
	"io/ioutil"
	"net/http"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	chatgpt "github.com/xyhelper/chatgpt-go"
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
	// 注册聊天处理函数
	mux.HandleFunc("/api/chat-process", ApiChatProcess)
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

// ApiChatProcessReq 聊天请求
type ApiChatProcessReq struct {
	Prompt  string `json:"prompt"`
	Options *struct {
		ConversationId  string `json:"conversationId,omitempty"`
		ParentMessageId string `json:"parentMessageId,omitempty"`
	} `json:"options,omitempty"`
}

// ApiChatProcess 处理聊天请求
func ApiChatProcess(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	// 获取请求的方法
	method := req.Method
	// 如果不是POST请求，则返回错误
	if method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// 解析请求的 JSON 数据
	var reqData ApiChatProcessReq
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		g.Log().Errorf(ctx, "读取请求的 JSON 数据失败: %s", err.Error())
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	err = gconv.Struct(reqBody, &reqData)
	if err != nil {
		g.Log().Errorf(ctx, "解析请求的 JSON 数据失败: %s", err.Error())
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	g.Dump(reqData)

	cli := chatgpt.NewClient(
		chatgpt.WithAccessToken("hello world"),
	)
	message := reqData.Prompt
	stream, err := cli.GetChatStream(message)
	if err != nil {
		g.Log().Errorf(ctx, "获取聊天内容失败: %s", err.Error())
		res.WriteHeader(http.StatusTooManyRequests)
		return
	}
	var answer string
	res.Header().Set("Content-Type", "text")
	res.Header().Set("Cache-Control", "no-cache")
	res.Header().Set("Connection", "keep-alive")

	for text := range stream.Stream {
		g.Log().Printf(ctx, "stream text: %s\n", text.Content)
		answer = text.Content
		// 构造要返回的 JSON 数据
		responseData := map[string]interface{}{
			"role":            "assistant",
			"id":              text.MessageID,
			"parentMessageId": reqData.Options.ParentMessageId,
			"conversationId":  text.ConversationID,
			"text":            text.Content,
		}
		responseJson := gjson.New(responseData)
		// runtime.EventsEmit(ctx, "chat", answer)
		// 将 JSON 数据写入响应
		res.Write(responseJson.MustToJson())
		// 增加换行符
		res.Write([]byte("\n"))
	}
	if stream.Err != nil {
		g.Log().Errorf(ctx, "stream closed with error: %v\n", stream.Err)
	}

	g.Log().Infof(ctx, "q: %s, a: %s\n", message, answer)
	// //  构造要返回的 JSON 数据
	// responseData := map[string]interface{}{
	// 	"role":            "assistant",
	// 	"id":              text.MessageID,
	// 	"parentMessageId": reqData.Options.ParentMessageId,
	// 	"conversationId":  text.ConversationID,
	// 	"text":            text.Content,
	// }
	// responseJson := gjson.New(responseData)
	// // 将 JSON 数据写入响应
	// res.Header().Set("Content-Type", "text/event-stream")
	// res.Write(responseJson.MustToJson())

}
