package remote

import (
	"github.com/google/uuid"
	"time"
)

type Device struct {
	Name       string `json:"name"`
	Hostname   string `json:"hostname"`
	Port       uint   `json:"port"`
	User       string `json:"user"`
	Password   string `json:"password"`
	PrivateKey string `json:"private_key"`
	OSType     string `json:"os_type"`
}

type Remote struct {
	Type                 uint   //0:ssh, 1:telnet
	Device               Device `json:"device"`
	InCh                 chan []byte
	InCancelCh           chan struct{}
	OutCancelCh          chan struct{}
	CusCh                map[uuid.UUID]*CusChannel
	ExitCh               chan struct{}
	ClosedCh             ClosedChStatus
	StdoutLastUpdateTime time.Time
	IsTemplate           bool
	Log                  []Log
	Error                error
	LastUpdatedAt        time.Time
}

type Log struct {
	InputByte  []byte
	OutputStr  string
	OutputByte []byte
}

type Command struct {
	Type    string
	Command string
	Code    int
	Option1 string
	Option2 string
}

type CusChannel struct {
	OutCh          chan []byte
	CusInCancelCh  chan struct{}
	CusOutCancelCh chan struct{}
	ClosedCusCh    *ClosedCusChStatus
}

type ClosedChStatus struct {
	ClosedInCancelCh  bool
	ClosedOutCancelCh bool
	CloseExitCh       bool
}
type ClosedCusChStatus struct {
	ClosedOutCh          bool
	ClosedCusInCancelCh  bool
	ClosedCusOutCancelCh bool
}
