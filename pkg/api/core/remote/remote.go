package remote

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"time"
)

func (r *Remote) SSHShell() {
	consoleLog := ""
	authMethod := []ssh.AuthMethod{}
	if r.Device.Password != "" {
		authMethod = append(authMethod, ssh.Password(r.Device.Password))
		authMethod = append(authMethod, ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
			if len(questions) == 0 {
				return nil, nil
			}
			if len(questions) != 1 {
				return nil, fmt.Errorf("too complex questionnaire: %#v", questions)
			}
			return []string{r.Device.Password}, nil
		}))
	}

	if r.Device.PrivateKey != "" {
		signer, err := ssh.ParsePrivateKey([]byte(r.Device.PrivateKey))
		if err != nil {
			r.Error = status.Error(codes.InvalidArgument, fmt.Sprintf("[private key is wrong...]", err))
			return
		}
		authMethod = append([]ssh.AuthMethod{ssh.PublicKeys(signer)}, authMethod...)
	}

	sshConfig := &ssh.ClientConfig{
		User:            r.Device.User,
		Auth:            authMethod,
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
		r.Error = status.Error(codes.Unknown, fmt.Sprintf("[SSH Dial]", err))
		return
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		r.Error = status.Error(codes.Unknown, fmt.Sprintf("[client.NewSession]", err))
		return
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		r.Error = status.Error(codes.Unknown, fmt.Sprintf("[session.StdinPipe]", err))
		return
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		r.Error = status.Error(codes.Unknown, fmt.Sprintf("[session.StdoutPipe]", err))
		return
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
		r.Error = status.Error(codes.Unknown, fmt.Sprintf("[session.RequestPty]", err))
		return
	}

	err = session.Shell()
	if err != nil {
		r.Error = status.Error(codes.Unknown, fmt.Sprintf("[session.Shell]", err))
		return
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
				r.ClosedCh.ClosedOutCancelCh = true
				break OutCancel
			default:
				n, err := stdout.Read(buf)
				// Update time
				r.StdoutLastUpdateTime = time.Now()
				for _, cusCh := range r.CusCh {
					if !cusCh.ClosedCusCh.ClosedCusOutCancelCh {
						cusCh.OutCh <- buf[:n]
					}
				}
				consoleLog += string(buf[:n])
				r.LastUpdatedAt = time.Now()
				r.Log[len(r.Log)-1].OutputByte = append(r.Log[len(r.Log)-1].OutputByte, buf[:n]...)
				r.Log[len(r.Log)-1].OutputStr += string(buf[:n])

				//log.Printf("\n[%d]> %d\n", n, string(buf[:n]))
				//fmt.Printf(string(buf[:n]))
				if err != nil {
					//fmt.Println("[*normal* stdout finish]", err)
					close(r.InCancelCh)
					r.ClosedCh.ClosedInCancelCh = true
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
			r.ClosedCh.ClosedInCancelCh = true
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

	for _, cusCh := range r.CusCh {
		if !cusCh.ClosedCusCh.ClosedCusInCancelCh {
			close(cusCh.CusInCancelCh)
			cusCh.ClosedCusCh.ClosedCusInCancelCh = true
		}
		if !cusCh.ClosedCusCh.ClosedCusOutCancelCh {
			close(cusCh.CusOutCancelCh)
			cusCh.ClosedCusCh.ClosedCusOutCancelCh = true
		}
	}

	if !r.ClosedCh.CloseExitCh {
		close(r.ExitCh)
		r.ClosedCh.CloseExitCh = true
	}
	return
}
