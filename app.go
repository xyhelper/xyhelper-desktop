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
}

// ChatProcess
func (a *App) ChatProcess(req *ChatProcessReq) (err error) {
	// a.chatStop <- false
	ctx := a.ctx
	g.DumpWithType(req)
	cli := chatgpt.NewClient(
		chatgpt.WithAccessToken("hello world"),
		chatgpt.WithTimeout(120*time.Second),
	)
	message := req.Prompt
	errMsg := map[string]interface{}{
		"role":            "assistant",
		"id":              "",
		"parentMessageId": req.Options.ParentMessageId,
		"conversationId":  req.Options.ConversationId,
		"text":            "OPENAI服务器限流,请稍后或刷新重试,点这里 ➚",
	}
	var stream *chatgpt.ChatStream
	if req.Options.ConversationId == "" || req.Options.ParentMessageId == "" {
		stream, err = cli.GetChatStream(message)
	} else {
		stream, err = cli.GetChatStream(message, req.Options.ConversationId, req.Options.ParentMessageId)
	}

	if err != nil {
		g.Log().Errorf(ctx, "获取聊天内容失败: %s", err.Error())

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
		chatChannel := "chat"
		if req.Options.ParentMessageId != "" {
			chatChannel = req.Options.ParentMessageId
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
		runtime.EventsEmit(a.ctx, "chat", errMsg)
	}

	g.Log().Infof(ctx, "q: %s, a: %s\n", message, answer)

	return
}

// StopChat
func (a *App) StopChat() {
	g.Log().Info(a.ctx, "StopChat~~~~~~~~~~~~~~~~~~~~")
}
