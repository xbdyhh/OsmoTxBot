package main

import (
	"github.com/xbdyhh/OsmoTxBot/logic"
	"github.com/xbdyhh/OsmoTxBot/tool"
	"github.com/xbdyhh/OsmoTxBot/tool/osmo"
)

func main() {
	ctx := tool.InitMyContext()
	osmo.InitCcontext()
	ctx.Wg.Add(2)
	go logic.FreshPoolMap(ctx)
	go logic.SendOsmoTriTx(ctx)
	ctx.Wg.Wait()
}
