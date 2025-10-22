package auth

import (
	"ivory/src/clients/auth/basic"
	"ivory/src/clients/auth/ldap"
	"ivory/src/clients/auth/oidc"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Router struct {
	authService *Service

	path       string
	tlsEnabled bool
}

func NewRouter(authService *Service, path string, tlsEnabled bool) *Router {
	return &Router{authService: authService, path: path, tlsEnabled: tlsEnabled}
}

func (r *Router) SessionMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		_, errCookie := context.Cookie("session")
		if errCookie != nil {
			session, errToken := uuid.NewUUID()
			if errToken != nil {
				context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": errToken.Error()})
				return
			}
			r.setCookieSession(context, session.String())
		}
		context.Next()
	}
}

func (r *Router) AuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		valid, errValid := r.authService.ValidateAuthToken(context)
		if !valid {
			context.Header("WWW-Authenticate", "Bearer JWT realm="+r.authService.getIssuer())
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errValid.Error()})
			return
		}
		context.Next()
	}
}

func (r *Router) LdapConnect(context *gin.Context) {
	var config ldap.Config
	parseErr := context.ShouldBindJSON(&config)
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}

	err := r.authService.ldapProvider.Connect(config)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": "connected"})
}

func (r *Router) OidcConnect(context *gin.Context) {
	var config oidc.Config
	parseErr := context.ShouldBindJSON(&config)
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}

	err := r.authService.oidcProvider.Connect(config)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": "connected"})
}

func (r *Router) BasicLogin(context *gin.Context) {
	var login basic.Login
	parseErr := context.ShouldBindJSON(&login)
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}

	token, exp, errToken := r.authService.GenerateBasicAuthToken(login)
	if errToken != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": errToken.Error()})
		return
	}

	r.setCookieToken(context, token, exp)
	context.JSON(http.StatusOK, gin.H{"response": gin.H{"token": token, "expire": exp.String()}})
}

func (r *Router) LdapLogin(context *gin.Context) {
	var login ldap.Login
	parseErr := context.ShouldBindJSON(&login)
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}

	token, exp, errToken := r.authService.GenerateLdapAuthToken(login)
	if errToken != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": errToken.Error()})
		return
	}

	r.setCookieToken(context, token, exp)
	context.JSON(http.StatusOK, gin.H{"response": gin.H{"token": token, "expire": exp.String()}})
}

func (r *Router) OidcLogin(context *gin.Context) {
	state, errState := uuid.NewUUID()
	if errState != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": errState.Error()})
		return
	}
	codeUrl, errCode := r.authService.getOAuthCodeURL(state.String())
	if errCode != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": errCode.Error()})
		return
	}

	r.setCookieState(context, state.String())
	http.Redirect(context.Writer, context.Request, codeUrl, http.StatusFound)
}

func (r *Router) OidcCallback(context *gin.Context) {
	state, errState := context.Cookie("state")
	if errState != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "state cookie not found"})
		return
	}

	if context.Query("state") != state {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid state parameter"})
		return
	}

	token, exp, errToken := r.authService.GenerateOidcAuthToken(context.Query("code"))
	if errToken != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": errToken.Error()})
		return
	}

	r.setCookieToken(context, token, exp)
	http.Redirect(context.Writer, context.Request, r.path, http.StatusFound)
}

func (r *Router) Logout(context *gin.Context) {
	context.SetCookie("token", "", -1, r.path, "", r.tlsEnabled, true)
}

func (r *Router) setCookieSession(context *gin.Context, value string) {
	// NOTE: maxAge is provided in seconds 2592000 sec = 30 days
	context.SetCookie("session", value, 2592000, r.path, "", r.tlsEnabled, true)
}

func (r *Router) setCookieToken(context *gin.Context, value string, exp *time.Time) {
	seconds := int(time.Until(*exp).Seconds())
	context.SetSameSite(http.SameSiteStrictMode)
	context.SetCookie("token", value, seconds, r.path, "", r.tlsEnabled, true)
}

func (r *Router) setCookieState(context *gin.Context, value string) {
	path := ""
	if r.path != "/" {
		path = r.path
	}
	context.SetCookie("state", value, 600, path+"/api/oidc", "", r.tlsEnabled, true)
}
