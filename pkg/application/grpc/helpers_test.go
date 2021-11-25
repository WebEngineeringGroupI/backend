package grpc_test

import (
	"context"
	"net"

	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gogrpc "google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/WebEngineeringGroupI/backend/pkg/application/grpc"
)

func newTestingConnection(config grpc.Config) (*gogrpc.ClientConn, context.CancelFunc) {
	listener := bufconn.Listen(1024 * 1024)
	ctx, cancel := context.WithCancel(context.Background())

	server := grpc.NewServer(config)
	go func() {
		defer ginkgo.GinkgoRecover()
		err := server.Serve(listener)
		Expect(err).ToNot(HaveOccurred())
	}()

	conn, err := gogrpc.DialContext(ctx, "bufnet", gogrpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
		return listener.Dial()
	}), gogrpc.WithInsecure(), gogrpc.WithBlock())
	ExpectWithOffset(1, err).ToNot(HaveOccurred())

	go func() {
		<-ctx.Done()
		conn.Close()
		server.Stop()
	}()

	return conn, cancel
}

type FakeURLValidator struct {
	returnValidURL bool
	returnError    error
}

func (f *FakeURLValidator) shouldReturnValidURL(validURL bool) {
	f.returnValidURL = validURL
}

func (f *FakeURLValidator) shouldReturnError(err error) {
	f.returnError = err
}

func (f *FakeURLValidator) ValidateURL(url string) (bool, error) {
	return f.returnValidURL, f.returnError
}

func (f *FakeURLValidator) ValidateURLs(urls []string) (bool, error) {
	return f.returnValidURL, f.returnError
}
