package user

import (
	"ivory/core/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	userService *Service
}

func NewRouter(userService *Service) *Router {
	return &Router{userService: userService}
}

func (r *Router) GetUserList(context *gin.Context) {
	users, err := r.userService.List()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": users})
}

func (r *Router) DeleteUser(context *gin.Context) {
	err := r.userService.Delete(context.Param("username"), r.getRequester(context))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "user was deleted"})
}

func (r *Router) PostUserPassword(context *gin.Context) {
	var request PasswordUpdateRequest
	if errBind := context.ShouldBindJSON(&request); errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	err := r.userService.UpdatePassword(r.getRequester(context), request.PreviousPassword, request.NewPassword)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "password was updated"})
}

func (r *Router) GetUserLinkList(context *gin.Context) {
	links, err := r.userService.LinkList()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": links})
}

func (r *Router) PostUserLink(context *gin.Context) {
	var request LinkRequest
	if errBind := context.ShouldBindJSON(&request); errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	link, err := r.userService.LinkCreateInvite(request, r.getRequester(context))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": link})
}

func (r *Router) PostUserResetLink(context *gin.Context) {
	var request LinkResetRequest
	if errBind := context.ShouldBindJSON(&request); errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	link, err := r.userService.LinkCreateReset(request.Username, r.getRequester(context))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": link})
}

func (r *Router) DeleteUserLink(context *gin.Context) {
	err := r.userService.LinkRevoke(context.Param("uuid"), r.getRequester(context))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "link is not valid any more"})
}

func (r *Router) PostUserLinkVerify(context *gin.Context) {
	var request LinkVerifyRequest
	if errBind := context.ShouldBindJSON(&request); errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	payload, err := r.userService.LinkVerify(request.Token)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": payload})
}

func (r *Router) PostUserLinkPassword(context *gin.Context) {
	var request LinkPasswordRequest
	if errBind := context.ShouldBindJSON(&request); errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	user, err := r.userService.LinkApply(request.Token, request.Password)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// NOTE: this endpoint sets a password and nothing else. Signing the person
	// in afterwards is the login endpoint's job, and the page calls it with the
	// credentials it already holds - so neither feature has to know the other
	context.JSON(http.StatusOK, gin.H{"response": user})
}

// getRequester names who is asking. It is empty only when Ivory runs without
// authentication, where there is no identity and everything is permitted.
func (r *Router) getRequester(context *gin.Context) string {
	return context.GetString(config.AuthContextKey.Username)
}
