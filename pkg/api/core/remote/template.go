package remote

import (
	"bufio"
	"bytes"
	"log"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var funcs = template.FuncMap{
	"str": str,
}

func (r *Remote) TemplateApply(templateStr string) error {
	templateArray, err := LoadTemplate(templateStr)
	if err != nil {
		return err
	}
	switch r.Type {
	case 0:
		go r.SSHShell()
	}

	//input := make(chan string)
	//go InputKeyLines(input)
	var commandByte []byte
	var prevCommandByte []byte
InputCancel:
	for _, templateOne := range templateArray {
		for {
			time.Sleep(1 * time.Second)
			if time.Now().Before(r.StdoutLastUpdateTime.Add(time.Second * 5)) {
				continue
			}

			select {
			case <-r.InputCancelCh:
				break InputCancel
			default:
				switch templateOne.Type {
				case "CMD":
					commandByte = []byte(templateOne.Command)
				case "KEY":
					if len(prevCommandByte) != 0 && byteEqual(convertToLF(prevCommandByte), convertToLF(r.Log[len(r.Log)-1].OutputByte)) {
						continue
					}
					commandByte = append(commandByte, byte(templateOne.Code))
					r.InCh <- commandByte
					prevCommandByte = commandByte
					commandByte = []byte{}
				case "OPT":
					switch templateOne.Command {
					case "disconnect":
						close(r.InCancelCh)
						close(r.OutCancelCh)
						break
					default:
						if strings.Contains(templateOne.Command, "wait:") {
							waitTime, err := strconv.Atoi(templateOne.Command[5:])
							if err != nil {
								waitTime = 1
							}
							time.Sleep(time.Duration(waitTime) * time.Second)
						} else if strings.Contains(templateOne.Command, "wait_str:") {
							if len(prevCommandByte) != 0 && strings.Contains(string(r.Log[len(r.Log)-1].OutputByte), templateOne.Command[9:]) {
								continue
							}
						}
					}
				}
			}
			break
		}

	}

	return nil
}

func LoadTemplate(templateData string) ([]command, error) {
	var commands []command = []command{}

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
			commands = append(commands, command{Type: "CMD", Command: text[5:]})
		case "OPT: ":
			commands = append(commands, command{Type: "OPT", Command: text[5:]})
		case "KEY: ":
			keyCode := 10
			switch strings.Join(strings.Fields(text[5:]), "") {
			case "enter":
				keyCode = 10
			}
			commands = append(commands, command{Type: "KEY", Code: keyCode})
		}
	}
	return commands, nil
}

func str(command string) string {
	return command
}
