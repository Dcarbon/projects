package repo

import (
	"errors"
	"testing"

	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/go-shared/libs/utils"
	"github.com/Dcarbon/projects/internal/domain"
)

func TestProjectCreate(t *testing.T) {
	db, err := ConnectMockPostgresql()
	if !errors.Is(err, nil) {
		t.Errorf("fail to init Databases.")
		return
	}
	service, err := NewProjectImpl(db)

	specs := map[string]float64{"a": 12344, "b": 121232, "c": 121212}
	req := domain.RProjectCreate{
		Owner:    dmodels.EthAddress("0x5348a62dc343a9fa6881a64b51ea6137968506c5"),
		Location: dmodels.NewCoord4326(111.2222, 232323),
		Specs: &domain.RProjectUpdateSpecs{
			Specs: specs,
		},
		Descs: []*domain.RProjectUpdateDesc{
			{
				Language: "vi",
				Name:     "Description Name",
				Desc:     "Description",
			},
		},
		Area:         1000,
		LocationName: "LOCATION_NAME",
	}

	t.Run("test create project success", func(t *testing.T) {
		project, err := service.Create(&req)
		if err != nil {
			t.Errorf("fail to create project, err= %s", err)
		}
		if project.Id == 0 {
			t.Errorf("fail when create project. ")
		}
		if project.Descs[0].Id == 0 {
			t.Errorf("fail when create project descs ")
		}
		if project.Specs.Id == 0 {
			t.Errorf("fail when create project specs ")
		}
	})
	utils.PanicError("", err)
}

func TestProjectUpdateDesc(t *testing.T) {
	db, err := ConnectMockPostgresql()
	if !errors.Is(err, nil) {
		t.Errorf("fail to init Databases.")
		return
	}
	service, err := NewProjectImpl(db)

	specs := map[string]float64{"a": 12344, "b": 121232, "c": 121212}
	req := domain.RProjectCreate{
		Owner:    dmodels.EthAddress("0x5348a62dc343a9fa6881a64b51ea6137968506c5"),
		Location: dmodels.NewCoord4326(111.2222, 232323),
		Specs: &domain.RProjectUpdateSpecs{
			Specs: specs,
		},
		Descs: []*domain.RProjectUpdateDesc{
			{
				Language: "vi",
				Name:     "Description Name",
				Desc:     "Description",
			},
		},
		Area:         1000,
		LocationName: "LOCATION_NAME",
	}
	t.Run("test update description project fail when project not exists", func(t *testing.T) {
		_, err := service.UpdateDesc(&domain.RProjectUpdateDesc{
			ProjectId: 999,
			Language:  "vi",
			Name:      "name",
			Desc:      "description",
		})
		if err == nil {
			t.Errorf("Update project description fail.")
			return
		}
		//It must not create a new project.
		if _, err := service.GetByID(999); err == nil {
			t.Errorf("get project by id fail when id not exist")
			return
		}
	})

	t.Run("test update description project success", func(t *testing.T) {
		prj, err := service.Create(&req)
		if err != nil {
			t.Errorf("Create new project fail.")
			return
		}
		desc, err := service.UpdateDesc(&domain.RProjectUpdateDesc{
			ProjectId: prj.Id,
			Language:  "vi",
			Name:      "name",
			Desc:      "description",
		})
		if err != nil {
			t.Errorf("Update project description fail.")
			return
		}
		if desc.Id == 0 {
			t.Errorf("Update project description fail.")
			return
		}
	})
	utils.PanicError("", err)
}
func TestProjectUpdateSpecs(t *testing.T) {
	db, err := ConnectMockPostgresql()
	if !errors.Is(err, nil) {
		t.Errorf("fail to init Databases.")
		return
	}
	service, err := NewProjectImpl(db)

	specs := map[string]float64{"a": 12344, "b": 121232, "c": 121212}
	req := domain.RProjectCreate{
		Owner:    dmodels.EthAddress("0x5348a62dc343a9fa6881a64b51ea6137968506c5"),
		Location: dmodels.NewCoord4326(111.2222, 232323),
		Specs: &domain.RProjectUpdateSpecs{
			Specs: specs,
		},
		Descs: []*domain.RProjectUpdateDesc{
			{
				Language: "vi",
				Name:     "Description Name",
				Desc:     "Description",
			},
		},
		Area:         1000,
		LocationName: "LOCATION_NAME",
	}
	t.Run("test update specs project fail when project not exists", func(t *testing.T) {
		_, err := service.UpdateSpecs(&domain.RProjectUpdateSpecs{
			ProjectId: 999,
			Specs:     map[string]float64{"d": 1010101, "asasasa": 101010},
		})
		if err == nil {
			t.Errorf("Update project description fail.")
			return
		}
		//It must not create a new project.
		if _, err := service.GetByID(999); err == nil {
			t.Errorf("get project by id fail when id not exist")
			return
		}
	})

	t.Run("test update specs project success", func(t *testing.T) {
		prj, err := service.Create(&req)
		if err != nil {
			t.Errorf("Create new project fail.")
			return
		}
		desc, err := service.UpdateSpecs(&domain.RProjectUpdateSpecs{
			ProjectId: prj.Id,
			Specs:     map[string]float64{"as": 10101, "bc": 2345},
		})
		if err != nil {
			t.Errorf("Update project description fail.")
			return
		}
		if desc.Id == 0 {
			t.Errorf("Update project description fail.")
			return
		}
	})
	utils.PanicError("", err)
}
