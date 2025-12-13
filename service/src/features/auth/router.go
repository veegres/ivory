package auth

import (
	"errors"
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

func (r *Router) ValidateMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set("auth", true)
		valid, username, authType, errParse := r.authService.ParseAuthToken(context)
		if !valid {
			context.Header("WWW-Authenticate", "Bearer JWT realm="+r.authService.getIssuer())
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errParse.Error()})
			return
		}
		if errors.Is(errParse, ErrAuthDisabled) {
			context.Set("auth", false)
		} else {
			if username == "" {
				context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "username cannot be empty"})
				return
			}
			if authType == nil {
				context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth type"})
				return
			}
		}
		context.Set("username", username)
		context.Set("authType", authType.String())
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
	r.Logout(context)
	var login basic.Login
	parseErr := context.ShouldBindJSON(&login)
	if parseErr != nil {
		r.setCookieTokenError(context, parseErr.Error())
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}

	token, exp, errToken := r.authService.GenerateBasicAuthToken(login)
	if errToken != nil {
		r.setCookieTokenError(context, errToken.Error())
		context.JSON(http.StatusUnauthorized, gin.H{"error": errToken.Error()})
		return
	}

	r.setCookieToken(context, token, exp)
	context.JSON(http.StatusOK, gin.H{"response": gin.H{"token": token, "expire": exp.String()}})
}

func (r *Router) LdapLogin(context *gin.Context) {
	r.Logout(context)
	var login ldap.Login
	parseErr := context.ShouldBindJSON(&login)
	if parseErr != nil {
		r.setCookieTokenError(context, parseErr.Error())
		context.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}

	token, exp, errToken := r.authService.GenerateLdapAuthToken(login)
	if errToken != nil {
		r.setCookieTokenError(context, errToken.Error())
		context.JSON(http.StatusUnauthorized, gin.H{"error": errToken.Error()})
		return
	}

	r.setCookieToken(context, token, exp)
	context.JSON(http.StatusOK, gin.H{"response": gin.H{"token": token, "expire": exp.String()}})
}

func (r *Router) OidcLogin(context *gin.Context) {
	r.Logout(context)
	state, errState := uuid.NewUUID()
	if errState != nil {
		r.handleTokenError(context, errState.Error())
		return
	}
	codeUrl, errCode := r.authService.oidcProvider.GetCode(state.String())
	if errCode != nil {
		r.handleTokenError(context, errCode.Error())
		return
	}

	r.setCookieState(context, state.String())
	http.Redirect(context.Writer, context.Request, codeUrl, http.StatusFound)
}

func (r *Router) OidcCallback(context *gin.Context) {
	err := context.Query("error")
	if err != "" {
		errDesc := context.Query("error_description")
		r.handleTokenError(context, err+": "+errDesc)
		return
	}

	state, errState := context.Cookie("state")
	if errState != nil {
		r.handleTokenError(context, "state cookie not found")
		return
	}

	if context.Query("state") != state {
		r.handleTokenError(context, "invalid state parameter")
		return
	}

	token, exp, errToken := r.authService.GenerateOidcAuthToken(context.Query("code"))
	if errToken != nil {
		r.handleTokenError(context, errToken.Error())
		return
	}

	r.setCookieToken(context, token, exp)
	http.Redirect(context.Writer, context.Request, r.path, http.StatusFound)
}

func (r *Router) Logout(context *gin.Context) {
	context.SetCookie("token", "", -1, r.path, "", r.tlsEnabled, true)
	context.SetCookie("token_error", "", -1, r.path, "", r.tlsEnabled, true)
}

func (r *Router) handleTokenError(context *gin.Context, err string) {
	r.setCookieTokenError(context, err)
	path := r.path + "?error=" + err
	http.Redirect(context.Writer, context.Request, path, http.StatusFound)
}

func (r *Router) setCookieSession(context *gin.Context, value string) {
	// NOTE: maxAge is provided in seconds 2592000 sec = 30 days
	context.SetCookie("session", value, 2592000, r.path, "", r.tlsEnabled, true)
}

func (r *Router) setCookieToken(context *gin.Context, value string, exp *time.Time) {
	// NOTE: we add 60 sec just to see proper error when token has expired
	seconds := int(time.Until(*exp).Seconds()) + 60
	context.SetSameSite(http.SameSiteStrictMode)
	context.SetCookie("token_error", "", -1, r.path, "", r.tlsEnabled, true)
	context.SetCookie("token", value, seconds, r.path, "", r.tlsEnabled, true)

}

func (r *Router) setCookieTokenError(context *gin.Context, value string) {
	context.SetSameSite(http.SameSiteStrictMode)
	context.SetCookie("token", "", -1, r.path, "", r.tlsEnabled, true)
	context.SetCookie("token_error", value, 0, r.path, "", r.tlsEnabled, true)
}

func (r *Router) setCookieState(context *gin.Context, value string) {
	path := ""
	if r.path != "/" {
		path = r.path
	}
	context.SetCookie("state", value, 600, path+"/api/oidc", "", r.tlsEnabled, true)
}
