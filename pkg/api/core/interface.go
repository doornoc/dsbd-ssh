package core

import (
	"github.com/doornoc/dsbd-ssh/pkg/api/core/remote"
	"github.com/google/uuid"
	"time"
)

// channel定義(websocketで使用)
var Clients = make(map[uuid.UUID]*Remote)
var Broadcast = make(chan RemoteResult)

type Remote struct {
	Remote        *remote.Remote
	StartAt       time.Time
	LastUpdatedAt time.Time
}

// websocket用
type RemoteResult struct {
	ID          uint      `json:"id"`
	Err         string    `json:"error"`
	CreatedAt   time.Time `json:"created_at"`
	TicketID    uint      `json:"ticket_id"`
	UserToken   string    `json:"user_token"`
	AccessToken string    `json:"access_token"`
	UserID      uint      `json:"user_id"`
	UserName    string    `json:"user_name"`
	GroupID     uint      `json:"group_id"`
	Admin       bool      `json:"admin"`
	Message     string    `json:"message"`
}
