package tool

import (
	"context"
	"go.uber.org/zap"
)

type MyContext struct {
	context.Context
	Logger *zap.SugaredLogger
}

func InitMyContext() *MyContext {
	lg := InitLogger()
	return &MyContext{
		Context: context.Background(),
		Logger:  lg,
	}
}
