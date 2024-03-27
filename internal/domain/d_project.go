package domain

import (
	"time"

	"github.com/Dcarbon/go-shared/dmodels"
)

type IProject interface {
	Create(req *RProjectCreate) (*Project, error)

	UpdateDesc(req *RProjectUpdateDesc) (*ProjectDesc, error)
	UpdateSpecs(req *RProjectUpdateSpecs) (*ProjectSpecs, error)

	GetById(id int64, lang string) (*Project, error)
	GetList(filter *RProjectGetList) ([]*Project, error)
	GetOwner(projectId int64) (string, error)

	AddImage(*RProjectAddImage) (*ProjectImage, error)
	ChangeStatus(id string, status ProjectStatus) error
}

type RProjectCreate struct {
	Owner        dmodels.EthAddress    ``
	Location     *dmodels.Coord        ``
	Specs        *RProjectUpdateSpecs  ``
	Descs        []*RProjectUpdateDesc ``
	Area         float64               ``
	LocationName string                ``
}

type RProjectUpdateDesc struct {
	ProjectId int64  ``
	Language  string ``
	Name      string ``
	Desc      string ``
}

type RProjectUpdateSpecs struct {
	ProjectId int64              `json:"projectId"`
	Specs     map[string]float64 `json:"specs"`
}

type RProjectGetList struct {
	Skip  int    `json:"skip" form:"skip"`
	Limit int    `json:"limit" form:"limit;max=50"`
	Owner string `json:"owner" form:"owner"`
}

type RProjectAddImage struct {
	ProjectId int64  `json:"projectId"`
	ImgPath   string `json:"imgPath"`
}

func (rproject *RProjectCreate) ToProject() *Project {
	var project = &Project{
		Id:           0,
		LocationName: rproject.LocationName,
		Status:       ProjectStatusRegister,
		Owner:        rproject.Owner,
		Location:     rproject.Location,
		Specs:        rproject.Specs.ToProjectSpecs(),
		Descs:        make([]*ProjectDesc, len(rproject.Descs)),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	for i, desc := range rproject.Descs {
		project.Descs[i] = desc.ToProjectDesc()
	}

	return project
}

func (rdesc *RProjectUpdateDesc) ToProjectDesc() *ProjectDesc {
	if rdesc == nil {
		return nil
	}
	return &ProjectDesc{
		Id:        0,
		ProjectId: rdesc.ProjectId,
		Language:  rdesc.Language,
		Name:      rdesc.Name,
		Desc:      rdesc.Desc,
	}
}

func (rspec *RProjectUpdateSpecs) ToProjectSpecs() *ProjectSpecs {
	if rspec.Specs == nil {
		return nil
	}
	return &ProjectSpecs{
		Id:        0,
		ProjectId: rspec.ProjectId,
		Specs:     rspec.Specs,
	}
}
