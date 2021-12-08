package grpc

import (
	"fmt"
	"io"

	genproto "github.com/WebEngineeringGroupI/genproto-go/api/v1alpha1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

type server struct {
	genproto.UnimplementedURLShorteningServer
	urlShortener *url.SingleURLShortener
	baseDomain   string
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

		shortURL, err := s.urlShortener.HashFromURL(request.Url)
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

type Config struct {
	BaseDomain         string
	ShortURLRepository url.ShortURLRepository
	CustomMetrics      url.Metrics
}

func NewServer(config Config) *grpc.Server {
	grpcServer := grpc.NewServer()
	srv := &server{
		baseDomain:   config.BaseDomain,
		urlShortener: url.NewSingleURLShortener(config.ShortURLRepository, config.CustomMetrics),
	}

	genproto.RegisterURLShorteningServer(grpcServer, srv)

	reflection.Register(grpcServer)
	return grpcServer
}
