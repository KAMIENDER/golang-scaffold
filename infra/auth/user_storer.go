package auth

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/KAMIENDER/golang-scaffold/infra/database/nosql"
	"github.com/KAMIENDER/golang-scaffold/infra/persistent/basic"
	"github.com/KAMIENDER/golang-scaffold/infra/persistent/po"
	"github.com/pkg/errors"
	"github.com/volatiletech/authboss/v3"
	"github.com/volatiletech/authboss/v3/otp/twofactor/sms2fa"
	"github.com/volatiletech/authboss/v3/otp/twofactor/totp2fa"
	"gorm.io/gorm"
)

var flagDebug = flag.Bool("debug", false, "output debugging information")

func debugln(args ...interface{}) {
	if !*flagDebug {
		return
	}

	fmt.Println(args...)
}

func debugf(format string, args ...interface{}) {
	if !*flagDebug {
		return
	}

	fmt.Printf(format, args...)
}

// User struct for authboss
type User struct {
	*po.User
	// ID int

	// // Non-authboss related field
	// Name string

	// // Auth
	// Email    string
	// Password string

	// // Confirm
	// ConfirmSelector string
	// ConfirmVerifier string
	// Confirmed       bool

	// // Lock
	// AttemptCount int
	// LastAttempt  time.Time
	// Locked       time.Time

	// // Recover
	// RecoverSelector    string
	// RecoverVerifier    string
	// RecoverTokenExpiry time.Time

	// // OAuth2
	// OAuth2UID          string
	// OAuth2Provider     string
	// OAuth2AccessToken  string
	// OAuth2RefreshToken string
	// OAuth2Expiry       time.Time

	// // 2fa
	// TOTPSecretKey      string
	// SMSPhoneNumber     string
	// SMSSeedPhoneNumber string
	// RecoveryCodes      string

	// Remember is in another table
}

// This pattern is useful in real code to ensure that
// we've got the right interfaces implemented.
var (
	assertUser   = &User{}
	assertStorer = &DBStorer{}

	_ authboss.User            = assertUser
	_ authboss.AuthableUser    = assertUser
	_ authboss.ConfirmableUser = assertUser
	_ authboss.LockableUser    = assertUser
	_ authboss.RecoverableUser = assertUser
	_ authboss.ArbitraryUser   = assertUser

	_ totp2fa.User = assertUser
	_ sms2fa.User  = assertUser

	_ authboss.CreatingServerStorer    = assertStorer
	_ authboss.ConfirmingServerStorer  = assertStorer
	_ authboss.RecoveringServerStorer  = assertStorer
	_ authboss.RememberingServerStorer = assertStorer
)

// PutPID into user
func (u *User) PutPID(pid string) { u.Email = pid }

// PutPassword into user
func (u *User) PutPassword(password string) { u.Password = password }

// PutEmail into user
func (u *User) PutEmail(email string) { u.Email = email }

// PutConfirmed into user
func (u *User) PutConfirmed(confirmed bool) { u.Confirmed = confirmed }

// PutConfirmSelector into user
func (u *User) PutConfirmSelector(confirmSelector string) { u.ConfirmSelector = confirmSelector }

// PutConfirmVerifier into user
func (u *User) PutConfirmVerifier(confirmVerifier string) { u.ConfirmVerifier = confirmVerifier }

// PutLocked into user
func (u *User) PutLocked(locked time.Time) { u.Locked = &locked }

// PutAttemptCount into user
func (u *User) PutAttemptCount(attempts int) { u.AttemptCount = attempts }

// PutLastAttempt into user
func (u *User) PutLastAttempt(last time.Time) { u.LastAttempt = &last }

// PutRecoverSelector into user
func (u *User) PutRecoverSelector(token string) { u.RecoverSelector = token }

// PutRecoverVerifier into user
func (u *User) PutRecoverVerifier(token string) { u.RecoverVerifier = token }

// PutRecoverExpiry into user
func (u *User) PutRecoverExpiry(expiry time.Time) { u.RecoverTokenExpiry = &expiry }

// PutTOTPSecretKey into user
func (u *User) PutTOTPSecretKey(key string) { u.TOTPSecretKey = key }

// PutSMSPhoneNumber into user
func (u *User) PutSMSPhoneNumber(key string) { u.SMSPhoneNumber = key }

// PutRecoveryCodes into user
func (u *User) PutRecoveryCodes(key string) { u.RecoveryCodes = key }

// PutOAuth2UID into user
func (u *User) PutOAuth2UID(uid string) { u.OAuth2UID = uid }

// PutOAuth2Provider into user
func (u *User) PutOAuth2Provider(provider string) { u.OAuth2Provider = provider }

// PutOAuth2AccessToken into user
func (u *User) PutOAuth2AccessToken(token string) { u.OAuth2AccessToken = token }

// PutOAuth2RefreshToken into user
func (u *User) PutOAuth2RefreshToken(refreshToken string) { u.OAuth2RefreshToken = refreshToken }

// PutOAuth2Expiry into user
func (u *User) PutOAuth2Expiry(expiry time.Time) { u.OAuth2Expiry = &expiry }

// PutArbitrary into user
func (u *User) PutArbitrary(values map[string]string) {
	if n, ok := values["name"]; ok {
		u.Name = n
	}
	// Add user information: deal with your register
	if n, ok := values["bard_ID"]; ok {
		u.BardID, _ = strconv.Atoi(n)
	}
}

// GetPID from user
func (u User) GetPID() string { return u.Email }

// GetPassword from user
func (u User) GetPassword() string { return u.Password }

// GetEmail from user
func (u User) GetEmail() string { return u.Email }

// GetConfirmed from user
func (u User) GetConfirmed() bool { return u.Confirmed }

// GetConfirmSelector from user
func (u User) GetConfirmSelector() string { return u.ConfirmSelector }

// GetConfirmVerifier from user
func (u User) GetConfirmVerifier() string { return u.ConfirmVerifier }

// GetLocked from user
func (u User) GetLocked() time.Time { return *u.Locked }

// GetAttemptCount from user
func (u User) GetAttemptCount() int { return u.AttemptCount }

// GetLastAttempt from user
func (u User) GetLastAttempt() time.Time { return *u.LastAttempt }

// GetRecoverSelector from user
func (u User) GetRecoverSelector() string { return u.RecoverSelector }

// GetRecoverVerifier from user
func (u User) GetRecoverVerifier() string { return u.RecoverVerifier }

// GetRecoverExpiry from user
func (u User) GetRecoverExpiry() time.Time { return *u.RecoverTokenExpiry }

// GetTOTPSecretKey from user
func (u User) GetTOTPSecretKey() string { return u.TOTPSecretKey }

// GetSMSPhoneNumber from user
func (u User) GetSMSPhoneNumber() string { return u.SMSPhoneNumber }

// GetSMSPhoneNumberSeed from user
func (u User) GetSMSPhoneNumberSeed() string { return u.SMSSeedPhoneNumber }

// GetRecoveryCodes from user
func (u User) GetRecoveryCodes() string { return u.RecoveryCodes }

// IsOAuth2User returns true if the user was created with oauth2
func (u User) IsOAuth2User() bool { return len(u.OAuth2UID) != 0 }

// GetOAuth2UID from user
func (u User) GetOAuth2UID() (uid string) { return u.OAuth2UID }

// GetOAuth2Provider from user
func (u User) GetOAuth2Provider() (provider string) { return u.OAuth2Provider }

// GetOAuth2AccessToken from user
func (u User) GetOAuth2AccessToken() (token string) { return u.OAuth2AccessToken }

// GetOAuth2RefreshToken from user
func (u User) GetOAuth2RefreshToken() (refreshToken string) { return u.OAuth2RefreshToken }

// GetOAuth2Expiry from user
func (u User) GetOAuth2Expiry() (expiry time.Time) { return *u.OAuth2Expiry }

// GetArbitrary from user
func (u User) GetArbitrary() map[string]string {
	return map[string]string{
		"name": u.Name,
		// Add user information: Processing your custom fields during registration.
		"bard_ID": strconv.Itoa(u.BardID),
	}
}

// DBStorer stores users in memory
type DBStorer struct {
	db   *gorm.DB
	noDB nosql.NoSQLDB
}

// NewDBStorer constructor
func NewDBStorer(db *gorm.DB, noDB nosql.NoSQLDB) *DBStorer {
	return &DBStorer{
		db:   db.Debug(),
		noDB: noDB,
	}
}

// Save the user
func (m *DBStorer) Save(ctx context.Context, user authboss.User) error {
	db := basic.Use(m.db).User
	u := user.(*User)
	return db.WithContext(ctx).Save(u.User)
}

// SaveOAuth2 user
func (m DBStorer) SaveOAuth2(ctx context.Context, user authboss.OAuth2User) error {
	db := basic.Use(m.db).User
	u := user.(*User)
	return db.WithContext(ctx).Save(u.User)
}

// Load the user
func (m DBStorer) Load(ctx context.Context, key string) (authboss.User, error) {
	// Check to see if our key is actually an oauth2 pid
	_, uid, err := authboss.ParseOAuth2PID(key)
	user := basic.Use(m.db).User
	if err == nil {
		out, err := user.WithContext(ctx).Where(user.OAuth2UID.Eq(uid)).Find()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) || len(out) == 0 {
			return nil, authboss.ErrUserFound
		}
		return &User{
			out[0],
		}, err
	}
	out, err := user.WithContext(ctx).Where(user.Email.Eq(key)).Find()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) || len(out) == 0 {
		return nil, authboss.ErrUserFound
	}
	return &User{
		out[0],
	}, err
}

// New user creation
func (m DBStorer) New(_ context.Context) authboss.User {
	return &User{
		&po.User{},
	}
}

// Create the user
func (m DBStorer) Create(ctx context.Context, user authboss.User) error {
	db := basic.Use(m.db).User
	u := user.(*User)
	return db.WithContext(ctx).Save(u.User)
}

// LoadByConfirmSelector looks a user up by confirmation token
func (m DBStorer) LoadByConfirmSelector(ctx context.Context, selector string) (authboss.ConfirmableUser, error) {
	user := basic.Use(m.db).User
	out, err := user.WithContext(ctx).Where(user.ConfirmSelector.Eq(selector)).Find()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) || len(out) == 0 {
		return nil, authboss.ErrUserFound
	}
	return &User{
		out[0],
	}, err
}

// LoadByRecoverSelector looks a user up by confirmation selector
func (m DBStorer) LoadByRecoverSelector(ctx context.Context, selector string) (authboss.RecoverableUser, error) {
	user := basic.Use(m.db).User
	out, err := user.WithContext(ctx).Where(user.RecoverSelector.Eq(selector)).Find()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) || len(out) == 0 {
		return nil, authboss.ErrUserFound
	}
	return &User{
		out[0],
	}, err
}

func buildRememberTokenKey(pid string) string {
	return "remember_token_key_" + pid
}

// AddRememberToken to a user
func (m DBStorer) AddRememberToken(ctx context.Context, pid, token string) error {
	oriTokens := make([]string, 0)
	_, err := m.noDB.Get(ctx, buildRememberTokenKey(pid), &oriTokens)
	if err != nil {
		return errors.Wrap(err, "")
	}
	oriTokens = append(oriTokens, token)

	m.noDB.Set(ctx, buildRememberTokenKey(pid), oriTokens, -1)
	return nil
}

// DelRememberTokens removes all tokens for the given pid
func (m DBStorer) DelRememberTokens(ctx context.Context, pid string) error {
	return m.noDB.Del(ctx, buildRememberTokenKey(pid))
}

// UseRememberToken finds the pid-token pair and deletes it.
// If the token could not be found return ErrTokenNotFound
func (m DBStorer) UseRememberToken(ctx context.Context, pid, token string) error {
	tokens := make([]string, 0)
	ok, err := m.noDB.Get(ctx, pid, &tokens)
	if err != nil {
		return err
	}
	if !ok {
		return authboss.ErrTokenNotFound
	}

	for i, tok := range tokens {
		if tok == token {
			tokens[len(tokens)-1] = tokens[i]
			tokens = tokens[:len(tokens)-1]
			return m.noDB.Set(ctx, buildRememberTokenKey(pid), tokens, -1)
		}
	}

	return authboss.ErrTokenNotFound
}

// NewFromOAuth2 creates an oauth2 user (but not in the database, just a blank one to be saved later)
// func (m DBStorer) NewFromOAuth2(_ context.Context, provider string, details map[string]string) (authboss.OAuth2User, error) {
// 	switch provider {
// 	case "google":
// 		email := details[aboauth.OAuth2Email]

// 		var user *User
// 		if u, ok := m.Users[email]; ok {
// 			user = &u
// 		} else {
// 			user = &User{}
// 		}

// 		// Google OAuth2 doesn't allow us to fetch real name without more complicated API calls
// 		// in order to do this properly in your own app, look at replacing the authboss oauth2.GoogleUserDetails
// 		// method with something more thorough.
// 		user.Name = "Unknown"
// 		user.Email = details[aboauth.OAuth2Email]
// 		user.OAuth2UID = details[aboauth.OAuth2UID]
// 		user.Confirmed = true

// 		return user, nil
// 	}

// 	return nil, errors.Errorf("unknown provider %s", provider)
// }

/*
func (s DBStorer) PutOAuth(uid, provider string, attr authboss.Attributes) error {
	return s.Create(uid+provider, attr)
}

func (s DBStorer) GetOAuth(uid, provider string) (result interface{}, err error) {
	user, ok := s.Users[uid+provider]
	if !ok {
		return nil, authboss.ErrUserNotFound
	}

	return &user, nil
}

func (s DBStorer) AddToken(key, token string) error {
	s.Tokens[key] = append(s.Tokens[key], token)
	fmt.Println("AddToken")
	spew.Dump(s.Tokens)
	return nil
}

func (s DBStorer) DelTokens(key string) error {
	delete(s.Tokens, key)
	fmt.Println("DelTokens")
	spew.Dump(s.Tokens)
	return nil
}

func (s DBStorer) UseToken(givenKey, token string) error {
	toks, ok := s.Tokens[givenKey]
	if !ok {
		return authboss.ErrTokenNotFound
	}

	for i, tok := range toks {
		if tok == token {
			toks[i], toks[len(toks)-1] = toks[len(toks)-1], toks[i]
			s.Tokens[givenKey] = toks[:len(toks)-1]
			return nil
		}
	}

	return authboss.ErrTokenNotFound
}

func (s DBStorer) ConfirmUser(tok string) (result interface{}, err error) {
	fmt.Println("==============", tok)

	for _, u := range s.Users {
		if u.ConfirmToken == tok {
			return &u, nil
		}
	}

	return nil, authboss.ErrUserNotFound
}

func (s DBStorer) RecoverUser(rec string) (result interface{}, err error) {
	for _, u := range s.Users {
		if u.RecoverToken == rec {
			return &u, nil
		}
	}

	return nil, authboss.ErrUserNotFound
}
*/
