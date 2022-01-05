package grpc

import (
	"context"
	"fmt"
	"io"

	genproto "github.com/WebEngineeringGroupI/genproto-go/api/v1alpha1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/clock"
)

type server struct {
	genproto.UnimplementedURLShorteningServer
	baseDomain   string
	urlShortener *url.SingleURLShortener
	loadBalancer *url.LoadBalancerService
}

func (s *server) ShortURLs(shortURLsServer genproto.URLShortening_ShortURLsServer) error {
	for {
		request, err := shortURLsServer.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return status.Errorf(codes.Internal, err.Error())
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
				return status.Errorf(codes.Internal, err.Error())
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
			return status.Errorf(codes.Internal, err.Error())
		}
	}
}

func (s *server) ShortSingleURL(ctx context.Context, req *genproto.ShortSingleURLRequest) (*genproto.ShortSingleURLResponse, error) {
	if req.GetUrl() == "" {
		return nil, status.Errorf(codes.FailedPrecondition, "empty URL provided")
	}

	shortURL, err := s.urlShortener.HashFromURL(ctx, req.GetUrl())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &genproto.ShortSingleURLResponse{
		ShortUrl: fmt.Sprintf("%s/r/%s", s.baseDomain, shortURL.Hash),
		LongUrl:  req.GetUrl(),
	}, nil
}

func (s *server) BalanceURLs(ctx context.Context, req *genproto.BalanceURLsRequest) (*genproto.BalanceURLsResponse, error) {
	//fixme(fede): use the ctx for cancellation
	balancedURL, err := s.loadBalancer.ShortURLs(ctx, req.GetUrls())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &genproto.BalanceURLsResponse{ShortUrl: fmt.Sprintf("%s/lb/%s", s.baseDomain, balancedURL.Hash)}, nil
}

type Config struct {
	BaseDomain                 string
	ShortURLRepository         event.Repository
	CustomMetrics              url.Metrics
	LoadBalancedURLsRepository event.Repository
}

func NewServer(config Config) *grpc.Server {
	grpcServer := grpc.NewServer()
	srv := &server{
		baseDomain:   config.BaseDomain,
		urlShortener: url.NewSingleURLShortener(config.ShortURLRepository, clock.NewFromSystem(), config.CustomMetrics),
		loadBalancer: url.NewLoadBalancer(config.LoadBalancedURLsRepository, clock.NewFromSystem()),
	}

	genproto.RegisterURLShorteningServer(grpcServer, srv)

	reflection.Register(grpcServer)
	return grpcServer
}
