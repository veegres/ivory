package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"ivory/persistence"
	"net/http"
)

func (r routes) CertGroup(group *gin.RouterGroup) {
	cert := group.Group("/cert")
	cert.GET("", getCertList)
	cert.DELETE("/:uuid", deleteCert)
	cert.POST("/upload", postUploadCert)
}

func getCertList(context *gin.Context) {
	list, err := persistence.Database.Cert.List()
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		context.JSON(http.StatusOK, gin.H{"response": list})
	}
}

func deleteCert(context *gin.Context) {
	certUuid, err := uuid.Parse(context.Param("uuid"))
	err = persistence.Database.Cert.Delete(certUuid)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		context.JSON(http.StatusOK, gin.H{"response": "deleted"})
	}
}

func postUploadCert(context *gin.Context) {
	file, err := context.FormFile("file")
	cert, err := persistence.Database.Cert.Create(file.Filename)
	err = context.SaveUploadedFile(file, cert.Path)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"error": err.Error()})
	}
	context.JSON(http.StatusOK, gin.H{"response": "no impl"})
}
