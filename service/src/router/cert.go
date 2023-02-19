package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "ivory/src/model"
	"ivory/src/persistence"
	"net/http"
	"strconv"
)

type CertRouter struct {
	repository *persistence.CertRepository
}

func NewCertRouter(repository *persistence.CertRepository) *CertRouter {
	return &CertRouter{repository: repository}
}

func (r *CertRouter) GetCertList(context *gin.Context) {
	certType := context.Request.URL.Query().Get("type")

	var err error
	var list map[string]CertModel
	if certType != "" {
		number, _ := strconv.Atoi(certType)
		list, err = r.repository.ListByType(CertType(number))
	} else {
		list, err = r.repository.List()
	}

	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": list})
}

func (r *CertRouter) DeleteCert(context *gin.Context) {
	certUuid, err := uuid.Parse(context.Param("uuid"))
	err = r.repository.Delete(certUuid)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "deleted"})
}

func (r *CertRouter) PostUploadCert(context *gin.Context) {
	certType, err := strconv.Atoi(context.PostForm("type"))
	file, err := context.FormFile("file")
	if file.Size > 1000000 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "maximum size is 1MB"})
		return
	}

	cert, err := r.repository.Create(file.Filename, CertType(certType), FileUsageType(UPLOAD))
	err = context.SaveUploadedFile(file, cert.Path)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": cert})
}

func (r *CertRouter) PostAddCert(context *gin.Context) {
	var certRequest CertRequest
	parseErr := context.ShouldBindJSON(&certRequest)
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}

	cert, err := r.repository.Create(certRequest.Path, certRequest.Type, FileUsageType(PATH))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": cert})
}
