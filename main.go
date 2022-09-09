package main

import (
	"github.com/xbdyhh/OsmoTxBot/logic"
	"github.com/xbdyhh/OsmoTxBot/tool"
	"github.com/xbdyhh/OsmoTxBot/tool/osmo"
)

func main() {
	ctx := tool.InitMyContext()
	osmo.InitCcontext()
	logic.FreshPoolMap(ctx)
}
