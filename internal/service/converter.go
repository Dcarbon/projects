package service

import (
	"github.com/Dcarbon/arch-proto/pb"
	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/go-shared/libs/utils"
	"github.com/Dcarbon/projects/internal/domain"
)

func convertProject(in *domain.Project) *pb.Project {
	types := map[int]string{0: "None", 1: "Biomass to Gasification", 2: "Biogas to Electricity", 3: "Model S"}

	if nil == in {
		return nil
	}
	var rs = &pb.Project{
		Thumbnail:    in.Thumbnail,
		Id:           in.Id,
		Owner:        in.OwnerId,
		LocationName: in.LocationName,
		Location:     convertGPS(in.Location),
		Status:       int32(in.Status),
		Ca:           in.CreatedAt.UnixMilli(),
		Ua:           in.UpdatedAt.UnixMilli(),
		Images:       convertImage(in.Images),
		Specs:        convertProjectSpecs(in.Specs),
		Descs:        convertArr[domain.ProjectDesc, pb.ProjectDesc](in.Descs, convertProjectDesc),
		Area:         in.Area,
		Type:         pb.ProjectType(in.Type),
		Unit:         in.Unit,
		Country:      convertCountry(in.Country),
		Iframe:       in.Iframe,
		DetailType: &pb.Type{
			Id:   int32(in.Type),
			Name: types[int(in.Type)],
		},
	}
	return rs
}

func convertCountry(in *domain.Country) *pb.Country {
	if nil == in {
		return nil
	}
	var rs = &pb.Country{
		//Id:          in.Id,
		Name:        in.Name,
		CountryCode: in.Code,
	}
	return rs
}

func convertProjectDesc(in *domain.ProjectDesc) *pb.ProjectDesc {
	if nil == in {
		return nil
	}
	var rs = &pb.ProjectDesc{
		Id:       in.Id,
		Language: in.Language,
		Name:     in.Name,
		Desc:     in.Desc,
	}
	return rs
}

func convertProjectSpecs(in *domain.ProjectSpecs) *pb.ProjectSpecs {
	if nil == in {
		return nil
	}
	var rs = &pb.ProjectSpecs{
		Id:    in.Id,
		Specs: in.Specs,
	}
	return rs
}

func convertGPS(in *dmodels.Coord) *pb.GPS {
	if nil == in {
		return nil
	}
	var rs = &pb.GPS{
		Latitude:  in.Lat,
		Longitude: in.Lng,
	}
	return rs
}

func convertImage(in []*domain.ProjectImage) []string {
	if nil == in {
		return nil
	}
	var rs = make([]string, len(in))
	for i, it := range in {
		rs[i] = utils.StringEnv("STORAGE_URL", "") + it.Image
	}
	return rs
}

func convertArr[T any, T2 any](arr []*T, fn func(*T) *T2) []*T2 {
	var rs = make([]*T2, len(arr))
	for i, it := range arr {
		rs[i] = fn(it)
	}
	return rs
}
