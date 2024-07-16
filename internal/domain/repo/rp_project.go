package repo

import (
	"encoding/json"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/projects/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProjectImpl struct {
	db *gorm.DB
}

func NewProjectImpl(db *gorm.DB) (*ProjectImpl, error) {
	err := db.AutoMigrate(
		&domain.Project{},
		&domain.ProjectImage{},
		&domain.ProjectSpecs{},
		&domain.ProjectDesc{},
	)
	if nil != err {
		return nil, err
	}

	var pp = &ProjectImpl{
		db: db,
	}
	return pp, nil
}

func (pImpl *ProjectImpl) Create(req *domain.RProjectCreate,
) (*domain.Project, error) {
	project := req.ToProject()
	if err := pImpl.tblProject().Transaction(func(dbTx *gorm.DB) error {
		if err := dbTx.Table(domain.TableNameProject).Omit("Images").
			Create(project).Error; err != nil {
			return dmodels.ParsePostgresError("Create project", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return project, nil
}

func (pImpl *ProjectImpl) UpdateDesc(req *domain.RProjectUpdateDesc,
) (*domain.ProjectDesc, error) {
	desc := req.ToProjectDesc()
	if err := pImpl.tblProjectDesc().
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{
					{Name: "project_id"}, {Name: "language"},
				},
				UpdateAll: true,
			}).
		Create(desc).Error; nil != err {
		return nil, dmodels.ParsePostgresError("Update project desc", err)
	}
	return desc, nil
}

func (pImpl *ProjectImpl) UpdateSpecs(req *domain.RProjectUpdateSpecs,
) (*domain.ProjectSpecs, error) {
	var spec = req.ToProjectSpecs()

	if err := pImpl.tblProjectSpec().
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{{Name: "project_id"}},
				DoUpdates: clause.AssignmentColumns(
					[]string{"specs", "updated_at"},
				),
			},
		).Create(spec).Error; nil != err {
		return nil, dmodels.ParsePostgresError("Update project desc", err)
	}
	return spec, nil
}

func (pImpl *ProjectImpl) GetById(id int64, lang string,
) (*domain.Project, error) {
	var project = &domain.Project{}
	var query = pImpl.tblProject().Where("id = ?", id).
		Preload("Images", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("project_id, image")
		}).
		Preload("Specs")
	if lang == "" {
		lang = "vi"
	}
	query.Preload("Descs", func(tx *gorm.DB) *gorm.DB {
		return tx.Where("language = ?", lang)
	})
	var err = query.First(project).Error
	if nil != err {
		return nil, dmodels.ParsePostgresError("Project", err)
	}
	country, _ := pImpl.GetCountry(int(project.CountryId), lang)
	project.Country = country
	return project, nil
}

func (pImpl *ProjectImpl) GetList(filter *domain.RProjectGetList,
) ([]*domain.Project, error) {
	var tbl = pImpl.tblProject().Offset(filter.Skip)
	if filter.Limit > 0 {
		tbl = tbl.Limit(filter.Limit)
	}

	if filter.Owner != "" {
		tbl = tbl.Where("owner = ?", filter.Owner)
	}

	var data = make([]*domain.Project, 0)
	var err = tbl.Preload("Descs").Find(&data).Error
	if nil != err {
		return nil, dmodels.ParsePostgresError("Project", err)
	}
	for _, dat := range data {
		country, _ := pImpl.GetCountry(int(dat.CountryId), "vi")
		dat.Country = country

	}
	return data, nil
}

func (pImpl *ProjectImpl) GetByID(id int64) (*domain.Project, error) {
	var data = &domain.Project{}
	var err = pImpl.tblProject().Where("id = ?", id).First(data).Error
	if nil != err {
		return nil, dmodels.ParsePostgresError("Project", err)
	}
	return data, nil
}

func (pImpl *ProjectImpl) ChangeStatus(id string, status domain.ProjectStatus,
) error {
	var err = pImpl.tblProject().
		Where("id = ?", id).
		Update("status", status).
		Error
	return dmodels.ParsePostgresError("Project", err)
}

func (pImpl *ProjectImpl) GetOwner(projectId int64) (string, error) {
	var owner = ""
	var err = pImpl.tblProject().
		Where("id = ?", projectId).
		Pluck("owner", &owner).Error
	if nil != err {
		return "", dmodels.ParsePostgresError("Get owner ", err)
	}
	return owner, nil
}

func (pImpl *ProjectImpl) AddImage(req *domain.RProjectAddImage) (*domain.ProjectImage, error) {
	if req.Type != 0 { // Add thumbnail
		err := pImpl.tblProject().Where("id = ?", req.ProjectId).Update("thumbnail", req.ImgPath).Error
		if err != nil {
			return nil, dmodels.ParsePostgresError("AddImage", err)
		}
		return &domain.ProjectImage{
			ProjectId: req.ProjectId,
			Image:     req.ImgPath,
			CreatedAt: time.Now(),
		}, nil
	}
	img := &domain.ProjectImage{
		ProjectId: req.ProjectId,
		Image:     req.ImgPath,
		CreatedAt: time.Now(),
	}

	err := pImpl.tblImage().Create(img).Error
	if err != nil {
		return nil, dmodels.ParsePostgresError("AddImage", err)
	}

	return nil, nil
}

func (pImpl *ProjectImpl) GetCountry(id int, locale string) (*domain.Country, error) {

	jsonPath := "json/country.json"
	jsonFile, err := os.Open(jsonPath)
	if err != nil {
		return &domain.Country{}, dmodels.ParsePostgresError("Get Country ", err)
	}
	defer jsonFile.Close()
	jsonByte, _ := io.ReadAll(jsonFile)
	countries := map[string][]domain.Language{}
	_ = json.Unmarshal(jsonByte, &countries)
	for _, language := range countries[strconv.Itoa(id)] {
		if language.Locale == locale {
			return &domain.Country{
				Id:   int64(id),
				Name: language.Name,
			}, nil
		}
	}
	return &domain.Country{}, nil
}

func (pImpl *ProjectImpl) tblProject() *gorm.DB {
	return pImpl.db.Table(domain.TableNameProject)
}

func (pImpl *ProjectImpl) tblProjectDesc() *gorm.DB {
	return pImpl.db.Table(domain.TableNameProjectDesc)
}

func (pImpl *ProjectImpl) tblProjectSpec() *gorm.DB {
	return pImpl.db.Table(domain.TableNameProjectSpecs)
}

func (pImpl *ProjectImpl) tblImage() *gorm.DB {
	return pImpl.db.Table(domain.TableNameProjectImage)
}
