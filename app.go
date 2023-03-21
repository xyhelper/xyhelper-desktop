package main

import (
	"context"
	"fmt"
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
