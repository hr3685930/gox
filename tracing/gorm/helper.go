package gorm


import (
	"context"
	//"github.com/blinkbean/jaeger-module/jaegersql"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"gorm.io/gorm"
)

const (
	_prefix      = "gorm"
	_errorTagKey = "error"
)

var (
	// span.Tag keys
	_tableTagKey = keyWithPrefix("table")
	// span.Log keys
	//_errorLogKey        = keyWithPrefix("error")
	_resultLogKey       = keyWithPrefix("result")
	_sqlLogKey          = keyWithPrefix("sql")
	_rowsAffectedLogKey = keyWithPrefix("rowsAffected")
)

func keyWithPrefix(key string) string {
	return _prefix + "." + key
}

var (
	opentracingSpanKey = "opentracing:span"
)

func (op opentracingPlugin) injectBefor(db *gorm.DB, opName operationName) {
	if db == nil {
		return
	}
	if db.Statement == nil || db.Statement.Context == nil {
		db.Logger.Error(context.TODO(), "could not inject sp from nil Statement.Context or nil Statement")
		return
	}

	sp, _ := opentracing.StartSpanFromContext(db.Statement.Context, op.Name())
	ext.DBType.Set(sp, "db.mysql")
	db.InstanceSet(opentracingSpanKey, sp)
}

func (op opentracingPlugin) extractAfter(db *gorm.DB) {
	if db == nil {
		return
	}
	if db.Statement == nil || db.Statement.Context == nil {
		db.Logger.Error(context.TODO(), "could not extract sp from nil Statement.Context or nil Statement")
		return
	}

	v, ok := db.InstanceGet(opentracingSpanKey)
	if !ok || v == nil {
		return
	}
	sp, ok := v.(opentracing.Span)
	if !ok || sp == nil {
		return
	}
	sp.SetOperationName(QuerySignature(db.Statement.SQL.String()))
	defer sp.Finish()

	op.tag(sp, db)
}

func (op opentracingPlugin) tag(sp opentracing.Span, db *gorm.DB) {
	// 可以加个选择开关
	if err := db.Error; err != nil {
		sp.SetTag(_errorTagKey, true)
	}
	sp.SetTag(_sqlLogKey, db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...))
}
