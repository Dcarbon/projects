package infra

import (
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
	var project = req.ToProject()
	var e1 = pImpl.tblProject().Transaction(func(dbTx *gorm.DB) error {
		err := dbTx.Table(domain.TableNameProject).Create(project).Error
		if nil != err {
			return dmodels.ParsePostgresError("Create project", err)
		}

		return nil
	})

	if nil != e1 {
		return nil, e1
	}

	return project, nil
}

func (pImpl *ProjectImpl) UpdateDesc(req *domain.RProjectUpdateDesc,
) (*domain.ProjectDesc, error) {
	var desc = req.ToProjectDesc()

	var err = pImpl.tblProjectDesc().
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{
					{Name: "project_id"}, {Name: "language"},
				},
				UpdateAll: true,
			},
			clause.Insert{},
		).
		Create(desc).Error
	if nil != err {
		return nil, dmodels.ParsePostgresError("Update project desc", err)
	}
	return desc, nil
}

func (pImpl *ProjectImpl) UpdateSpecs(req *domain.RProjectUpdateSpecs,
) (*domain.ProjectSpecs, error) {
	var spec = req.ToProjectSpecs()

	var err = pImpl.tblProjectSpec().
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{{Name: "project_id"}},
				DoUpdates: clause.AssignmentColumns(
					[]string{"specs", "updated_at"},
				),
			},
		).Create(spec).Error
	if nil != err {
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
	var err = tbl.Find(&data).Error
	if nil != err {
		return nil, dmodels.ParsePostgresError("Project", err)
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

func (pImpl *ProjectImpl) AddImage(req *domain.RProjectAddImage,
) (*domain.ProjectImage, error) {
	var img = &domain.ProjectImage{
		Id:        0,
		ProjectId: req.ProjectId,
		Image:     req.ImgPath,
		CreatedAt: time.Now(),
	}
	var err = pImpl.tblImage().Create(img).Error
	if nil != err {
		return nil, dmodels.ParsePostgresError("AddImage ", err)
	}
	return img, nil
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
