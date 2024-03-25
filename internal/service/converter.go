package service

import (
	"github.com/Dcarbon/arch-proto/pb"
	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/projects/internal/domain"
)

func convertProject(in *domain.Project) *pb.Project {
	if nil == in {
		return nil
	}
	var rs = &pb.Project{
		Id:           in.Id,
		Owner:        in.Owner.String(),
		LocationName: in.LocationName,
		Location:     convertGPS(in.Location),
		Status:       int32(in.Status),
		Ca:           in.CreatedAt.UnixMilli(),
		Ua:           in.UpdatedAt.UnixMilli(),
		Images:       []string{},
		Specs:        convertProjectSpecs(in.Specs),
		Descs:        convertArr[domain.ProjectDesc, pb.ProjectDesc](in.Descs, convertProjectDesc),
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

func convertArr[T any, T2 any](arr []*T, fn func(*T) *T2) []*T2 {
	var rs = make([]*T2, len(arr))
	for i, it := range arr {
		rs[i] = fn(it)
	}
	return rs
}
