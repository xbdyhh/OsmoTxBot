package tool

import (
	"context"
	"go.uber.org/zap"
	"sync"
)

type MyContext struct {
	context.Context
	Logger *zap.SugaredLogger
	Wg     sync.WaitGroup
}

func InitMyContext() *MyContext {
	lg := InitLogger()
	wg := sync.WaitGroup{}
	return &MyContext{
		Context: context.Background(),
		Logger:  lg,
		Wg:      wg,
	}

}
