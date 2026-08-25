// Package service implements the {{ ProjectName }} gRPC service — the p6m
// standard CRUD surface over the {{ EntityName }} entity. The generated bindings live in
// gen/{{ project_name }} (rendered by `make proto`; the production image
// generates them during the build).
package service

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "{{ module_path }}/gen/{{ project_name }}"
	"{{ module_path }}/internal/repository"
)

// {{ ProjectName }}Service implements pb.{{ ProjectName }}Server.
type {{ ProjectName }}Service struct {
	pb.Unimplemented{{ ProjectName }}Server
	store *repository.Store
}

// New returns the service backed by the given store.
func New(store *repository.Store) *{{ ProjectName }}Service {
	return &{{ ProjectName }}Service{store: store}
}

func (s *{{ ProjectName }}Service) Create{{ EntityName }}(ctx context.Context, req *pb.Create{{ EntityName }}Request) (*pb.{{ EntityName }}, error) {
	e, err := s.store.Create(ctx, req.GetDisplayName())
	if err != nil {
		return nil, storeError(err)
	}
	return toProto(e), nil
}

func (s *{{ ProjectName }}Service) Get{{ EntityName }}(ctx context.Context, req *pb.Get{{ EntityName }}Request) (*pb.{{ EntityName }}, error) {
	e, err := s.store.Get(ctx, req.GetId())
	if err != nil {
		return nil, storeError(err)
	}
	return toProto(e), nil
}

func (s *{{ ProjectName }}Service) List{{ EntityName }}s(ctx context.Context, req *pb.List{{ EntityName }}sRequest) (*pb.List{{ EntityName }}sResponse, error) {
	items, err := s.store.List(ctx)
	if err != nil {
		return nil, storeError(err)
	}
	out := &pb.List{{ EntityName }}sResponse{}
	for _, e := range items {
		out.Items = append(out.Items, toProto(e))
	}
	return out, nil
}

func (s *{{ ProjectName }}Service) Update{{ EntityName }}(ctx context.Context, req *pb.Update{{ EntityName }}Request) (*pb.{{ EntityName }}, error) {
	e, err := s.store.Update(ctx, req.GetId(), req.GetDisplayName())
	if err != nil {
		return nil, storeError(err)
	}
	return toProto(e), nil
}

func (s *{{ ProjectName }}Service) Delete{{ EntityName }}(ctx context.Context, req *pb.Delete{{ EntityName }}Request) (*pb.Delete{{ EntityName }}Response, error) {
	if err := s.store.Delete(ctx, req.GetId()); err != nil {
		return nil, storeError(err)
	}
	return &pb.Delete{{ EntityName }}Response{}, nil
}

func toProto(e repository.{{ EntityName }}) *pb.{{ EntityName }} {
	return &pb.{{ EntityName }}{Id: e.ID, DisplayName: e.DisplayName}
}

func storeError(err error) error {
	if errors.Is(err, repository.ErrNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}
	return status.Error(codes.Internal, err.Error())
}
