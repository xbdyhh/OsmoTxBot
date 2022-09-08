package logic

import (
	"github.com/xbdyhh/OsmoTxBot/tool"
	"github.com/xbdyhh/OsmoTxBot/tool/osmo"
	"testing"
)

func TestFreshPoolMap(t *testing.T) {
	ctx := tool.InitMyContext()
	osmo.InitCcontext()
	FreshPoolMap(ctx)
	SendOsmoTriTx(ctx)
}
