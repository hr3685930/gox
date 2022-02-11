package gorm

type gormOption string

func (g gormOption) Name() string {
	return string(g)
}

const (
	_gormCreate gormOption = "gorm:create"
	_gormUpdate gormOption = "gorm:update"
	_gormQuery  gormOption = "gorm:query"
	_gormDelete gormOption = "gorm:delete"
	_gormRow    gormOption = "gorm:row"
	_gormRaw    gormOption = "gorm:raw"
)

type operationName string

func (op operationName) String() string {
	return string(op)
}

const (
	_createOp operationName = "create"
	_updateOp operationName = "update"
	_queryOp  operationName = "query"
	_deleteOp operationName = "delete"
	_rowOp    operationName = "row"
	_rawOp    operationName = "raw"
)

type operationStage string

func (op operationStage) Name() string {
	return string(op)
}

const (
	_stageBeforeCreate operationStage = "opentracing:before_create"
	_stageAfterCreate  operationStage = "opentracing:after_create"
	_stageBeforeUpdate operationStage = "opentracing:before_update"
	_stageAfterUpdate  operationStage = "opentracing:after_update"
	_stageBeforeQuery  operationStage = "opentracing:before_query"
	_stageAfterQuery   operationStage = "opentracing:after_query"
	_stageBeforeDelete operationStage = "opentracing:before_delete"
	_stageAfterDelete  operationStage = "opentracing:after_delete"
	_stageBeforeRow    operationStage = "opentracing:before_row"
	_stageAfterRow     operationStage = "opentracing:after_row"
	_stageBeforeRaw    operationStage = "opentracing:before_raw"
	_stageAfterRaw     operationStage = "opentracing:after_raw"
)