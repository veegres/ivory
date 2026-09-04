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

func (r *Router) PostUser(context *gin.Context) {
	var request CreateRequest
	if errBind := context.ShouldBindJSON(&request); errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	user, err := r.userService.Create(request, r.getRequester(context))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": user})
}

func (r *Router) PutUser(context *gin.Context) {
	var request UpdateRequest
	if errBind := context.ShouldBindJSON(&request); errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	user, err := r.userService.Update(context.Param("username"), request.AuthTypes, r.getRequester(context))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": user})
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

func (r *Router) PostUserPasswordReset(context *gin.Context) {
	registration, err := r.userService.PasswordResetIssue(context.Param("username"), r.getRequester(context))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": registration})
}

func (r *Router) DeleteUserPasswordReset(context *gin.Context) {
	err := r.userService.PasswordResetRevoke(context.Param("username"), r.getRequester(context))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": "registration is not valid any more"})
}

func (r *Router) PostUserRegistrationVerify(context *gin.Context) {
	var request RegistrationVerifyRequest
	if errBind := context.ShouldBindJSON(&request); errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	payload, err := r.userService.RegistrationVerify(request.Token)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": payload})
}

func (r *Router) PostUserRegistrationPassword(context *gin.Context) {
	var request RegistrationPasswordRequest
	if errBind := context.ShouldBindJSON(&request); errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}
	user, err := r.userService.RegistrationApply(request.Token, request.Password)
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
