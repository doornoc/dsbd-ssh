package remote

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"strconv"
	"time"
)

func (r *Remote) SSHShell() error {
	consoleLog := ""
	sshConfig := &ssh.ClientConfig{
		User: r.Device.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(r.Device.Password),
			ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
				if len(questions) == 0 {
					return nil, nil
				}
				if len(questions) != 1 {
					return nil, fmt.Errorf("too complex questionnaire: %#v", questions)
				}
				return []string{r.Device.Password}, nil
			}),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		HostKeyAlgorithms: []string{
			ssh.KeyAlgoRSA,
			ssh.KeyAlgoDSA,
			ssh.KeyAlgoECDSA256,
			ssh.KeyAlgoSKECDSA256,
			ssh.KeyAlgoECDSA384,
			ssh.KeyAlgoECDSA521,
			ssh.KeyAlgoED25519,
			ssh.KeyAlgoSKED25519,
			ssh.KeyAlgoRSASHA256,
			ssh.KeyAlgoRSASHA512,
		},
		Timeout: 30 * time.Second,
	}

	sshConfig.KeyExchanges = []string{
		"diffie-hellman-group1-sha1",
		"diffie-hellman-group14-sha1",
		"diffie-hellman-group14-sha256",
		"diffie-hellman-group16-sha512",
		"ecdh-sha2-nistp256",
		"ecdh-sha2-nistp384",
		"ecdh-sha2-nistp521",
		"curve25519-sha256@libssh.org",
		"curve25519-sha256",
		"diffie-hellman-group-exchange-sha256",
		"diffie-hellman-group-exchange-sha1",
	}

	client, err := ssh.Dial("tcp", r.Device.Hostname+":"+strconv.Itoa(int(r.Device.Port)), sshConfig)
	if err != nil {
		fmt.Println("[SSH Dial]", err)
		return nil
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		fmt.Println("[client.NewSession]", err)
		return nil
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		fmt.Println("[session.StdinPipe]", err)
		return nil
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		fmt.Println("[session.StdoutPipe]", err)
		return nil
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	//term := os.Getenv("TERM")
	term := "xterm-256color"
	// request pty session
	//width, height, err := getTerminalSize()
	//if err != nil {
	//	err = fmt.Errorf("get terminal size failed: %v", err)
	//	return
	//}
	err = session.RequestPty(term, 40, 100, modes)
	if err != nil {
		fmt.Println("[session.RequestPty]", err)
		return nil
	}

	err = session.Shell()
	if err != nil {
		fmt.Println("[session.Shell]", err)
		return nil
	}

	//stdoutBeforeTime := time.Now().Add(time.Second * -5)
	r.StdoutLastUpdateTime = time.Now()
	if len(r.Log) == 0 {
		r.Log = append(r.Log, Log{
			OutputStr:  "",
			OutputByte: []byte{},
		})
	}

	// stdout
	go func() {
		buf := make([]byte, 1000000)

		var err error = nil
	OutCancel:
		for err == nil {
			select {
			case <-r.OutCancelCh:
				break OutCancel
			default:
				n, err := stdout.Read(buf)
				// Update time
				r.StdoutLastUpdateTime = time.Now()
				for _, outCh := range r.OutCh {
					outCh <- buf[:n]
				}
				consoleLog += string(buf[:n])
				r.Log[len(r.Log)-1].OutputByte = append(r.Log[len(r.Log)-1].OutputByte, buf[:n]...)
				r.Log[len(r.Log)-1].OutputStr += string(buf[:n])

				//log.Printf("\n[%d]> %d\n", n, string(buf[:n]))
				fmt.Printf(string(buf[:n]))
				if err != nil {
					fmt.Println("[*normal* stdout finish]", err)
					close(r.InCancelCh)
					break OutCancel
				}
			}
		}
	}()

	//stdin
InCancel:
	for {
		select {
		case <-r.InCancelCh:
			break InCancel
		case b := <-r.InCh:
			stdin.Write(b)
			r.Log = append(r.Log, []Log{
				{
					InputByte:  b,
					OutputStr:  "",
					OutputByte: []byte{},
				},
			}...)
			consoleLog += string(b)
		}
	}

	for _, cusInCancelCh := range r.CusInCancelCh {
		close(cusInCancelCh)
	}

	for _, cusOutCancelCh := range r.CusOutCancelCh {
		close(cusOutCancelCh)
	}

	return nil
}
