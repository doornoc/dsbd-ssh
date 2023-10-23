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

func (s server) DisConnect(ctx context.Context, disConnectReq *DisconnectRequest) (*DisconnectResponse, error) {
	uuid := uuid2.MustParse(disConnectReq.Uuid)
	if _, isExist := Clients[uuid]; !isExist {
		errorValue := fmt.Sprintf("UUID is not exists. The value is %#v", uuid)
		return nil, status.Errorf(codes.Unimplemented, errorValue)
	}
	close(Clients[uuid].Remote.InCancelCh)
	close(Clients[uuid].Remote.OutCancelCh)
	delete(Clients, uuid)

	return &DisconnectResponse{}, nil

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
	r.InCh = make(chan []byte)
	r.OutCh = remote.CreateCh(sessionID)
	r.InCancelCh = make(chan struct{})
	r.OutCancelCh = make(chan struct{})
	r.CusInCancelCh = remote.CreateCusCancelCh(sessionID)
	r.CusOutCancelCh = remote.CreateCusCancelCh(sessionID)

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
