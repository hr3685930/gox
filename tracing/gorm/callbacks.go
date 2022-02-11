package gorm

import (
	"gorm.io/gorm"
)

func (op opentracingPlugin)beforeCreate(db *gorm.DB){
	op.injectBefor(db, _createOp)
}

func (op opentracingPlugin)beforeUpdate(db *gorm.DB){
	op.injectBefor(db, _updateOp)
}

func (op opentracingPlugin)beforeQuery(db *gorm.DB){
	op.injectBefor(db, _queryOp)
}

func (op opentracingPlugin)beforeDelete(db *gorm.DB){
	op.injectBefor(db, _deleteOp)
}

func (op opentracingPlugin)beforeRow(db *gorm.DB){
	op.injectBefor(db, _rowOp)
}

func (op opentracingPlugin)beforeRaw(db *gorm.DB){
	op.injectBefor(db, _rawOp)
}

func (op opentracingPlugin)after(db *gorm.DB){
	op.extractAfter(db)
}