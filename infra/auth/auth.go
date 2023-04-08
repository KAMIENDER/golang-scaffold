package auth

import (
	"encoding/base64"
	"net/http"
	"regexp"
	"time"

	"github.com/KAMIENDER/golang-scaffold/infra/config"
	"github.com/KAMIENDER/golang-scaffold/infra/database/nosql"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	abclientstate "github.com/volatiletech/authboss-clientstate"
	abrenderer "github.com/volatiletech/authboss-renderer"
	"github.com/volatiletech/authboss/v3"
	boss "github.com/volatiletech/authboss/v3/auth"
	"github.com/volatiletech/authboss/v3/confirm"
	"github.com/volatiletech/authboss/v3/defaults"
	"github.com/volatiletech/authboss/v3/otp/twofactor"
	"github.com/volatiletech/authboss/v3/otp/twofactor/sms2fa"
	"github.com/volatiletech/authboss/v3/otp/twofactor/totp2fa"
	"github.com/volatiletech/authboss/v3/recover"
	"github.com/volatiletech/authboss/v3/register"
	"gorm.io/gorm"
)

const (
	sessionCookieName = "session"
)

var (
	_ register.Register
	_ boss.Auth
	_ confirm.Confirm
	_ recover.Recover
)

type AuthManager struct {
	ab *authboss.Authboss
}

// NewAuthManager You can customize the settings of Auth according to your own needs.
func NewAuthManager(conf *config.Config, db *gorm.DB, nosqlDB nosql.NoSQLDB) (*AuthManager, error) {
	authDataBase := NewDBStorer(db, nosqlDB)
	var (
		sessionStore abclientstate.SessionStorer
		cookieStore  abclientstate.CookieStorer
	)

	cookieStoreKey, _ := base64.StdEncoding.DecodeString(`NpEPi8pEjKVjLGJ6kYCS+VTCzi6BUuDzU0wrwXyf5uDPArtlofn2AG6aTMiPmN3C909rsEWMNqJqhIVPGP3Exg==`)
	sessionStoreKey, _ := base64.StdEncoding.DecodeString(`AbfYwmmt8UCwUuhd9qvfNA9UCuN1cVcKJN1ofbiky6xCyyBj20whe40rJa3Su0WOWLWcPpO1taqJdsEI/65+JA==`)
	cookieStore = abclientstate.NewCookieStorer(cookieStoreKey, nil)
	cookieStore.HTTPOnly = false
	cookieStore.Secure = false
	sessionStore = abclientstate.NewSessionStorer(sessionCookieName, sessionStoreKey, nil)
	cstore := sessionStore.Store.(*sessions.CookieStore)
	cstore.Options.HttpOnly = false
	cstore.Options.Secure = false
	cstore.MaxAge(int((30 * 24 * time.Hour) / time.Second))

	ab := authboss.New()
	ab.Config.Paths.RootURL = "http://localhost:8080"
	ab.Config.Modules.LogoutMethod = "GET"
	ab.Config.Core.ViewRenderer = defaults.JSONRenderer{}
	ab.Config.Storage.Server = authDataBase
	ab.Config.Storage.SessionState = sessionStore
	ab.Config.Storage.CookieState = cookieStore
	ab.Config.Core.MailRenderer = abrenderer.NewEmail("/auth", "auth_template")
	ab.Config.Modules.RegisterPreserveFields = []string{"email", "name"}

	ab.Config.Modules.TOTP2FAIssuer = "ABBlog"
	ab.Config.Modules.ResponseOnUnauthed = authboss.RespondRedirect

	ab.Config.Modules.TwoFactorEmailAuthRequired = true

	// This instantiates and uses every default implementation
	// in the Config.Core area that exist in the defaults package.
	// Just a convenient helper if you don't want to do anything fancy.
	defaults.SetCore(&ab.Config, true, false)

	// email sender config
	ab.Config.Mail.From = "ones_kami_sama@qq.com"
	ab.Config.Mail.FromName = "bard"
	ab.Config.Core.Mailer = NewEmailSender(conf)

	// Here we initialize the bodyreader as something customized in order to accept a name
	// parameter for our user as well as the standard e-mail and password.
	//
	// We also change the validation for these fields
	// to be something less secure so that we can use test data easier.
	emailRule := defaults.Rules{
		FieldName: "email", Required: true,
		MatchError: "Must be a valid e-mail address",
		MustMatch:  regexp.MustCompile(`.*@.*\.[a-z]+`),
	}
	passwordRule := defaults.Rules{
		FieldName: "password", Required: true,
		MinLength: 4,
	}
	nameRule := defaults.Rules{
		FieldName: "name", Required: true,
		MinLength: 2,
	}
	// You can add custom information that you need during registration.
	bardIDRule := defaults.Rules{
		FieldName: "bard_ID", Required: true,
		MustMatch:  regexp.MustCompile("\\d{5,}"),
		MatchError: "错误的bard id格式",
	}

	ab.Config.Core.BodyReader = defaults.HTTPBodyReader{
		ReadJSON: true,
		Rulesets: map[string][]defaults.Rules{
			"register":    {emailRule, passwordRule, nameRule, bardIDRule},
			"recover_end": {passwordRule},
		},
		Confirms: map[string][]string{
			"register":    {"password", authboss.ConfirmPrefix + "password"},
			"recover_end": {"password", authboss.ConfirmPrefix + "password"},
		},
		Whitelist: map[string][]string{
			"register": {"email", "name", "password", "bard_ID"},
		},
	}

	// Set up 2fa
	twofaRecovery := &twofactor.Recovery{Authboss: ab}
	if err := twofaRecovery.Setup(); err != nil {
		return nil, err
	}

	totp := &totp2fa.TOTP{Authboss: ab}
	if err := totp.Setup(); err != nil {
		return nil, err
	}

	sms := &sms2fa.SMS{Authboss: ab, Sender: SMSLogSender{}}
	if err := sms.Setup(); err != nil {
		return nil, err
	}

	if err := ab.Init(); err != nil {
		return nil, err
	}
	return &AuthManager{
		ab: ab,
	}, nil
}

// SetupAuthBoss Used to set the path prefix used by Auth.
func (m *AuthManager) SetupAuthBoss(r *gin.Engine) {
	r.Any("/auth/*w", gin.WrapH(http.StripPrefix("/auth", m.ab.LoadClientStateMiddleware(m.ab.Config.Core.Router))))
}

type nextRequestHandler struct {
	c *gin.Context
}

func (h *nextRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.c.Next()
}

// WrapHandler Add Authboss validation middleware for the handler.
func (m *AuthManager) WrapHandler(handler gin.HandlerFunc) gin.HandlerFunc {
	// The execution order is consistent with the order in the list.
	authMiddleware := []func(h http.Handler) http.Handler{
		m.ab.LoadClientStateMiddleware,
		authboss.Middleware(m.ab, true, false, false),
	}
	return func(c *gin.Context) {
		if len(authMiddleware) > 0 {
			now := authMiddleware[len(authMiddleware)-1](&nextRequestHandler{c: c})
			for t := len(authMiddleware) - 2; t > -1; t-- {
				now = authMiddleware[t](now)
			}
			now.ServeHTTP(m.ab.NewResponse(c.Writer), c.Request)
		}
		handler(c)
	}
}
