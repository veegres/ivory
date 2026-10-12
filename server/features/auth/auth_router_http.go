package auth

import (
	"errors"
	"ivory/clients/auth/basic"
	"ivory/clients/auth/ldap"
	"ivory/clients/auth/oidc"
	"ivory/core/config"
	"ivory/features/user"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HttpRouter struct {
	authService *Service

	path       string
	tlsEnabled bool
}

func NewHttpRouter(authService *Service, path string, tlsEnabled bool) *HttpRouter {
	return &HttpRouter{authService: authService, path: path, tlsEnabled: tlsEnabled}
}

func (r *HttpRouter) SessionMiddleware() gin.HandlerFunc {
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

func (r *HttpRouter) AllowMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		if context.GetBool(config.AuthContextKey.Enabled) && context.GetBool(config.AuthContextKey.Authorised) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": ErrAlreadyAuthenticated.Error()})
			return
		}
		context.Next()
	}
}

func (r *HttpRouter) RejectMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		if !context.GetBool(config.AuthContextKey.Enabled) {
			context.Next()
			return
		}
		valid := context.GetBool(config.AuthContextKey.Authorised)
		validError := context.GetError(config.AuthContextKey.Error)
		if !valid || validError != nil {
			context.Header("WWW-Authenticate", "Bearer JWT realm="+r.authService.getIssuer())
			message := ErrNoAuthorizationToken.Error()
			if validError != nil {
				message = validError.Error()
			}
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": message})
			return
		}
		username := context.GetString(config.AuthContextKey.Username)
		if username == "" {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrUsernameEmpty.Error()})
			return
		}
		authType := context.GetString(config.AuthContextKey.Type)
		if !user.AuthType(authType).Valid() {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidAuthType.Error()})
			return
		}
		context.Next()
	}
}

func (r *HttpRouter) ValidateWithContextMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set(config.AuthContextKey.Enabled, true)
		valid, username, authType, errParse := r.resolveAuth(context)
		context.Set(config.AuthContextKey.Authorised, valid)
		if valid && errParse == nil {
			r.refreshCookieToken(context)
		}
		if errors.Is(errParse, ErrAuthDisabled) {
			context.Set(config.AuthContextKey.Enabled, false)
		} else {
			if errParse != nil {
				context.Set(config.AuthContextKey.Error, errParse)
			}
			context.Set(config.AuthContextKey.Username, username)
			context.Set(config.AuthContextKey.Superuser, r.authService.IsSuperuser(username))
			if authType != nil {
				context.Set(config.AuthContextKey.Type, string(*authType))
			}
		}
		context.Next()
	}
}

func (r *HttpRouter) resolveAuth(context *gin.Context) (bool, string, *AuthType, error) {
	headerToken, errHeader := r.getHeaderToken(context)
	cookieToken, errCookie := r.getCookieToken(context)
	return r.authService.ParseAuthTokenWithFallback(headerToken, errHeader, cookieToken, errCookie)
}

// refreshCookieToken renews only a token that arrived in the cookie; one sent in the header is never touched.
func (r *HttpRouter) refreshCookieToken(context *gin.Context) {
	headerToken, errHeader := r.getHeaderToken(context)
	if errHeader == nil {
		if valid, _, _, _ := r.authService.ParseAuthToken(headerToken, nil); valid {
			return
		}
	}
	cookieToken, errCookie := r.getCookieToken(context)
	if errCookie != nil {
		return
	}
	token, exp, errRefresh := r.authService.RefreshToken(cookieToken)
	if errRefresh != nil {
		slog.Warn("failed to refresh the session cookie", "error", errRefresh)
		return
	}
	if token != "" {
		r.setCookieToken(context, token, exp)
	}
}

func (r *HttpRouter) getHeaderToken(context *gin.Context) (string, error) {
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

func (r *HttpRouter) getCookieToken(context *gin.Context) (string, error) {
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

func (r *HttpRouter) LdapConnect(context *gin.Context) {
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

func (r *HttpRouter) OidcConnect(context *gin.Context) {
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

func (r *HttpRouter) BasicLogin(context *gin.Context) {
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

func (r *HttpRouter) LdapLogin(context *gin.Context) {
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

func (r *HttpRouter) OidcLogin(context *gin.Context) {
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

func (r *HttpRouter) OidcCallback(context *gin.Context) {
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

func (r *HttpRouter) Logout(context *gin.Context) {
	context.SetCookie("token", "", -1, r.path, "", r.tlsEnabled, true)
	context.SetCookie("token_error", "", -1, r.path, "", r.tlsEnabled, true)
}

func (r *HttpRouter) handleTokenError(context *gin.Context, err string) {
	r.setCookieTokenError(context, err)
	path := r.path + "?error=" + err
	http.Redirect(context.Writer, context.Request, path, http.StatusFound)
}

func (r *HttpRouter) setCookieSession(context *gin.Context, value string) {
	// NOTE: maxAge is provided in seconds 2592000 sec = 30 days
	context.SetCookie("session", value, 2592000, r.path, "", r.tlsEnabled, true)
}

func (r *HttpRouter) setCookieToken(context *gin.Context, value string, exp *time.Time) {
	// NOTE: we add 60 sec just to see proper error when token has expired
	seconds := int(time.Until(*exp).Seconds()) + 60
	context.SetSameSite(http.SameSiteStrictMode)
	context.SetCookie("token_error", "", -1, r.path, "", r.tlsEnabled, true)
	context.SetCookie("token", value, seconds, r.path, "", r.tlsEnabled, true)

}

func (r *HttpRouter) setCookieTokenError(context *gin.Context, value string) {
	context.SetSameSite(http.SameSiteStrictMode)
	context.SetCookie("token", "", -1, r.path, "", r.tlsEnabled, true)
	context.SetCookie("token_error", value, 0, r.path, "", r.tlsEnabled, true)
}

func (r *HttpRouter) setCookieState(context *gin.Context, value string) {
	path := ""
	if r.path != "/" {
		path = r.path
	}
	context.SetCookie("state", value, 600, path+"/api/oidc", "", r.tlsEnabled, true)
}

// PostAuthToken hands the caller a fresh session token for a client that cannot
// carry Ivory's cookie - an MCP client reaching /api/mcp. It needs no permission
// for the same reason changing your own password does not: it is the account you
// are already signed in as. It deliberately writes no cookie, so minting a token
// for an agent leaves the browser's own session exactly as it was.
func (r *HttpRouter) PostAuthToken(context *gin.Context) {
	username := context.GetString(config.AuthContextKey.Username)
	callerType := AuthType(context.GetString(config.AuthContextKey.Type))
	var request RouterTokenRequest
	if errBind := context.ShouldBindJSON(&request); errBind != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errBind.Error()})
		return
	}

	expiration := time.Duration(request.ExpiresIn) * time.Second
	token, exp, errToken := r.authService.GenerateSessionToken(username, callerType, expiration)
	if errToken != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errToken.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"response": RouterTokenResponse{Token: token, ExpiresAt: *exp}})
}
