package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	v1 "github.com/aaabhilash97/op/pkg/api/v1"
	"github.com/aaabhilash97/op/pkg/config"
	"github.com/aaabhilash97/op/pkg/db"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

// toDoServiceServer is implementation of v1.ToDoServiceServer proto interface
type OpServiceServer struct {
	config config.Config
	db     *db.DB
}

type Options struct {
	Config config.Config
	DB     *db.DB
}

// NewOpServiceServer creates ToDo service
func NewOpServiceServer(opt Options) v1.OpServiceServer {
	return &OpServiceServer{
		config: opt.Config,
		db:     opt.DB,
	}
}

// checkAPI checks if the API version requested by client is supported by server
func (s *OpServiceServer) checkAPI(api string) error {
	// API version is "" means use current version of the service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}
	return nil
}

// Create new todo task
func (s *OpServiceServer) NewUserRegistration(ctx context.Context, req *v1.NewUserRegistrationRequest) (*v1.NewUserRegistrationResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(""); err != nil {
		return nil, err
	}

	return &v1.NewUserRegistrationResponse{}, nil
}
