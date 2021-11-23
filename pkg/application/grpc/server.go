package grpc

import (
	"context"
	"fmt"

	genproto "github.com/WebEngineeringGroupI/genproto-go/api/v1alpha1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

type server struct {
	genproto.UnimplementedURLShorteningServer
	urlShortener *url.SingleURLShortener
	baseDomain   string
}

func (s *server) ShortURLs(ctx context.Context, request *genproto.ShortURLsRequest) (*genproto.ShortURLsResponse, error) {
	results := []*genproto.ShortURLsResponse_Result{}

	for _, longURL := range request.GetUrls() {
		shortURL, err := s.urlShortener.HashFromURL(longURL)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		results = append(results, &genproto.ShortURLsResponse_Result{
			ShortUrl: fmt.Sprintf("%s/r/%s", s.baseDomain, shortURL.Hash),
			LongUrl:  shortURL.LongURL,
		})
	}

	return &genproto.ShortURLsResponse{
		Results: results,
	}, nil
}

type Config struct {
	BaseDomain         string
	ShortURLRepository url.ShortURLRepository
	URLValidator       url.Validator
}

func NewServer(config Config) *grpc.Server {
	grpcServer := grpc.NewServer()
	srv := &server{
		baseDomain:   config.BaseDomain,
		urlShortener: url.NewSingleURLShortener(config.ShortURLRepository, config.URLValidator),
	}
	genproto.RegisterURLShorteningServer(grpcServer, srv)

	reflection.Register(grpcServer)
	return grpcServer
}
