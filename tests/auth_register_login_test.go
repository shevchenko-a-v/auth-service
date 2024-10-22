package tests

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt"
	"github.com/shevchenko-a-v/auth-service/tests/suite"
	ssov1 "github.com/shevchenko-a-v/protofiles/gen/go/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyAppID = 0
	appID      = 1
	appSecret  = "test-secret"

	passDefaultLen = 10
)

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}

func TestRegisterLogin_LoginHappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword()

	respRegister, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respRegister.GetUserId())

	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})
	require.NoError(t, err)

	loginTime := time.Now()

	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, respRegister.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))

	const deltaSeconds = 1
	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func TestRegisterLogin_RegisterExistingUser(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword()

	respRegister, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respRegister.GetUserId())

	respRegister, err = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})

	require.Error(t, err)
	assert.Empty(t, respRegister)
	assert.ErrorContains(t, err, "user already exists")
}

func TestRegisterLogin_LoginUserNotExists(t *testing.T) {
	ctx, st := suite.New(t)

	email := "notexistingemail@in.database"
	password := randomFakePassword()

	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})
	require.Error(t, err)
	assert.Empty(t, respLogin)
	assert.ErrorContains(t, err, "invalid user or password")
}

func TestRegisterLogin_LoginWrongPassword(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword()

	respRegister, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respRegister.GetUserId())

	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password + "wrongpass",
		AppId:    appID,
	})
	require.Error(t, err)
	assert.Empty(t, respLogin)
	assert.ErrorContains(t, err, "invalid user or password")
}

func TestRegisterLogin_LoginWrongApp(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword()

	respRegister, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respRegister.GetUserId())

	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    777,
	})
	require.Error(t, err)
	assert.Empty(t, respLogin)
	assert.ErrorContains(t, err, "invalid application id")
}

func TestRegister_FailCases(t *testing.T) {
	ctx, st := suite.New(t)
	tests := []struct {
		name        string
		email       string
		password    string
		expectedErr string
	}{
		{
			name:        "Empty password",
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "empty password",
		},
		{
			name:        "Empty email",
			email:       "",
			password:    randomFakePassword(),
			expectedErr: "empty email",
		},
		{
			name:        "Empty email and password",
			email:       "",
			password:    "",
			expectedErr: "empty email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respRegister, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    tt.email,
				Password: tt.password,
			})

			require.Error(t, err)
			assert.Empty(t, respRegister)
			assert.ErrorContains(t, err, tt.expectedErr)
		})
	}
}

func TestLogin_FailCases(t *testing.T) {
	ctx, st := suite.New(t)
	tests := []struct {
		name        string
		email       string
		password    string
		appID       int
		expectedErr string
	}{
		{
			name:        "Empty password",
			email:       gofakeit.Email(),
			password:    "",
			appID:       appID,
			expectedErr: "empty password",
		},
		{
			name:        "Empty email",
			email:       "",
			password:    randomFakePassword(),
			appID:       appID,
			expectedErr: "empty email",
		},
		{
			name:        "Empty email and password",
			email:       "",
			password:    "",
			appID:       appID,
			expectedErr: "empty email",
		},
		{
			name:        "Empty app id",
			email:       gofakeit.Email(),
			password:    randomFakePassword(),
			appID:       emptyAppID,
			expectedErr: "empty app_id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
				AppId:    int32(tt.appID),
			})
			require.Error(t, err)
			assert.Empty(t, respLogin)
			assert.ErrorContains(t, err, tt.expectedErr)
		})
	}
}
