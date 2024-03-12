package remote

import (
	"bufio"
	"bytes"
	"github.com/google/uuid"
	"log"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var funcs = template.FuncMap{
	"str": str,
}

func (r *Remote) Exec(uuid uuid.UUID, commands []Command) (int, error) {
	var commandByte []byte
	var prevCommandByte []byte

InputCancel:
	for _, command := range commands {
		for {
			time.Sleep(1 * time.Second)
			if r.IsTemplate && time.Now().Before(r.StdoutLastUpdateTime.Add(time.Second*5)) {
				continue
			}

			select {
			case <-r.CusCh[uuid].CusInCancelCh:
				r.CusCh[uuid].ClosedCusCh.ClosedCusInCancelCh = true
				break InputCancel
			default:
				switch command.Type {
				case "CMD":
					commandByte = []byte(command.Command)
					r.InCh <- commandByte
				case "KEY":
					if len(prevCommandByte) != 0 && byteEqual(convertToLF(prevCommandByte), convertToLF(r.Log[len(r.Log)-1].OutputByte)) {
						continue
					}
					commandByte = []byte{byte(command.Code)}
					r.InCh <- commandByte
					prevCommandByte = commandByte
					commandByte = []byte{}
				case "OPT":
					switch command.Command {
					case "disconnect":
						close(r.InCancelCh)
						r.ClosedCh.ClosedInCancelCh = true
						close(r.OutCancelCh)
						r.ClosedCh.ClosedOutCancelCh = true
						break
					default:
						if strings.Contains(command.Command, "wait:") {
							waitTime, err := strconv.Atoi(command.Command[5:])
							if err != nil {
								waitTime = 1
							}
							time.Sleep(time.Duration(waitTime) * time.Second)
						} else if strings.Contains(command.Command, "wait_str:") {
							if len(prevCommandByte) != 0 && strings.Contains(string(r.Log[len(r.Log)-1].OutputByte), command.Command[9:]) {
								continue
							}
						}
					}
				}
			}
			break
		}
	}

	return 0, nil
}

func (r *Remote) TemplateApply(templateStr string) error {
	templateArray, err := LoadTemplate(templateStr)
	r.IsTemplate = true
	if err != nil {
		return err
	}
	switch r.Type {
	case 0:
		go r.SSHShell()
	}

	sessionID, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	r.Exec(sessionID, templateArray)

	return nil
}

func LoadTemplate(templateData string) ([]Command, error) {
	var commands []Command = []Command{}

	cmdParse, err := template.New("").Funcs(funcs).Parse(templateData)
	if err != nil {
		log.Println(err)
		return commands, err
	}
	var b bytes.Buffer
	if err = cmdParse.Execute(&b, templateData); err != nil {
		return commands, err
	}

	buf := bytes.NewBufferString(b.String())
	scanner := bufio.NewScanner(buf)

	for scanner.Scan() {
		text := scanner.Text()
		if len(text) <= 4 {
			continue
		}
		switch text[0:5] {
		case "CMD: ":
			commands = append(commands, Command{Type: "CMD", Command: text[5:]})
		case "OPT: ":
			commands = append(commands, Command{Type: "OPT", Command: text[5:]})
		case "KEY: ":
			keyCode := 10
			switch strings.Join(strings.Fields(text[5:]), "") {
			case "enter":
				keyCode = 10
			}
			commands = append(commands, Command{Type: "KEY", Code: keyCode})
		}
	}
	return commands, nil
}

func str(command string) string {
	return command
}
