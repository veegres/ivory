package auth

import (
	"errors"
	"ivory/clients/auth/basic"
	"ivory/clients/auth/ldap"
	"ivory/clients/auth/oidc"
	"ivory/core/config"
	"net/http"
	"strings"
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
		session, errCookie := context.Cookie("session")
		if errCookie != nil {
			id, errToken := uuid.NewUUID()
			if errToken != nil {
				context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": errToken.Error()})
				return
			}
			session = id.String()
			r.setCookieSession(context, session)
		}
		context.Set(config.AuthContextKey.Session, session)
		context.Next()
	}
}

func (r *Router) ValidateMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set(config.AuthContextKey.Enabled, true)
		valid, username, authType, errParse := r.resolveAuth(context)
		if !valid {
			context.Header("WWW-Authenticate", "Bearer JWT realm="+r.authService.getIssuer())
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errParse.Error()})
			return
		}
		if errors.Is(errParse, ErrAuthDisabled) {
			context.Set(config.AuthContextKey.Enabled, false)
		} else {
			if username == "" {
				context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrUsernameEmpty.Error()})
				return
			}
			if authType == nil {
				context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidAuthType.Error()})
				return
			}
			context.Set(config.AuthContextKey.Username, username)
			context.Set(config.AuthContextKey.Type, authType.String())
			context.Set(config.AuthContextKey.Superuser, r.authService.IsSuperuser(username))
		}
		context.Next()
	}
}

func (r *Router) AuthContextMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set(config.AuthContextKey.Enabled, true)
		valid, username, authType, errParse := r.resolveAuth(context)
		context.Set(config.AuthContextKey.Authorised, valid)
		if errors.Is(errParse, ErrAuthDisabled) {
			context.Set(config.AuthContextKey.Enabled, false)
		} else {
			if errParse != nil {
				context.Set(config.AuthContextKey.Error, errParse.Error())
			}
			context.Set(config.AuthContextKey.Username, username)
			context.Set(config.AuthContextKey.Superuser, r.authService.IsSuperuser(username))
			if authType != nil {
				context.Set(config.AuthContextKey.Type, authType.String())
			}
		}
		context.Next()
	}
}

func (r *Router) resolveAuth(context *gin.Context) (bool, string, *AuthType, error) {
	headerToken, errHeader := r.getHeaderToken(context)
	cookieToken, errCookie := r.getCookieToken(context)
	return r.authService.ParseAuthTokenWithFallback(headerToken, errHeader, cookieToken, errCookie)
}

func (r *Router) getHeaderToken(context *gin.Context) (string, error) {
	authHeader := context.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthorizationToken
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", ErrInvalidAuthorizationHeader
	}
	return parts[1], nil
}

func (r *Router) getCookieToken(context *gin.Context) (string, error) {
	cookieToken, errToken := context.Cookie("token")
	if cookieToken == "" || errToken != nil {
		cookieTokenError, _ := context.Cookie("token_error")
		if cookieTokenError != "" {
			return "", errors.New(cookieTokenError)
		}
		return "", ErrNoAuthorizationToken
	}
	return cookieToken, nil
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
		r.handleTokenError(context, ErrStateCookieNotFound.Error())
		return
	}

	if context.Query("state") != state {
		r.handleTokenError(context, ErrInvalidStateParameter.Error())
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
