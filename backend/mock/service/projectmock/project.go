package projectmock

import (
	"context"

	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"

	projectv1core "github.com/lyft/clutch/backend/api/core/project/v1"
	projectv1 "github.com/lyft/clutch/backend/api/project/v1"
	"github.com/lyft/clutch/backend/service"
	projectservice "github.com/lyft/clutch/backend/service/project"
)

type svc struct{}

func New() projectservice.Service {
	return &svc{}
}

func NewAsService(*anypb.Any, *zap.Logger, tally.Scope) (service.Service, error) {
	return New(), nil
}

func (s *svc) GetProjects(ctx context.Context, req *projectv1.GetProjectsRequest) (*projectv1.GetProjectsResponse, error) {
	// scenario: we call the api with a user's email
	if len(req.Users) != 0 && len(req.Projects) == 0 {
		return &projectv1.GetProjectsResponse{
			Results: map[string]*projectv1.ProjectResult{
				"blog": {
					From: &projectv1.ProjectResult_From{
						Selected: true,
						Users:    []string{"cat@meow.com"},
					},
					Project: &projectv1core.Project{
						Name: "blog",
						Dependencies: &projectv1core.ProjectDependencies{
							Upstreams: map[string]*projectv1core.Dependency{},
							Downstreams: map[string]*projectv1core.Dependency{
								"type.googleapis.com/clutch.core.project.v1.Project": {
									Ids: []string{"posts"},
								},
							},
						},
					},
				},
				"posts": {
					From: &projectv1.ProjectResult_From{
						Selected: false,
						Users:    []string{},
					},
					Project: &projectv1core.Project{
						Name: "posts",
						Dependencies: &projectv1core.ProjectDependencies{
							Upstreams: map[string]*projectv1core.Dependency{
								"type.googleapis.com/clutch.core.project.v1.Project": {
									Ids: []string{"blog"},
								},
							},
							Downstreams: map[string]*projectv1core.Dependency{},
						},
					},
				},
			},
		}, nil
	}

	// scenario: we call the api with a user's email and custom project "comments"
	if len(req.Users) != 0 {
		for _, p := range req.Projects {
			if p == "comments" {
				return &projectv1.GetProjectsResponse{
					Results: map[string]*projectv1.ProjectResult{
						"blog": {
							From: &projectv1.ProjectResult_From{
								Selected: true,
								Users:    []string{"cat@meow.com"},
							},
							Project: &projectv1core.Project{
								Name: "blog",
								Dependencies: &projectv1core.ProjectDependencies{
									Upstreams: map[string]*projectv1core.Dependency{},
									Downstreams: map[string]*projectv1core.Dependency{
										"type.googleapis.com/clutch.core.project.v1.Project": {
											Ids: []string{"posts"},
										},
									},
								},
							},
						},
						"posts": {
							From: &projectv1.ProjectResult_From{
								Selected: false,
								Users:    []string{},
							},
							Project: &projectv1core.Project{
								Name: "posts",
								Dependencies: &projectv1core.ProjectDependencies{
									Upstreams: map[string]*projectv1core.Dependency{
										"type.googleapis.com/clutch.core.project.v1.Project": {
											Ids: []string{"blog"},
										},
									},
									Downstreams: map[string]*projectv1core.Dependency{},
								},
							},
						},
						"comments": {
							From: &projectv1.ProjectResult_From{
								Selected: true,
								Users:    []string{},
							},
							Project: &projectv1core.Project{
								Name: "comments",
								Dependencies: &projectv1core.ProjectDependencies{
									Upstreams: map[string]*projectv1core.Dependency{
										"type.googleapis.com/clutch.core.project.v1.Project": {
											Ids: []string{"users"},
										},
									},
									Downstreams: map[string]*projectv1core.Dependency{
										"type.googleapis.com/clutch.core.project.v1.Project": {
											Ids: []string{"likes"},
										},
									},
								},
							},
						},
						"users": {
							From: &projectv1.ProjectResult_From{
								Selected: false,
								Users:    []string{},
							},
							Project: &projectv1core.Project{
								Name: "users",
								Dependencies: &projectv1core.ProjectDependencies{
									Upstreams: map[string]*projectv1core.Dependency{},
									Downstreams: map[string]*projectv1core.Dependency{
										"type.googleapis.com/clutch.core.project.v1.Project": {
											Ids: []string{"comments"},
										},
									},
								},
							},
						},
						"likes": {
							From: &projectv1.ProjectResult_From{
								Selected: false,
								Users:    []string{},
							},
							Project: &projectv1core.Project{
								Name: "likes",
								Dependencies: &projectv1core.ProjectDependencies{
									Upstreams: map[string]*projectv1core.Dependency{
										"type.googleapis.com/clutch.core.project.v1.Project": {
											Ids: []string{"comments"},
										},
									},
									Downstreams: map[string]*projectv1core.Dependency{},
								},
							},
						},
					},
				}, nil
			}
		}
	}

	return &projectv1.GetProjectsResponse{
		Results: map[string]*projectv1.ProjectResult{
			"clutch": {
				From: &projectv1.ProjectResult_From{
					Selected: true,
					Users:    []string{"cat@meow.com"},
				},
				Project: &projectv1core.Project{
					Name: "clutch",
					Dependencies: &projectv1core.ProjectDependencies{
						Upstreams:   map[string]*projectv1core.Dependency{},
						Downstreams: map[string]*projectv1core.Dependency{},
					},
				},
			},
			"choiceapi": {
				From: &projectv1.ProjectResult_From{
					Selected: true,
					Users:    []string{"nom@meow.com"},
				},
				Project: &projectv1core.Project{
					Name: "choiceapi",
					Dependencies: &projectv1core.ProjectDependencies{
						Upstreams:   map[string]*projectv1core.Dependency{},
						Downstreams: map[string]*projectv1core.Dependency{},
					},
				},
			},
		},
	}, nil
}
