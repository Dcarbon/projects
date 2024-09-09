package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Dcarbon/projects/internal/domain/repo"

	"github.com/Dcarbon/arch-proto/pb"
	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/go-shared/gutils"
	"github.com/Dcarbon/go-shared/libs/sclient"
	"github.com/Dcarbon/projects/internal/domain"
	"github.com/Dcarbon/projects/internal/rss"
)

type Service struct {
	pb.UnimplementedProjectServiceServer
	*gutils.GService
	iProject domain.IProject
	storage  sclient.IStorage
}

func NewProjectService(config *gutils.Config,
) (*Service, error) {
	rss.SetUrl(config.GetDBUrl())

	iProject, err := repo.NewProjectImpl(rss.GetDB())
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
	var sv = &Service{
		iProject: iProject,
		storage:  storage,
	}

	return sv, nil
}

func (sv *Service) Create(ctx context.Context, req *pb.RPCreate,
) (*pb.Project, error) {
	var descs []*domain.RProjectUpdateDesc
	for _, desc := range req.Descs {
		descs = append(descs, &domain.RProjectUpdateDesc{
			Language: desc.Language,
			Name:     desc.Name,
			Desc:     desc.Desc})
	}
	project, err := sv.iProject.Create(&domain.RProjectCreate{
		Owner:        dmodels.EthAddress(req.Owner),
		Location:     dmodels.NewCoord4326(req.Location.Longitude, req.Location.Latitude),
		Specs:        &domain.RProjectUpdateSpecs{Specs: req.Specs.GetSpecs()},
		Descs:        descs,
		Area:         0, //TODO: fix Area
		LocationName: req.LocationName,
		Type:         int32(req.Type),
		Unit:         float32(req.Unit),
		CountryId:    req.CountryId,
		OwnerId:      req.OwnerId,
		Iframe:       req.Iframe,
		OwnerAddress: req.OwnerAddress,
	})
	if err != nil {
		return nil, err
	}
	return convertProject(project), nil
}

func (sv *Service) UpdateDesc(ctx context.Context, req *pb.RPUpdateDesc,
) (*pb.ProjectDesc, error) {
	desc, err := sv.iProject.UpdateDesc(&domain.RProjectUpdateDesc{
		ProjectId: req.ProjectId,
		Language:  req.Language,
		Name:      req.Name,
		Desc:      req.Desc,
	})
	if err != nil {
		return nil, err
	}
	return convertProjectDesc(desc), nil
}

func (sv *Service) UpdateSpecs(ctx context.Context, req *pb.RPUpdateSpecs,
) (*pb.ProjectSpecs, error) {
	spec, err := sv.iProject.UpdateSpecs(&domain.RProjectUpdateSpecs{
		ProjectId: req.ProjectId,
		Specs:     req.Specs,
	})
	if err != nil {
		return nil, err
	}
	return convertProjectSpecs(spec), nil
}

func (sv *Service) AddImage(ctx context.Context, req *pb.RPAddImage,
) (*pb.String, error) {
	image, err := sv.iProject.AddImage(&domain.RProjectAddImage{ProjectId: req.ProjectId, ImgPath: req.Image, Type: req.Type})
	if err != nil {
		return nil, err
	}
	return &pb.String{Data: image.Image}, nil
}

func (sv *Service) GetById(ctx context.Context, req *pb.RPGetById,
) (*pb.Project, error) {
	data, err := sv.iProject.GetById(req.ProjectId, req.Lang)
	if nil != err {
		return nil, err
	}
	response := convertProject(data)
	response.Address = data.LocationName
	return response, nil
}

func (sv *Service) ChangeStatus(ctx context.Context, req *pb.RPChangeStatus,
) (*pb.Int64, error) {
	if err := sv.iProject.ChangeStatus(int(req.ProjectId), domain.ProjectStatus(req.Status)); nil != err {
		return nil, err
	}
	return &pb.Int64{Data: req.ProjectId}, nil
}

func (sv *Service) GetList(ctx context.Context, req *pb.RPGetList,
) (*pb.Projects, error) {
	intArray := []int{}
	if strings.TrimSpace(req.Ids) != "" {
		datas := strings.Split(req.Ids, ",")
		// Convert each string to an integer
		for _, s := range datas {
			// Convert string to integer
			num, err := strconv.Atoi(strings.TrimSpace(s))
			if err != nil {
				fmt.Printf("Error converting string to int: %v\n", err)
				return nil, err
			}
			intArray = append(intArray, num) // Append the converted integer to intArray
		}
	}
	count, data, err := sv.iProject.GetList(&domain.RProjectGetList{
		Skip:        int(req.Skip),
		Limit:       int(req.Limit),
		Owner:       req.OwnerId,
		Unit:        int64(req.Unit),
		CountryId:   req.CountryId,
		Type:        int64(req.Type),
		SearchValue: req.SearchValue,
		Location:    req.Location,
		Status:      int(req.Status),
		Ids:         intArray,
	})
	if nil != err {
		return nil, err
	}
	return &pb.Projects{
		Total: *count,
		Data:  convertArr[domain.Project, pb.Project](data, convertProject),
	}, nil
}

func (sv *Service) Update(ctx context.Context, req *pb.RPUpdate,
) (*pb.Int64, error) {
	id, err := sv.iProject.Update(&domain.RProjectUpdate{
		ProjectId:    req.ProjectId,
		Owner:        dmodels.EthAddress(req.Owner),
		Location:     dmodels.NewCoord4326(req.Location.Longitude, req.Location.Latitude),
		LocationName: req.LocationName,
		Type:         int64(req.Type),
		Unit:         float32(req.Unit),
		CountryId:    req.CountryId,
		OwnerId:      req.OwnerId,
		Iframe:       req.Iframe,
		OwnerAddress: req.OwnerAddress,
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &pb.Int64{Data: *id}, nil
}

func (sv *Service) UpsertDocument(ctx context.Context, req *pb.RUpsertDocument) (*pb.RPUpsertDocument, error) {
	documents := []*domain.Document{}
	for _, val := range req.Documents {
		documents = append(documents, &domain.Document{
			Url:          val.Url,
			ProjectId:    val.ProjectId,
			DocumentType: val.DocumentType,
			Id:           val.Id,
		})
	}
	data, err := sv.iProject.UpsertDocument(&domain.RProjectDocumentUpsert{Document: documents})
	if err != nil {
		return nil, err
	}

	return &pb.RPUpsertDocument{
		Documents: convertArr(data, convertDocument),
	}, nil
}

func (sv *Service) DeleteDocument(ctx context.Context, req *pb.RDeleteDocument) (*pb.Empty, error) {
	err := sv.iProject.DeleteDocument(&domain.RProjectDocumentDelete{Id: []int64{req.Id}})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
func (sv *Service) ListDocument(ctx context.Context, req *pb.RListDocument) (*pb.RPListDocument, error) {
	data, count, err := sv.iProject.ListDocument(&domain.RProjectDocumentList{
		Skip:  int(req.Skip),
		Limit: int(req.Limit),
		Ids:   req.Ids,
	})
	if err != nil {
		return nil, err
	}
	return &pb.RPListDocument{
		Documents: convertArr(data, convertDocument),
		Total:     count,
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
