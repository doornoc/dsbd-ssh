package core

import (
	"context"
	"fmt"
	"github.com/doornoc/dsbd-ssh/pkg/api/core/remote"
	"github.com/doornoc/dsbd-ssh/pkg/api/core/tool"
	uuid2 "github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
	"time"
)

type server struct {
	UnimplementedRemoteServiceServer
}

func NewServer() *server {
	return &server{}
}

func (s server) Connect(ctx context.Context, connectReq *ConnectRequest) (*ConnectResponse, error) {
	uuid, err := uuid2.NewUUID()
	if err != nil {
		return nil, err
	}

	if _, isExist := Clients[uuid]; isExist {
		errorValue := fmt.Sprintf("UUID exists. The value is %#v", uuid)
		return nil, status.Errorf(codes.Unimplemented, errorValue)
	}

	Clients[uuid] = &Remote{
		Remote: &remote.Remote{
			Device: remote.Device{
				Name:     uuid.String(),
				Hostname: connectReq.Account.Hostname,
				Port:     uint(connectReq.Account.Port),
				User:     connectReq.Account.Username,
				Password: connectReq.Account.Password,
			},
			InCh:           make(chan []byte),
			OutCh:          make(map[uuid2.UUID](chan []byte)),
			InCancelCh:     make(chan struct{}),
			OutCancelCh:    make(chan struct{}),
			CusInCancelCh:  make(map[uuid2.UUID](chan struct{})),
			CusOutCancelCh: make(map[uuid2.UUID](chan struct{})),
		},
		StartAt:       time.Time{},
		LastUpdatedAt: time.Time{},
	}

	switch connectReq.Account.Type {
	case Type_SSH:
		go Clients[uuid].Remote.SSHShell()
	}

	return &ConnectResponse{
		Uuid: uuid.String(),
	}, nil
}

func (s server) DisConnect(ctx context.Context, disConnectReq *DisconnectRequest) (*Result, error) {
	uuid := uuid2.MustParse(disConnectReq.Uuid)
	if _, isExist := Clients[uuid]; !isExist {
		errorValue := fmt.Sprintf("UUID is not exists. The value is %#v", uuid)
		return nil, status.Errorf(codes.Unimplemented, errorValue)
	}
	close(Clients[uuid].Remote.InCancelCh)
	close(Clients[uuid].Remote.OutCancelCh)
	delete(Clients, uuid)

	return &Result{}, nil

}
func (s server) Remote(stream RemoteService_RemoteServer) error {
	remoteReq, err := stream.Recv()
	if err != nil {
		return err
	}
	if remoteReq.Uuid == "" {
		return fmt.Errorf("No uuid...\n")
	}
	remoteUUID := uuid2.MustParse(remoteReq.Uuid)
	if _, isExist := Clients[remoteUUID]; !isExist {
		return fmt.Errorf("UUID is not exists. The value is %#v", remoteUUID)
	}

	sessionID, err := uuid2.NewUUID()
	if err != nil {
		return err
	}

	r := Clients[remoteUUID].Remote
	r.OutCh[sessionID] = make(chan []byte)
	r.CusInCancelCh[sessionID] = make(chan struct{})
	r.CusOutCancelCh[sessionID] = make(chan struct{})
	stream.Send(&RemoteResponse{
		Output: []byte("Loading...\n"),
	})
	go func() {
		for {
			select {
			case <-r.CusOutCancelCh[sessionID]:
				break
			case outCh := <-r.OutCh[sessionID]:
				stream.Send(&RemoteResponse{
					Output: outCh,
				})
			}
		}
	}()

	for {
		remoteReq, err = stream.Recv()
		if err == io.EOF {
			continue
		}

		if err != nil {
			return err
		}

		commands, err := remote.LoadTemplate(string(remoteReq.Input))
		if err != nil {
			return err
		}

		_, err = r.Exec(sessionID, commands)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s server) RemoteInput(ctx context.Context, remoteInReq *RemoteRequest) (*Result, error) {
	remoteUUID := uuid2.MustParse(remoteInReq.Uuid)
	if _, isExist := Clients[remoteUUID]; !isExist {
		errorValue := fmt.Sprintf("UUID is not exists. The value is %#v", remoteUUID)
		return &Result{Ok: false}, status.Errorf(codes.Unimplemented, errorValue)
	}

	sessionID, err := uuid2.NewUUID()
	if err != nil {
		return &Result{Ok: false}, err
	}

	r := Clients[remoteUUID].Remote

	commands, err := remote.LoadTemplate(string(remoteInReq.Input))
	if err != nil {
		return &Result{Ok: false}, err
	}
	_, err = r.Exec(sessionID, commands)
	if err != nil {
		return &Result{Ok: false}, err
	}

	return &Result{Ok: true}, nil
}

func (s server) RemoteOutputRemoteOutput(req *RemoteOutputRequest, stream RemoteService_RemoteOutputServer) error {
	if req.Uuid == "" {
		return fmt.Errorf("No uuid...\n")
	}
	remoteUUID := uuid2.MustParse(req.Uuid)
	if _, isExist := Clients[remoteUUID]; !isExist {
		return fmt.Errorf("UUID is not exists. The value is %#v", remoteUUID)
	}

	sessionID, err := uuid2.NewUUID()
	if err != nil {
		return err
	}

	r := Clients[remoteUUID].Remote
	r.OutCh[sessionID] = make(chan []byte)
	r.CusInCancelCh[sessionID] = make(chan struct{})
	r.CusOutCancelCh[sessionID] = make(chan struct{})
	stream.Send(&RemoteResponse{
		Output: []byte("Loading...\n"),
	})
	go func() {
		for {
			select {
			case <-r.CusOutCancelCh[sessionID]:
				break
			case outCh := <-r.OutCh[sessionID]:
				stream.Send(&RemoteResponse{
					Output: outCh,
				})
			}
		}
	}()

	return nil
}

func Server() {
	listenPort, err := net.Listen("tcp", fmt.Sprintf(":%d", tool.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	RegisterRemoteServiceServer(grpcServer, NewServer())

	if err = grpcServer.Serve(listenPort); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
