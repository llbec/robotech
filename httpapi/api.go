package httpapi

import (
	"encoding/json"
	"net/http"
	"robotech/action"
	"robotech/server"
	"robotech/session"
	"robotech/websocket"

	"github.com/julienschmidt/httprouter"
)

func writeJSON(w http.ResponseWriter, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(obj)
}

func openTCPServer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type Req struct {
		IP   string `json:"ip"`
		Port int    `json:"port"`
	}
	var req Req
	json.NewDecoder(r.Body).Decode(&req)
	go server.StartTCPServer(req.IP, req.Port)
	writeJSON(w, map[string]string{"status": "ok"})
}

func openTCPClient(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type Req struct {
		IP   string `json:"ip"`
		Port int    `json:"port"`
	}
	var req Req
	json.NewDecoder(r.Body).Decode(&req)
	go server.StartTCPClient(req.IP, req.Port)
	writeJSON(w, map[string]string{"status": "ok"})
}

func listSessions(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	writeJSON(w, session.List())
}

func websocketHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("sessionID")
	websocket.HandleWebSocket(w, r, id)
}

func addAction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req action.Action
	json.NewDecoder(r.Body).Decode(&req)
	id := action.Add(req)
	writeJSON(w, map[string]interface{}{"id": id})
}

func listActions(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	writeJSON(w, action.List())
}

func bindActions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type Req struct {
		SessionIDs []string `json:"session_ids"`
		ActionIDs  []int    `json:"action_ids"`
	}
	var req Req
	json.NewDecoder(r.Body).Decode(&req)
	action.Bind(req.SessionIDs, req.ActionIDs)
	writeJSON(w, map[string]string{"status": "bound"})
}

func closeSession(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type Req struct {
		ID string `json:"id"`
	}
	var req Req
	json.NewDecoder(r.Body).Decode(&req)
	session.Close(req.ID)
	writeJSON(w, map[string]string{"status": "closed"})
}

func closeTCPServer(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	server.StopAllServers()
	writeJSON(w, map[string]string{"status": "servers closed"})
}

func closeTCPClient(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	server.StopAllClients()
	writeJSON(w, map[string]string{"status": "clients closed"})
}
