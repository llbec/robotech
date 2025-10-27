package httpapi

import (
	"encoding/json"
	"net/http"

	"robotech/action"
	"robotech/logger"
	"robotech/server"
	"robotech/session"
	"robotech/utils"
	"robotech/websocket"

	"github.com/julienschmidt/httprouter"
)

type genericResp struct {
	Ok   bool   `json:"ok"`
	Msg  string `json:"msg,omitempty"`
	Data any    `json:"data,omitempty"`
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

// Handlers

func handleStartServer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var body struct {
		IP   string `json:"ip"`
		Port int    `json:"port"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if body.IP == "" {
		body.IP = "0.0.0.0"
	}
	if body.Port == 0 {
		writeJSON(w, genericResp{Ok: false, Msg: "port required"})
		return
	}
	if err := server.StartTCPServer(body.IP, body.Port); err != nil {
		writeJSON(w, genericResp{Ok: false, Msg: err.Error()})
		return
	}
	writeJSON(w, genericResp{Ok: true})
}

func handleStopServer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	addr := ps.ByName("addr")
	if addr == "" {
		writeJSON(w, genericResp{Ok: false, Msg: "addr required"})
		return
	}
	server.StopTCPServer(addr)
	writeJSON(w, genericResp{Ok: true})
}

func handleStartClient(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var body struct {
		IP   string `json:"ip"`
		Port int    `json:"port"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if body.IP == "" || body.Port == 0 {
		writeJSON(w, genericResp{Ok: false, Msg: "ip/port required"})
		return
	}
	if err := server.StartTCPClient(body.IP, body.Port); err != nil {
		writeJSON(w, genericResp{Ok: false, Msg: err.Error()})
		return
	}
	writeJSON(w, genericResp{Ok: true})
}

func handleSetSessionMsgType(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var body struct {
		MsgType session.MsgType `json:"msg_type"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if body.MsgType == "" {
		writeJSON(w, genericResp{Ok: false, Msg: "msg_type required"})
		return
	}
	id := ps.ByName("id")
	if id == "" {
		writeJSON(w, genericResp{Ok: false, Msg: "id required"})
		return
	}
	session.SetSessionMsgType(id, body.MsgType)
	writeJSON(w, genericResp{Ok: true})
}

func handleSend(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var body struct {
		Msg string `json:"msg"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if body.Msg == "" {
		writeJSON(w, genericResp{Ok: false, Msg: "msg required"})
		return
	}
	id := ps.ByName("id")
	if id == "" {
		writeJSON(w, genericResp{Ok: false, Msg: "id required"})
		return
	}
	msgType := session.GetSessionMsgType(id)
	if msgType == session.MsgTypeHex {
		body.Msg = utils.BytesToHexString([]byte(body.Msg))
	}
	if err := session.Send(id, []byte(body.Msg)); err != nil {
		writeJSON(w, genericResp{Ok: false, Msg: err.Error()})
		return
	}
	writeJSON(w, genericResp{Ok: true})
}

func handleStopClient(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	addr := ps.ByName("addr")
	if addr == "" {
		writeJSON(w, genericResp{Ok: false, Msg: "addr required"})
		return
	}
	server.StopTCPClient(addr)
	writeJSON(w, genericResp{Ok: true})
}

func handleListSessions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s := session.List()
	writeJSON(w, genericResp{Ok: true, Data: s})
}

func handleCloseSession(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if id == "" {
		writeJSON(w, genericResp{Ok: false, Msg: "id required"})
		return
	}
	session.Close(id)
	writeJSON(w, genericResp{Ok: true})
}

func handleWS(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if id == "" {
		writeJSON(w, genericResp{Ok: false, Msg: "id required"})
		return
	}
	websocket.HandleWS(w, r, id)
}

func handleAddAction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var a action.Action
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		writeJSON(w, genericResp{Ok: false, Msg: "invalid body"})
		return
	}
	if a.Type == "" {
		writeJSON(w, genericResp{Ok: false, Msg: "action type required"})
		return
	}
	id := action.Add(a)
	writeJSON(w, genericResp{Ok: true, Data: map[string]int{"id": id}})
}

func handleListActions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	list := action.List()
	writeJSON(w, genericResp{Ok: true, Data: list})
}

func handleBindActions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var body struct {
		SessionIDs []string `json:"session_ids"`
		ActionIDs  []int    `json:"action_ids"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if len(body.SessionIDs) == 0 || len(body.ActionIDs) == 0 {
		writeJSON(w, genericResp{Ok: false, Msg: "session_ids and action_ids required"})
		return
	}
	if err := action.Bind(body.SessionIDs, body.ActionIDs); err != nil {
		writeJSON(w, genericResp{Ok: false, Msg: err.Error()})
		return
	}
	writeJSON(w, genericResp{Ok: true})
}

func handleUnbindActions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var body struct {
		SessionID string `json:"session_id"`
		ActionID  int    `json:"action_id"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if body.SessionID == "" || body.ActionID == 0 {
		writeJSON(w, genericResp{Ok: false, Msg: "session_id and action_id required"})
		return
	}
	action.Unbind(body.SessionID, body.ActionID)
	writeJSON(w, genericResp{Ok: true})
}

func handleSetLogfile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var body struct {
		Path string `json:"path"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if body.Path == "" {
		writeJSON(w, genericResp{Ok: false, Msg: "path required"})
		return
	}
	logger.Init(body.Path)
	writeJSON(w, genericResp{Ok: true})
}
