package remote

import "time"

type Device struct {
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
	Port     uint   `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	OSType   string `json:"os_type"`
}

type Remote struct {
	Type                 uint   //0:ssh, 1:telnet
	Device               Device `json:"device"`
	InCh                 chan []byte
	OutCh                chan []byte
	InCancelCh           chan struct{}
	OutCancelCh          chan struct{}
	InputCancelCh        chan struct{}
	StdoutLastUpdateTime time.Time
	Log                  []Log
}

type Log struct {
	InputByte  []byte
	OutputStr  string
	OutputByte []byte
}

type command struct {
	Type    string
	Command string
	Code    int
	Option1 string
	Option2 string
}
