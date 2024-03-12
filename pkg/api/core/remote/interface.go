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
	OutCh                map[uuid.UUID](chan []byte)
	InCancelCh           chan struct{}
	OutCancelCh          chan struct{}
	CusInCancelCh        map[uuid.UUID](chan struct{})
	CusOutCancelCh       map[uuid.UUID](chan struct{})
	ExitCh               chan struct{}
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
