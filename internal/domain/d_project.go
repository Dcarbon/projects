package domain

import (
	"fmt"
	"time"

	"github.com/Dcarbon/arch-proto/pb"
	"github.com/Dcarbon/go-shared/dmodels"
)

type IProject interface {
	Create(req *RProjectCreate) (*Project, error)
	Update(req *RProjectUpdate) (*int64, error)
	UpdateDesc(req *RProjectUpdateDesc) (*ProjectDesc, error)
	UpdateSpecs(req *RProjectUpdateSpecs) (*ProjectSpecs, error)
	GetById(id int64, lang string) (*Project, error)
	GetList(filter *RProjectGetList) (*int64, []*Project, error)
	GetOwner(projectId int64) (string, error)
	AddImage(*RProjectAddImage) (*ProjectImage, error)
	ChangeStatus(id int, status ProjectStatus) error
	GetCountry(id string, vi string) (*Country, error)
	UpsertDocument(req *RProjectDocumentUpsert) ([]*ProjectDocument, error)
	ListDocument(req *RProjectDocumentList) ([]*ProjectDocument, int64, error)
	DeleteDocument(req *RProjectDocumentDelete) error
}

type RProjectCreate struct {
	Owner        dmodels.EthAddress    ``
	Location     *dmodels.Coord        ``
	Specs        *RProjectUpdateSpecs  ``
	Descs        []*RProjectUpdateDesc ``
	Area         float64               ``
	LocationName string                ``
	Type         int32                 ``
	Unit         float32               ``
	CountryId    string                ``
	OwnerId      string                ``
	Iframe       string                ``
	OwnerAddress string                ``
}

type RProjectUpdate struct {
	ProjectId    int64              ``
	CountryId    string             ``
	OwnerId      string             ``
	Type         int64              ``
	Unit         float32            ``
	Thumbnail    string             ``
	Owner        dmodels.EthAddress ``
	Location     *dmodels.Coord     ``
	LocationName string             ``
	Iframe       string             ``
	OwnerAddress string             ``
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
	Skip        int    `json:"skip" form:"skip"`
	Limit       int    `json:"limit" form:"limit;max=50"`
	Owner       string `json:"owner" form:"owner"`
	Unit        int64  ``
	Type        int64  ``
	CountryId   string ``
	SearchValue string ``
	Location    string ``
	Status      int    ``
	Ids         []int  ``
}

type RProjectAddImage struct {
	ProjectId int64  `json:"projectId"`
	ImgPath   string `json:"imgPath"`
	Type      int32
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
		OwnerId:      rproject.OwnerId,
		CountryId:    rproject.CountryId,
		Type:         int64(rproject.Type),
		Unit:         rproject.Unit,
		Iframe:       rproject.Iframe,
		Area:         rproject.Area,
		OwnerAddress: rproject.OwnerAddress,
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

func (p RProjectGetList) GetUnit() string {
	var query string
	var ranges [][2]int

	switch p.Type {
	case int64(*pb.ProjectType_PrjT_G.Enum()):
		ranges = [][2]int{{0, 20}, {20, 100}, {100, -1}}
	case int64(*pb.ProjectType_PrjT_E.Enum()), int64(*pb.ProjectType_PrjT_S.Enum()):
		ranges = [][2]int{{0, 90}, {90, 200}, {200, -1}}
	}

	if len(ranges) > 0 {
		switch p.Unit {
		case 1:
			query = fmt.Sprintf("unit >= %d AND unit < %d", ranges[0][0], ranges[0][1])
		case 2:
			query = fmt.Sprintf("unit >= %d AND unit <= %d", ranges[1][0], ranges[1][1])
		case 3:
			query = fmt.Sprintf("unit > %d", ranges[2][0])
		}
	}

	return query
}

type Document struct {
	Url          string
	DocumentName string
	DocumentType string
	ProjectId    int64
	Id           int64
}

type RProjectDocumentUpsert struct {
	Document []*Document ``
}

type RProjectDocumentDelete struct {
	Id []int64
}

type RProjectDocumentList struct {
	Skip  int     `json:"skip" form:"skip"`
	Limit int     `json:"limit" form:"limit;max=50"`
	Ids   []int64 ``
}
