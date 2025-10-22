package action

import (
	"sync"
	"time"
)

type ActionType string

const (
	ActionPeriodic              ActionType = "periodic"                // 5.1
	ActionPeriodicUntilResponse ActionType = "periodic_until_response" // 5.2
	ActionRespondOnReceive      ActionType = "respond_on_receive"      // 5.3
)

type Action struct {
	ID          int        `json:"id"`
	Type        ActionType `json:"type"`
	PeriodMs    int        `json:"period_ms,omitempty"` // for periodic types
	Message     string     `json:"message,omitempty"`
	Expect      string     `json:"expect,omitempty"` // for until-response
	Description string     `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

var (
	actions    = make(map[int]*Action)
	actionsMu  sync.Mutex
	nextID     = 1
	configFile = "actions.json"
)
