package action

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"robotech/logger"
	"robotech/session"
	"robotech/utils"
)

func init() {
	_ = LoadFromFile()
}

func Add(a Action) int {
	actionsMu.Lock()
	defer actionsMu.Unlock()
	a.ID = nextID
	nextID++
	a.CreatedAt = time.Now()
	actions[a.ID] = &a
	_ = saveToFileUnlocked()
	logger.Infof("action added id=%d type=%s desc=%s", a.ID, a.Type, a.Description)
	return a.ID
}

func List() []*Action {
	actionsMu.Lock()
	defer actionsMu.Unlock()
	out := make([]*Action, 0, len(actions))
	for _, a := range actions {
		out = append(out, a)
	}
	return out
}

func Get(id int) *Action {
	actionsMu.Lock()
	defer actionsMu.Unlock()
	return actions[id]
}

func SaveToFile() error {
	actionsMu.Lock()
	defer actionsMu.Unlock()
	return saveToFileUnlocked()
}

func saveToFileUnlocked() error {
	list := make([]*Action, 0, len(actions))
	for _, a := range actions {
		list = append(list, a)
	}

	b, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	logger.Infof("saving %d actions to file %s", len(list), filePath)
	return os.WriteFile(filePath, b, 0644)
}

func LoadFromFile() error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}
	b, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	var list []*Action
	if err := json.Unmarshal(b, &list); err != nil {
		return err
	}
	actionsMu.Lock()
	defer actionsMu.Unlock()
	for _, a := range list {
		actions[a.ID] = a
		if a.ID >= nextID {
			nextID = a.ID + 1
		}
	}
	logger.Infof("loaded %d actions", len(list))
	return nil
}

func Bind(sessionIDs []string, actionIDs []int) error {
	for _, sid := range sessionIDs {
		for _, aid := range actionIDs {
			a := Get(aid)
			if a == nil {
				return fmt.Errorf("action %d not found", aid)
			}
			// mark session-side assignment
			if err := session.AssignAction(sid, aid); err != nil {
				logger.Error("assign action fail:", err)
				continue
			}
			// start runner (deduped)
			logger.Infof("bind action %d to session %s", a.ID, sid)
			go startRunner(sid, a)
		}
	}
	return nil
}

func Unbind(sessionID string, actionID int) {
	// stop runner if exists
	key := runnerKey(sessionID, actionID)
	runnersMu.Lock()
	if ch, ok := runners[key]; ok {
		close(ch)
		delete(runners, key)
	}
	runnersMu.Unlock()
	// unassign on session
	session.UnassignAction(sessionID, actionID)
}

func runnerKey(sessionID string, aid int) string {
	return fmt.Sprintf("%s:%d", sessionID, aid)
}

func startRunner(sessionID string, a *Action) {
	key := runnerKey(sessionID, a.ID)
	// dedupe
	runnersMu.Lock()
	if _, ok := runners[key]; ok {
		runnersMu.Unlock()
		return
	}
	stop := make(chan struct{})
	runners[key] = stop
	runnersMu.Unlock()

	logger.Infof("runner start session=%s action=%d", sessionID, a.ID)

	switch a.Type {
	case ActionPeriodic:
		runPeriodic(sessionID, a, stop)
	case ActionPeriodicUntilResponse:
		runPeriodicUntilResponse(sessionID, a, stop)
	case ActionRespondOnReceive:
		runRespondOnReceive(sessionID, a, stop)
	default:
		logger.Error("unknown action type:", a.Type)
	}

	// cleanup
	runnersMu.Lock()
	if ch, ok := runners[key]; ok && ch == stop {
		delete(runners, key)
	}
	runnersMu.Unlock()
	logger.Infof("runner exit session=%s action=%d", sessionID, a.ID)
}

// parseMessage converts stored message string into raw bytes according to msgType
func parseMessage(msgType session.MsgType, raw string) ([]byte, error) {
	if msgType == session.MsgTypeHex {
		/*clean := strings.ReplaceAll(raw, " ", "")
		// tolerate 0x prefix
		clean = strings.TrimPrefix(clean, "0x")
		return hex.DecodeString(clean)*/
		return utils.HexStringToBytes(raw)
	}
	return []byte(raw), nil
}

func runPeriodic(sessionID string, a *Action, stop chan struct{}) {
	interval := time.Duration(a.PeriodMs) * time.Millisecond
	if interval <= 0 {
		interval = 1 * time.Second
	}
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-stop:
			return
		case <-t.C:
			if !session.IsActionAssigned(sessionID, a.ID) {
				return
			}
			b, err := parseMessage(session.GetSessionMsgType(sessionID), a.Message)
			if err != nil {
				logger.Error("parse message err:", err)
				return
			}
			if err := session.Send(sessionID, b); err != nil {
				logger.Error("session send err:", err)
				return
			}
		}
	}
}

func runPeriodicUntilResponse(sessionID string, a *Action, stop chan struct{}) {
	interval := time.Duration(a.PeriodMs) * time.Millisecond
	if interval <= 0 {
		interval = 1 * time.Second
	}
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-stop:
			return
		case <-t.C:
			if !session.IsActionAssigned(sessionID, a.ID) {
				return
			}
			b, err := parseMessage(session.GetSessionMsgType(sessionID), a.Message)
			if err != nil {
				logger.Error("parse message err:", err)
				return
			}
			if err := session.Send(sessionID, b); err != nil {
				logger.Error("session send err:", err)
				return
			}
			// wait up to 5s or until expected appears in RecvChan
			timeout := time.NewTimer(5 * time.Second)
			select {
			case <-stop:
				timeout.Stop()
				return
			case <-timeout.C:
				// try again next tick
			case recv, ok := <-sessionRecvChan(sessionID):
				if !ok {
					timeout.Stop()
					return
				}
				if a.Expect != "" && strings.Contains(recv, a.Expect) {
					logger.Infof("action %d got expected response on session %s", a.ID, sessionID)
					return
				}
			}
		}
	}
}

func runRespondOnReceive(sessionID string, a *Action, stop chan struct{}) {
	for {
		select {
		case <-stop:
			return
		case msg, ok := <-sessionRecvChan(sessionID):
			if !ok {
				return
			}
			// match Expect substring
			if a.Expect != "" && strings.Contains(msg, a.Expect) {
				// send reply
				b, err := parseMessage(session.GetSessionMsgType(sessionID), a.ReplyMsg)
				if err != nil {
					logger.Error("parse reply err:", err)
					continue
				}
				if err := session.Send(sessionID, b); err != nil {
					logger.Error("send reply err:", err)
					return
				}
			}
		}
	}
}

// helper to get session RecvChan, or closed chan if not exist
func sessionRecvChan(sessionID string) <-chan string {
	s := session.Get(sessionID)
	if s == nil {
		ch := make(chan string)
		close(ch)
		return ch
	}
	return s.RecvChan
}
