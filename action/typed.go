package action

import (
	"sync"
	"time"
)

type ActionType string

const (
	ActionPeriodic              ActionType = "periodic"
	ActionPeriodicUntilResponse ActionType = "periodic_until_response"
	ActionRespondOnReceive      ActionType = "respond_on_receive"
)

type Action struct {
	ID          int        `json:"id"`
	Type        ActionType `json:"type"`
	PeriodMs    int        `json:"period_ms,omitempty"`
	Message     string     `json:"message,omitempty"`   // message to send (for periodic or reply)
	Expect      string     `json:"expect,omitempty"`    // match string for until/auto
	ReplyMsg    string     `json:"reply_msg,omitempty"` // reply message when match (for respond_on_receive)
	Description string     `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at,omitempty"`
}

var (
	actionsMu sync.Mutex
	actions   = make(map[int]*Action)
	nextID    = 1
	filePath  = "actions.json"

	// runner management: key sessionID:actionID -> stop chan
	runnersMu sync.Mutex
	runners   = make(map[string]chan struct{})
)
