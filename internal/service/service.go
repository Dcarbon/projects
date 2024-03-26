package service

import (
	"context"

	"github.com/Dcarbon/arch-proto/pb"
	"github.com/Dcarbon/go-shared/gutils"
	"github.com/Dcarbon/go-shared/libs/sclient"
	"github.com/Dcarbon/projects/internal/domain"
	"github.com/Dcarbon/projects/internal/domain/infra"
	"github.com/Dcarbon/projects/internal/rss"
)

type Service struct {
	pb.UnimplementedProjectServiceServer
	*gutils.GService
	iProject   domain.IProject
	storage    sclient.IStorage
	mapService *MapService
}

func NewProjectService(config *gutils.Config,
) (*Service, error) {
	rss.SetUrl(config.GetDBUrl())

	iProject, err := infra.NewProjectImpl(rss.GetDB())
	if nil != err {
		return nil, err
	}

	gserice, err := gutils.NewGService(config, "")
	if nil != err {
		return nil, err
	}

	storage, err := sclient.NewStorage(
		config.GetStorageHost(), gserice.GetToken(),
	)
	if nil != err {
		return nil, err
	}

	mapService, err := NewMapService(config)
	if err != nil {
		return nil, err
	}
	var sv = &Service{
		iProject:   iProject,
		storage:    storage,
		mapService: mapService,
	}

	return sv, nil
}

func (sv *Service) Create(ctx context.Context, req *pb.RPCreate,
) (*pb.Project, error) {
	return nil, gutils.ErrorNotImplement
}

func (sv *Service) UpdateDesc(ctx context.Context, req *pb.RPUpdateDesc,
) (*pb.ProjectDesc, error) {
	return nil, nil
}

func (sv *Service) UpdateSpecs(ctx context.Context, req *pb.RPUpdateSpecs,
) (*pb.ProjectSpecs, error) {
	return nil, gutils.ErrorNotImplement
}

func (sv *Service) AddImage(ctx context.Context, req *pb.RPAddImage,
) (*pb.String, error) {
	return nil, gutils.ErrorNotImplement
}

func (sv *Service) GetById(ctx context.Context, req *pb.RPGetById,
) (*pb.Project, error) {
	data, err := sv.iProject.GetById(req.ProjectId, req.Lang)
	if nil != err {
		return nil, err
	}
	response := convertProject(data)
	address, _ := sv.mapService.GetAddress(data.Location.Lat, data.Location.Lng)
	response.Address = address
	return response, nil
}

func (sv *Service) GetList(ctx context.Context, req *pb.RPGetList,
) (*pb.Projects, error) {
	data, err := sv.iProject.GetList(&domain.RProjectGetList{
		Skip:  int(req.Skip),
		Limit: int(req.Limit),
		Owner: req.Owner,
	})
	if nil != err {
		return nil, err
	}
	return &pb.Projects{
		Data: convertArr[domain.Project, pb.Project](data, convertProject),
	}, nil
}

// func (sv *Service) isProjectOwner(ctx context.Context, projectId int64,
// ) error {
// 	user, err := mids.GetAuth(r.Request.Context())
// 	if nil != err {
// 		return dmodels.ErrInternal(errors.New("missing check authen in project add image"))
// 	}
// 	if user.Role == "super-admin" {
// 		return nil
// 	}
// 	owner, err := ctrl.repo.GetOwner(projectId)
// 	if nil != err {
// 		return err
// 	}
// 	if dmodels.EthAddress(user.EthAddress) != dmodels.EthAddress(owner) {
// 		return dmodels.ErrorPermissionDenied
// 	}
// 	return nil
// }
