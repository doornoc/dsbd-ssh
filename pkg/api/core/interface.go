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
	Remote  *remote.Remote
	StartAt time.Time
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

func Close(uuid string) {

}

func CloseCh(uuid uuid.UUID) {

	if !Clients[uuid].Remote.ClosedCh.ClosedInCancelCh {
		close(Clients[uuid].Remote.InCancelCh)
		Clients[uuid].Remote.ClosedCh.ClosedInCancelCh = true
	}
	if !Clients[uuid].Remote.ClosedCh.ClosedOutCancelCh {
		close(Clients[uuid].Remote.OutCancelCh)
		Clients[uuid].Remote.ClosedCh.ClosedOutCancelCh = true
	}

	for _, cusCh := range Clients[uuid].Remote.CusCh {
		if !cusCh.ClosedCusCh.ClosedCusInCancelCh {
			close(cusCh.CusInCancelCh)
			cusCh.ClosedCusCh.ClosedCusInCancelCh = true
		}
		if !cusCh.ClosedCusCh.ClosedCusOutCancelCh {
			close(cusCh.CusOutCancelCh)
			cusCh.ClosedCusCh.ClosedCusOutCancelCh = true
		}
		if !cusCh.ClosedCusCh.ClosedOutCh {
			close(cusCh.OutCh)
			cusCh.ClosedCusCh.ClosedOutCh = true
		}
	}

	if !Clients[uuid].Remote.ClosedCh.CloseExitCh {
		close(Clients[uuid].Remote.ExitCh)
		Clients[uuid].Remote.ClosedCh.CloseExitCh = true
	}
	delete(Clients, uuid)
}
