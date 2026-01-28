// Copyright (C) godoc. 2026-present.
//
// Created at 2026-01-28, by liasica

package godoc

import (
	"fmt"
)

type Logger struct {
}

func (l *Logger) Printf(format string, v ...interface{}) {
	fmt.Println(fmt.Sprintf(format, v...))
}

func NewLogger() *Logger {
	return &Logger{}
}
