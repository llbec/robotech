package httpapi

import (
	"github.com/julienschmidt/httprouter"
)

func NewRouter() *httprouter.Router {
	r := httprouter.New()

	r.POST("/api/tcp/server/open", openTCPServer)
	r.POST("/api/tcp/client/open", openTCPClient)
	r.GET("/api/tcp/sessions", listSessions)
	r.GET("/api/ws/:sessionID", websocketHandler)

	r.POST("/api/action/add", addAction)
	r.GET("/api/action/list", listActions)
	r.POST("/api/action/bind", bindActions)
	r.POST("/api/session/close", closeSession)
	r.POST("/api/server/close", closeTCPServer)
	r.POST("/api/client/close", closeTCPClient)

	return r
}
