package test

import (
	"fmt"

	"google.golang.org/grpc"
)

type TestDialer struct {
	Addr   string
	C      *grpc.ClientConn
	client TestClient
}

func (a *TestDialer) DialReinfra() (err error) {
	fmt.Println("Dialing", a.Addr)
	a.C, err = grpc.Dial(
		a.Addr,
		grpc.WithInsecure(),
		grpc.WithBackoffConfig(grpc.DefaultBackoffConfig),
	)
	if err != nil {
		return err
	}

	a.client = NewTestClient(a.C)
	return
}

func (a *TestDialer) Close() error {
	return a.C.Close()
}

func (a *TestDialer) TestClient() TestClient {
	return a.client
}
