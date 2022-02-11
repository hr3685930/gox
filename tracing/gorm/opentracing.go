package gorm

import (
	"gorm.io/gorm"
	"strings"
)


type opentracingPlugin struct {}

func NewOpentracingPlugin() *opentracingPlugin {
	return &opentracingPlugin{}
}

func (op opentracingPlugin) Name() string {
	return "gorm"
}

func (op opentracingPlugin) Initialize(db *gorm.DB) (err error) {
	e := myError{errs: make([]string, 0, 12)}

	err = db.Callback().Create().Before(_gormCreate.Name()).Register(_stageBeforeCreate.Name(), op.beforeCreate)
	e.add(_stageBeforeCreate, err)
	err = db.Callback().Create().After(_gormCreate.Name()).Register(_stageAfterCreate.Name(), op.after)
	e.add(_stageAfterCreate, err)

	err = db.Callback().Update().Before(_gormUpdate.Name()).Register(_stageBeforeUpdate.Name(), op.beforeUpdate)
	e.add(_stageBeforeUpdate, err)
	err = db.Callback().Update().After(_gormUpdate.Name()).Register(_stageAfterUpdate.Name(), op.after)
	e.add(_stageAfterUpdate, err)

	err = db.Callback().Query().Before(_gormUpdate.Name()).Register(_stageBeforeQuery.Name(), op.beforeQuery)
	e.add(_stageBeforeQuery, err)
	err = db.Callback().Query().After(_gormUpdate.Name()).Register(_stageAfterQuery.Name(), op.after)
	e.add(_stageAfterQuery, err)

	err = db.Callback().Delete().Before(_gormDelete.Name()).Register(_stageBeforeDelete.Name(), op.beforeDelete)
	e.add(_stageBeforeDelete, err)
	err = db.Callback().Delete().After(_gormDelete.Name()).Register(_stageAfterDelete.Name(),op.after)
	e.add(_stageAfterDelete, err)

	err = db.Callback().Row().Before(_gormRow.Name()).Register(_stageBeforeRow.Name(), op.beforeRow)
	e.add(_stageBeforeRow, err)
	err = db.Callback().Row().After(_gormRow.Name()).Register(_stageAfterRow.Name(), op.after)
	e.add(_stageAfterRow, err)

	err = db.Callback().Raw().Before(_gormRaw.Name()).Register(_stageBeforeRaw.Name(), op.beforeRaw)
	e.add(_stageBeforeRaw, err)
	err = db.Callback().Raw().After(_gormRaw.Name()).Register(_stageAfterRaw.Name(), op.after)
	e.add(_stageAfterRaw, err)

	return e.toError()
}

type myError struct {
	errs [] string
}

func (m *myError) add(stage operationStage,err error) {
	if err == nil{
		return
	}
	m.errs = append(m.errs, "stage="+stage.Name()+":"+err.Error())
}

func (m myError) toError() error {
	if len(m.errs)==0{
		return nil
	}
	return m
}

func (m myError) Error() string {
	return strings.Join(m.errs, ";")
}