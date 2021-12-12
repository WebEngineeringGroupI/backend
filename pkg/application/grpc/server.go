package grpc

import (
	"context"
	"fmt"
	"io"

	genproto "github.com/WebEngineeringGroupI/genproto-go/api/v1alpha1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

type server struct {
	genproto.UnimplementedURLShorteningServer
	baseDomain   string
	urlShortener *url.SingleURLShortener
	loadBalancer *url.LoadBalancer
}

func (s *server) ShortURLs(shortURLsServer genproto.URLShortening_ShortURLsServer) error {
	for {
		request, err := shortURLsServer.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		shortURL, err := s.urlShortener.HashFromURL(shortURLsServer.Context(), request.Url)
		if err != nil {
			err := shortURLsServer.Send(&genproto.ShortURLsResponse{
				Result: &genproto.ShortURLsResponse_Error_{
					Error: &genproto.ShortURLsResponse_Error{
						Url:   request.Url,
						Error: err.Error(),
					},
				},
			})
			if err != nil {
				return err
			}
			continue
		}

		err = shortURLsServer.Send(&genproto.ShortURLsResponse{
			Result: &genproto.ShortURLsResponse_Success_{
				Success: &genproto.ShortURLsResponse_Success{
					LongUrl:  shortURL.OriginalURL.URL,
					ShortUrl: fmt.Sprintf("%s/r/%s", s.baseDomain, shortURL.Hash),
				},
			},
		})
		if err != nil {
			return err
		}
	}
}

func (s *server) BalanceURLs(ctx context.Context, req *genproto.BalanceURLsRequest) (*genproto.BalanceURLsResponse, error) {
	//fixme(fede): use the ctx for cancellation
	balancedURL, err := s.loadBalancer.ShortURLs(ctx, req.GetUrls())
	if err != nil {
		return nil, err
	}
	return &genproto.BalanceURLsResponse{ShortUrl: fmt.Sprintf("%s/lb/%s", s.baseDomain, balancedURL.Hash)}, nil
}

type Config struct {
	BaseDomain                 string
	ShortURLRepository         url.ShortURLRepository
	CustomMetrics              url.Metrics
	LoadBalancedURLsRepository url.LoadBalancedURLsRepository
	EventEmitter               event.Emitter
}

func NewServer(config Config) *grpc.Server {
	grpcServer := grpc.NewServer()
	srv := &server{
		baseDomain:   config.BaseDomain,
		urlShortener: url.NewSingleURLShortener(config.ShortURLRepository, config.CustomMetrics, config.EventEmitter),
		loadBalancer: url.NewLoadBalancer(config.LoadBalancedURLsRepository),
	}

	genproto.RegisterURLShorteningServer(grpcServer, srv)

	reflection.Register(grpcServer)
	return grpcServer
}
