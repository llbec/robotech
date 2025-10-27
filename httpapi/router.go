package httpapi

import (
	"github.com/julienschmidt/httprouter"
)

func NewRouter() *httprouter.Router {
	r := httprouter.New()

	// TCP Server / Client
	r.POST("/api/tcp/server/start", handleStartServer)
	r.POST("/api/tcp/server/stop/:addr", handleStopServer) // addr is ip:port
	r.POST("/api/tcp/client/start", handleStartClient)
	r.POST("/api/tcp/client/stop/:addr", handleStopClient)

	// sessions
	r.GET("/api/sessions", handleListSessions)
	r.POST("/api/session/msg_type/:id", handleSetSessionMsgType)
	r.POST("/api/session/send/:id", handleSend)
	r.POST("/api/session/close/:id", handleCloseSession)

	// websocket subscribe
	r.GET("/ws/session/:id", handleWS)

	// actions
	r.POST("/api/action/add", handleAddAction)
	r.GET("/api/action/list", handleListActions)
	r.POST("/api/action/bind", handleBindActions)
	r.POST("/api/action/unbind", handleUnbindActions)

	// logfile
	r.POST("/api/logfile", handleSetLogfile)

	return r
}
