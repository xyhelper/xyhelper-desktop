package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	chatgpt "github.com/xyhelper/chatgpt-go"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// SessionRes is the response from the Session method
type SessionRes struct {
	Auth  bool   `json:"auth"`
	Model string `json:"model"`
}

// Session returns the session information
func (a *App) Session() *SessionRes {
	return &SessionRes{
		Auth:  false,
		Model: "ChatGPTUnofficialProxyAPI",
	}
}

// domReady is called when the DOM is ready
func (a *App) DomReady(ctx context.Context) {
	// Do something
	// a.ctx = ctx
}

// ChatProcessReq is the request for the ChatProcess method
type ChatProcessReq struct {
	Prompt  string `json:"prompt"`
	Options *struct {
		ConversationId  string `json:"conversationId,omitempty"`
		ParentMessageId string `json:"parentMessageId,omitempty"`
	} `json:"options,omitempty"`
	BaseURI     string `json:"baseURI,omitempty"`
	AccessToken string `json:"accessToken,omitempty"`
}

// ChatProcess
func (a *App) ChatProcess(req *ChatProcessReq) {
	// a.chatStop <- false
	ctx := a.ctx
	var err error
	// g.DumpWithType(req)
	g.Log().Debug(ctx, "ChatProcess", req)
	if req.BaseURI != "" {
		BaseURI = req.BaseURI
	}
	g.Log().Debug(ctx, "ChatProcess", BaseURI)
	if req.AccessToken != "" {
		AccessToken = req.AccessToken
	}
	g.Log().Debug(ctx, "ChatProcess", AccessToken)

	// 生成uuid

	cli := chatgpt.NewClient(
		chatgpt.WithAccessToken(AccessToken),
		chatgpt.WithTimeout(120*time.Second),
		chatgpt.WithBaseURI(BaseURI),
		// chatgpt.WithDebug(true),
	)
	message := req.Prompt
	errMsg := map[string]interface{}{
		"role":            "assistant",
		"id":              "",
		"parentMessageId": req.Options.ParentMessageId,
		"conversationId":  req.Options.ConversationId,
		"text":            "思考中...",
	}
	// 设置聊天频道
	chatChannel := "chat"
	if req.Options.ParentMessageId != "" {
		chatChannel = req.Options.ParentMessageId
	}
	// 发送思考中
	runtime.EventsEmit(a.ctx, chatChannel, errMsg)

	var stream *chatgpt.ChatStream
	if req.Options.ConversationId == "" || req.Options.ParentMessageId == "" {
		stream, err = cli.GetChatStream(message)
	} else {
		stream, err = cli.GetChatStream(message, req.Options.ConversationId, req.Options.ParentMessageId)
	}

	if err != nil {
		// 解决帐号切换后会话丢失的问题 使用新会话新求
		if err.Error() == "send message failed: 404 Not Found" {
			g.Log().Errorf(ctx, "获取聊天内容失败: %s,这是首次出现将重试", err.Error())
			errMsg["text"] = "会话丢失,正在重新生成..."
			runtime.EventsEmit(a.ctx, chatChannel, errMsg)
			stream, err = cli.GetChatStream(message)
		}
	}
	if err != nil {
		g.Log().Errorf(ctx, "获取聊天内容失败: %s", err.Error())
		errMsg["text"] = fmt.Sprintf("获取聊天内容失败: %s .", err.Error()) + "请稍后刷新重试,点这里 ➚"
		runtime.EventsEmit(a.ctx, chatChannel, errMsg)
		return
	}
	var answer string

	for text := range stream.Stream {
		g.Log().Printf(ctx, "stream text: %s\n", text.Content)
		answer = text.Content
		// 构造要返回的 JSON 数据
		responseData := map[string]interface{}{
			"role":            "assistant",
			"id":              text.MessageID,
			"parentMessageId": req.Options.ParentMessageId,
			"conversationId":  text.ConversationID,
			"text":            text.Content,
		}

		runtime.EventsEmit(a.ctx, chatChannel, responseData)
		// responseJson := gjson.New(responseData)
		// // runtime.EventsEmit(ctx, "chat", answer)
		// // 将 JSON 数据写入响应
		// res.Write(responseJson.MustToJson())
		// // 增加换行符
		// res.Write([]byte("\n"))
	}

	if stream.Err != nil {
		g.Log().Errorf(ctx, "stream closed with error: %v\n", stream.Err)
		runtime.EventsEmit(a.ctx, chatChannel, errMsg)
	}

	g.Log().Infof(ctx, "q: %s, a: %s\n", message, answer)

}

// StopChat
func (a *App) StopChat() {
	g.Log().Info(a.ctx, "StopChat~~~~~~~~~~~~~~~~~~~~")
}
