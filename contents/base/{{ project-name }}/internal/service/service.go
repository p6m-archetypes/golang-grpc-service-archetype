// Package service implements the {{ PrefixName }}{{ SuffixName }} gRPC service — the p6m
// standard CRUD surface over the {{ PrefixName }} entity. The generated bindings live in
// gen/{{ prefix_name }}_{{ suffix_name }} (rendered by `make proto`; the production image
// generates them during the build).
package service

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "{{ module_path }}/gen/{{ prefix_name }}_{{ suffix_name }}"
	"{{ module_path }}/internal/repository"
)

// {{ PrefixName }}{{ SuffixName }}Service implements pb.{{ PrefixName }}{{ SuffixName }}Server.
type {{ PrefixName }}{{ SuffixName }}Service struct {
	pb.Unimplemented{{ PrefixName }}{{ SuffixName }}Server
	store *repository.Store
}

// New returns the service backed by the given store.
func New(store *repository.Store) *{{ PrefixName }}{{ SuffixName }}Service {
	return &{{ PrefixName }}{{ SuffixName }}Service{store: store}
}

func (s *{{ PrefixName }}{{ SuffixName }}Service) Create{{ PrefixName }}(ctx context.Context, req *pb.Create{{ PrefixName }}Request) (*pb.{{ PrefixName }}, error) {
	e, err := s.store.Create(ctx, req.GetDisplayName())
	if err != nil {
		return nil, storeError(err)
	}
	return toProto(e), nil
}

func (s *{{ PrefixName }}{{ SuffixName }}Service) Get{{ PrefixName }}(ctx context.Context, req *pb.Get{{ PrefixName }}Request) (*pb.{{ PrefixName }}, error) {
	e, err := s.store.Get(ctx, req.GetId())
	if err != nil {
		return nil, storeError(err)
	}
	return toProto(e), nil
}

func (s *{{ PrefixName }}{{ SuffixName }}Service) List{{ PrefixName }}s(ctx context.Context, req *pb.List{{ PrefixName }}sRequest) (*pb.List{{ PrefixName }}sResponse, error) {
	items, err := s.store.List(ctx)
	if err != nil {
		return nil, storeError(err)
	}
	out := &pb.List{{ PrefixName }}sResponse{}
	for _, e := range items {
		out.Items = append(out.Items, toProto(e))
	}
	return out, nil
}

func (s *{{ PrefixName }}{{ SuffixName }}Service) Update{{ PrefixName }}(ctx context.Context, req *pb.Update{{ PrefixName }}Request) (*pb.{{ PrefixName }}, error) {
	e, err := s.store.Update(ctx, req.GetId(), req.GetDisplayName())
	if err != nil {
		return nil, storeError(err)
	}
	return toProto(e), nil
}

func (s *{{ PrefixName }}{{ SuffixName }}Service) Delete{{ PrefixName }}(ctx context.Context, req *pb.Delete{{ PrefixName }}Request) (*pb.Delete{{ PrefixName }}Response, error) {
	if err := s.store.Delete(ctx, req.GetId()); err != nil {
		return nil, storeError(err)
	}
	return &pb.Delete{{ PrefixName }}Response{}, nil
}

func toProto(e repository.{{ PrefixName }}) *pb.{{ PrefixName }} {
	return &pb.{{ PrefixName }}{Id: e.ID, DisplayName: e.DisplayName}
}

func storeError(err error) error {
	if errors.Is(err, repository.ErrNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}
	return status.Error(codes.Internal, err.Error())
}
