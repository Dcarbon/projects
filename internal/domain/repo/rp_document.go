package repo

import (
	"time"

	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/projects/internal/domain"
	"gorm.io/gorm/clause"
)

func (pImpl *ProjectImpl) UpsertDocument(req *domain.RProjectDocumentUpsert) ([]*domain.ProjectDocument, error) {
	documents := []*domain.ProjectDocument{}
	for _, val := range req.Document {
		documents = append(documents,
			&domain.ProjectDocument{
				Url:          val.Url,
				DocumentType: val.DocumentType,
				ProjectId:    val.ProjectId,
				Id:           val.Id,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now()})
	}

	if err := pImpl.tblDocument().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}}, // key column
		DoUpdates: clause.AssignmentColumns([]string{"url", "document_type", "updated_at"}),
	}).Create(&documents).Error; err != nil {
		return nil, err
	}
	return documents, nil
}

func (pImpl *ProjectImpl) DeleteDocument(req *domain.RProjectDocumentDelete) error {
	if err := pImpl.tblDocument().
		Where("id IN ?", req.Id).Delete(&domain.ProjectDocument{}).Error; nil != err {
		return dmodels.ParsePostgresError("Delete Document ", err)
	}
	return nil
}
func (pImpl *ProjectImpl) ListDocument(req *domain.RProjectDocumentList) ([]*domain.ProjectDocument, int64, error) {
	documents := []*domain.ProjectDocument{}
	var count int64
	var tbl = pImpl.tblDocument().Where("deleted_at IS NULL")
	if len(req.Ids) > 0 {
		tbl = tbl.Where("id = ?", req.Ids)
	}
	tbl.Count(&count).Offset(req.Skip)
	if req.Limit > 0 {
		tbl = tbl.Limit(req.Limit)
	}
	err := tbl.Order("created_at DESC").Find(&documents).Error
	if err != nil {
		return nil, 0, dmodels.ParsePostgresError("List Document", err)
	}
	return documents, count, nil
}
