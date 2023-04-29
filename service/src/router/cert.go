package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "ivory/src/model"
	"ivory/src/service"
	"net/http"
	"strconv"
)

type CertRouter struct {
	certService *service.CertService
}

func NewCertRouter(certService *service.CertService) *CertRouter {
	return &CertRouter{certService: certService}
}

func (r *CertRouter) GetCertList(context *gin.Context) {
	certType := context.Request.URL.Query().Get("type")

	var err error
	var list map[string]CertModel
	if certType != "" {
		number, _ := strconv.Atoi(certType)
		list, err = r.certService.ListByType(CertType(number))
	} else {
		list, err = r.certService.List()
	}

	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": list})
}

func (r *CertRouter) DeleteCert(context *gin.Context) {
	certUuid, err := uuid.Parse(context.Param("uuid"))
	err = r.certService.Delete(certUuid)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "deleted"})
}

func (r *CertRouter) PostUploadCert(context *gin.Context) {
	certType, errParse := strconv.Atoi(context.PostForm("type"))
	if errParse != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errParse.Error()})
		return
	}
	file, errFile := context.FormFile("file")
	if errFile != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errFile.Error()})
		return
	}
	if file.Size > 1000000 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "maximum size is 1MB"})
		return
	}

	cert, errCreate := r.certService.Create(file.Filename, CertType(certType), FileUsageType(UPLOAD))
	if errCreate != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errCreate.Error()})
		return
	}
	errSave := context.SaveUploadedFile(file, cert.Path)
	if errSave != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errSave.Error()})
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

	cert, err := r.certService.Create(certRequest.Path, certRequest.Type, FileUsageType(PATH))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": cert})
}
