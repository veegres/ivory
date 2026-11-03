package query

import (
	"errors"
	"ivory/plugins/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func writeQueryAdminError(ctx *gin.Context, err error) {
	if errors.Is(err, database.ErrAdminNotSupported) {
		ctx.JSON(http.StatusNotImplemented, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func (r *HttpRouter) PostDbRoleList(context *gin.Context) {
	var req RouterDbRoleListRequest
	if errBind := context.ShouldBindJSON(&req); errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	r.withQueryContext(context, req.Connection, func(queryContext Context) {
		res, err := r.service.DbRoleList(queryContext)
		if err != nil {
			writeQueryAdminError(context, err)
			return
		}
		context.JSON(http.StatusOK, gin.H{"response": res})
	})
}

func (r *HttpRouter) PostDbRoleCreate(context *gin.Context) {
	var req RouterDbRoleCreateRequest
	if errBind := context.ShouldBindJSON(&req); errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	r.withQueryContext(context, req.Connection, func(queryContext Context) {
		if err := r.service.DbRoleCreate(queryContext, req.Name, req.Password); err != nil {
			writeQueryAdminError(context, err)
			return
		}
		context.JSON(http.StatusOK, gin.H{"response": "done"})
	})
}

func (r *HttpRouter) PostDbRoleUpdate(context *gin.Context) {
	var req RouterDbRoleUpdateRequest
	if errBind := context.ShouldBindJSON(&req); errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	r.withQueryContext(context, req.Connection, func(queryContext Context) {
		if err := r.service.DbRoleUpdatePassword(queryContext, req.Name, req.Password); err != nil {
			writeQueryAdminError(context, err)
			return
		}
		context.JSON(http.StatusOK, gin.H{"response": "done"})
	})
}

func (r *HttpRouter) PostDbRoleDelete(context *gin.Context) {
	var req RouterDbRoleDeleteRequest
	if errBind := context.ShouldBindJSON(&req); errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	r.withQueryContext(context, req.Connection, func(queryContext Context) {
		if err := r.service.DbRoleDelete(queryContext, req.Name); err != nil {
			writeQueryAdminError(context, err)
			return
		}
		context.JSON(http.StatusOK, gin.H{"response": "done"})
	})
}

func (r *HttpRouter) PostDbDatabaseAdminList(context *gin.Context) {
	var req RouterDbDatabaseListRequest
	if errBind := context.ShouldBindJSON(&req); errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	r.withQueryContext(context, req.Connection, func(queryContext Context) {
		res, err := r.service.DbDatabaseList(queryContext)
		if err != nil {
			writeQueryAdminError(context, err)
			return
		}
		context.JSON(http.StatusOK, gin.H{"response": res})
	})
}

func (r *HttpRouter) PostDbDatabaseAdminCreate(context *gin.Context) {
	var req RouterDbDatabaseCreateRequest
	if errBind := context.ShouldBindJSON(&req); errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	r.withQueryContext(context, req.Connection, func(queryContext Context) {
		if err := r.service.DbDatabaseCreate(queryContext, req.Name); err != nil {
			writeQueryAdminError(context, err)
			return
		}
		context.JSON(http.StatusOK, gin.H{"response": "done"})
	})
}

func (r *HttpRouter) PostDbDatabaseAdminUpdate(context *gin.Context) {
	var req RouterDbDatabaseUpdateRequest
	if errBind := context.ShouldBindJSON(&req); errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	r.withQueryContext(context, req.Connection, func(queryContext Context) {
		if err := r.service.DbDatabaseRename(queryContext, req.Name, req.NewName); err != nil {
			writeQueryAdminError(context, err)
			return
		}
		context.JSON(http.StatusOK, gin.H{"response": "done"})
	})
}

func (r *HttpRouter) PostDbDatabaseAdminDelete(context *gin.Context) {
	var req RouterDbDatabaseDeleteRequest
	if errBind := context.ShouldBindJSON(&req); errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	r.withQueryContext(context, req.Connection, func(queryContext Context) {
		if err := r.service.DbDatabaseDelete(queryContext, req.Name); err != nil {
			writeQueryAdminError(context, err)
			return
		}
		context.JSON(http.StatusOK, gin.H{"response": "done"})
	})
}
