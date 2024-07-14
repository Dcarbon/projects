package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/Dcarbon/go-shared/dmodels"
)

const (
	TableNameProject      = "projects"
	TableNameProjectDesc  = "projects_desc"
	TableNameProjectSpecs = "projects_specs"
	TableNameProjectImage = "projects_image"
)

type ProjectStatus int

const (
	ProjectStatusReject   ProjectStatus = -1
	ProjectStatusRegister ProjectStatus = 1
	ProjectStatusActived  ProjectStatus = 20
)

type Project struct {
	Id           int64              `json:"id"                        gorm:"primaryKey"`                 //
	Owner        dmodels.EthAddress `json:"owner"                     gorm:"index"`                      // ETH address
	Status       ProjectStatus      `json:"status"                    `                                  //
	LocationName string             `json:"locationName,omitempty"    `                                  //
	Location     *dmodels.Coord     `json:"location"                  gorm:"type:geometry(POINT, 4326)"` //
	Specs        *ProjectSpecs      `json:"specs,omitempty"           gorm:"foreignKey:ProjectId"`       //
	Area         float64            `json:"area,omitempty"            `                                  //
	Descs        []*ProjectDesc     `json:"descs,omitempty"           gorm:"foreignKey:ProjectId"`       //
	Images       []*ProjectImage    `json:"images,omitempty"          gorm:"foreignKey:ProjectId"`       //
	Thumbnail    string             `json:"thumbnail"`
	CreatedAt    time.Time          `json:"createdAt"                 ` //
	UpdatedAt    time.Time          `json:"updatedAt"                 ` //
	Type         int64              `json:"type" gorm:"column:type"`
	Unit         float32            `json:"unit" gorm:"column:unit"`
	CountryId    int64              `gorm:"column:country_id"`
	Country      *Country           `json:"country" gorm:"-"`
} //@name Project

func (*Project) TableName() string { return TableNameProject }

type ProjectDesc struct {
	Id        int64     `gorm:"primaryKey"`
	ProjectId int64     `gorm:"index:idx_project_desc_lang,unique,priority:1"` //
	Language  string    `gorm:"index:idx_project_desc_lang,unique,priority:2"` //
	Name      string    ``
	Desc      string    ``
	CreatedAt time.Time ``
	UpdatedAt time.Time ``
} //@name ProjectDescription

func (*ProjectDesc) TableName() string { return TableNameProjectDesc }

type ProjectSpecs struct {
	Id        int64     `json:"id"              gorm:"primaryKey"`
	ProjectId int64     `json:"projectId"       gorm:"unique"`
	Specs     MapSFloat `json:"specs"           gorm:"type:json"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
} //@name ProjectSpec

func (*ProjectSpecs) TableName() string { return TableNameProjectSpecs }

type ProjectImage struct {
	Id        int64     `json:"id"`        //
	ProjectId int64     `json:"projectId"` //
	Image     string    `json:"image"`     // Image path
	CreatedAt time.Time `json:"createdAt"`
}

func (*ProjectImage) TableName() string { return TableNameProjectImage }

type MapSFloat map[string]float64 //@name MapSFloat

func (m *MapSFloat) Scan(value interface{}) error {
	if nil == m {
		m = new(MapSFloat)
	}
	switch vt := value.(type) {
	case string:
		return json.Unmarshal([]byte(vt), m)
	case []byte:
		return json.Unmarshal(vt, m)
	}
	return errors.New("scan value type for MapSFloat invalid")
}

type Country struct {
	Id   int64  `json:"id"  `
	Name string `json:"name"`
}
type Language struct {
	Locale string `json:"locale"`
	Name   string `json:"name"`
}

func (m MapSFloat) Value() (driver.Value, error) {
	if nil == m {
		return nil, nil
	}
	return json.Marshal(m)
}
