package core

import (
	"go/ast"
	"go/token"
)

// LogEntry представляет нормализованный лог вызов.
type LogEntry struct {
	Logger      string
	Level       string
	Message     string   // нормализованное текстовое сообщение
	MessageExpr ast.Expr // исходное выражение сообщения
	Pos         token.Pos
	End         token.Pos
	Call        *ast.CallExpr
}
