package middleware

import (
	"slingshot/db"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

func NewAdapter() (*gormadapter.Adapter, error) {
	if ga, err := gormadapter.NewAdapterByDB(db.DB); err != nil {
		return nil, err
	} else {
		return ga, nil
	}
}

// Construction function
func NewEnforcer(modelString string, adapter gormadapter.Adapter) (*casbin.Enforcer, error) {
	m, err := loadModelFromString(modelString)
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}
	return e, nil
}

// func loadModelFromString(modelString string) (model.Model, error) {
// 	m, err := model.NewModelFromString(modelString)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return m, nil
// }

func loadModelFromString(modelString string) (model.Model, error) {
	m := model.NewModel()
	err := m.LoadModelFromText(string(modelString))
	if err != nil {
		return nil, err
	}
	return m, nil
}
