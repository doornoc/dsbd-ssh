package remote

import (
	"testing"
)

func TestLoadTemplate(t *testing.T) {
	data := `
CMD: uname -a {{ str "_" }}
KEY: enter
CMD: uname -a {{ str "{}" }}
KEY: enter
CMD: uname {{ str "{}" }}
KEY: enter
KEY: ctrl-c
OPT: wait:10
OPT: disconnect
`
	LoadTemplate(data)
	t.Log("end")
}
