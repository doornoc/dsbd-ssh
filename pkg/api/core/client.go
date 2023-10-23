package core

import (
	"context"
	"fmt"
	"github.com/doornoc/dsbd-ssh/pkg/api/core/remote"
	"golang.org/x/term"
	"google.golang.org/grpc"
	"io"
	"log"
	"syscall"
)

func Client(hostname string, port int, username string) error {
	fmt.Println("Please input password...")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		panic(err)
	}

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Can't connect with server %v", err)
	}
	// create stream
	client := NewRemoteServiceClient(conn)
	connectRes, err := client.Connect(context.Background(), &ConnectRequest{
		Account: &Account{
			Type:     0,
			Hostname: hostname,
			Port:     uint32(port),
			Username: username,
			Password: string(bytePassword),
		},
	})
	if err != nil {
		log.Fatalf("Open request error: %v", err)
	}
	uuid := connectRes.Uuid
	log.Printf("[OK] UUID: %s\n", uuid)

	// create stream
	remoteStream, err := client.Remote(context.Background())
	if err != nil {
		log.Fatalf("Open stream error: %v", err)
	}

	remoteStream.Send(&RemoteRequest{Uuid: uuid})

	done := make(chan bool)

	go func() {
		for {
			resp, err := remoteStream.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				log.Fatalf("can not receive %v", err)
			}
			fmt.Printf("%s", string(resp.Output))
		}
	}()

	input := make(chan string)
	go remote.InputKeyLines(input)
	for {
		select {
		case <-done:
			break
		case inputLine := <-input:
			inputCmd := "CMD: " + inputLine + "\nKEY: enter"
			remoteStream.Send(&RemoteRequest{
				Uuid:  "",
				Input: []byte(inputCmd),
			})
		}
	}

	return nil
}
