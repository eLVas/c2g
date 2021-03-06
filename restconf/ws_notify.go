package restconf

import (
	"time"

	"github.com/c2stack/c2g/c2"
	"github.com/c2stack/c2g/node"
	"golang.org/x/net/websocket"
)

// Determined using default websocket settings and Chrome 49 and stop watch when it
// timesout out w/o any pings from server.
const PingRate = 30 * time.Second

// websocket library will kill connection after this time. This is mostly unnec.
// for our usage because we actively ping so this just has to be larger than ping rate
const serverSocketTimeout = 2 * PingRate

type WsNotifyService struct {
	Timeout int
	Factory node.Subscriber
}

func (self *WsNotifyService) Handle(ws *websocket.Conn) {
	var rate time.Duration
	if self.Timeout == 0 {
		rate = PingRate
	} else {
		rate = time.Duration(self.Timeout) * time.Millisecond
	}
	conn := &wsconn{
		pinger: time.NewTicker(rate),
		mgr:    node.NewSubscriptionManager(self.Factory, ws, ws),
	}
	defer conn.close()
	ws.Request().Body.Close()
	go conn.keepAlive(ws)
	if err := conn.mgr.Run(); err != nil {
		c2.Info.Printf("unclean terminination of web socket: (%s). other side may have close browser. closing socket.", err)
	}
	// ignore error, other-side is free to disappear at will
	ws.Close()
}

type wsconn struct {
	pinger *time.Ticker
	mgr    *node.SubscriptionManager
}

func (self *wsconn) keepAlive(ws *websocket.Conn) {
	for {
		ws.SetDeadline(time.Now().Add(serverSocketTimeout))
		if fw, err := ws.NewFrameWriter(websocket.PingFrame); err != nil {
			//self.Close()
			return
		} else if _, err = fw.Write([]byte{}); err != nil {
			//self.Close()
			return
		}
		if _, running := <-self.pinger.C; !running {
			return
		}
	}
}

func (self *wsconn) close() {
	self.pinger.Stop()
	self.mgr.Close()
}
