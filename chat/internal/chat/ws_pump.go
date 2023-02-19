package chat

import (
	"context"
	"log"

	"github.com/irth/chanutil"
	"github.com/irth/wsrpc"
)

func websocketPump(ctx context.Context, conn *wsrpc.Conn) chan wsrpc.Errable {
	cmds := make(chan wsrpc.Errable, 1)
	go func() {
		for ctx.Err() == nil {
			cmd, err := conn.Decode()
			if err != nil {
				log.Printf("chat: websocket decode error: %s", err.Error())
				close(cmds)
				return
			}
			chanutil.Put(ctx, cmds, cmd)
		}
	}()
	return cmds
}
