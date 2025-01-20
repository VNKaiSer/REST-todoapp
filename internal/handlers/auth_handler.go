package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"todo-app/bunapp"
	"todo-app/httputil/httperror"
	"todo-app/httputil/httpresponse"
	"todo-app/internal/db"
	"todo-app/internal/dtos"
	handlers "todo-app/internal/services"
	"todo-app/pkg/utils"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"github.com/twinj/uuid"
)

type AuthHandler struct {
	app *bunapp.App
}

var _ handlers.AuthHandlerService = (*AuthHandler)(nil)

func NewAuthHandler(app *bunapp.App) *AuthHandler {
	return &AuthHandler{app: app}
}


// Login implements handlers.AuthHandlerService.
func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Body)
}

// RefreshToken implements handlers.AuthHandlerService.
func (a *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// Register implements handlers.AuthHandlerService.
func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var authDTO dtos.AuthDTO
	err := json.NewDecoder(r.Body).Decode(&authDTO)
	if err != nil {
		render.Render(w, r, httperror.ErrInvalidRequest(err))
		return
	}

	exists, err := a.app.DB().NewSelect().Model((*db.User)(nil)).Where("username = ?", authDTO.Username).Exists(r.Context())
	if err != nil {
		render.Render(w, r, httperror.ErrInternalError(err))
		return
	}

	if exists {
		render.Render(w, r, httperror.ErrBadRequest(fmt.Errorf("user already exists")))
		return
	}

	passwordHash, err := utils.HashPassword(authDTO.Password)
	if err != nil {
		render.Render(w, r, httperror.ErrInternalError(err))
		return
	}

	user := &db.User{
        Username:     authDTO.Username,
        PasswordHash: passwordHash,
    }
    _, err = a.app.DB().NewInsert().Model(user).Returning("*").Exec(r.Context())
	
	if err != nil {
		render.Render(w, r, httperror.ErrInternalError(err))
		return
	}

	Token, RefreshToken, err := a.NewJWT().GenerateTokenPair(authDTO.Username, user.ID, false)
	if err != nil {
		render.Render(w, r, httperror.ErrInternalError(err))
		return
	}
	// save db
    _, err = a.app.DB().NewInsert().Model(&db.Session{
        AccessToken: Token,
        RefreshToken: RefreshToken,
        UserID: user.ID,
    }).Exec(r.Context())

	if err != nil {
		render.Render(w, r, httperror.ErrInternalError(err))
		return
	}

	render.JSON(w, r, httpresponse.SingleResponse{
		Data: map[string]any{
			"assess_token": Token,
		},
	})
}

type JwtPayload struct {
	Username string `json:"username"`
	Sub          int64               `json:"sub"`
	Exp          int64                `json:"exp"`
	Iat          int64                `json:"iat"`
	IsAnonymous  bool                 `json:"is_anonymous"`
	jwt.RegisteredClaims
}

type JWT struct {
    secretKey string
    refreshSecretKey string
    accessDuration time.Duration
    refreshDuration time.Duration
}

func (a *AuthHandler) NewJWT() *JWT {
    return &JWT{
        secretKey: a.app.Config().Jwt.Secret,
        refreshSecretKey: a.app.Config().Jwt.Secret,
        accessDuration: time.Hour * 24,
        refreshDuration: time.Hour * 24 * 7,
    }
}


func (j *JWT) GenerateTokenPair(username string, uid int64, isAnonymous bool) (accessToken string, refreshToken string, err error) {
    now := time.Now()
    // Create access token
    accessClaims := JwtPayload{
		Username: username,
        Sub:      uid,      
        Exp:      now.Add(time.Duration(j.accessDuration)).Unix(),
        Iat:      now.Unix(),
        IsAnonymous: isAnonymous,
        RegisteredClaims: jwt.RegisteredClaims{
           ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.accessDuration))),
        },
    }

    accessTokenObject := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
    accessToken, err = accessTokenObject.SignedString([]byte(j.secretKey))
    if err != nil {
        return "", "", fmt.Errorf("failed to sign access token: %w", err)
    }

    // Create refresh token
    refreshClaims := JwtPayload{
		Username: username,
        Sub:      uid,
        Exp:      now.Add(time.Duration(j.refreshDuration)).Unix(),
        Iat:      now.Unix(),
        IsAnonymous: isAnonymous,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.refreshDuration))),
            ID: uuid.NewV4().String(),
        },
    }


    refreshTokenObject := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    refreshToken, err = refreshTokenObject.SignedString([]byte(j.refreshSecretKey))
	// save to db
    if err != nil {
        return "", "", fmt.Errorf("failed to sign refresh token: %w", err)
    }
    return accessToken, refreshToken, nil
}

// VerifyAccessToken verify access token
func (j *JWT) VerifyAccessToken(tokenString string) (*JwtPayload, error) {
    claims := &JwtPayload{}

    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
       if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
           return nil, fmt.Errorf("invalid signing method")
       }
       return []byte(j.secretKey), nil
   })
   if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }
   
   return claims, nil
}

// VerifyRefreshToken verify refresh token
func (j *JWT) VerifyRefreshToken(tokenString string) (*JwtPayload, error) {
    claims := &JwtPayload{}

    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
       if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
           return nil, fmt.Errorf("invalid signing method")
       }
       return []byte(j.refreshSecretKey), nil
   })
   if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }
   
   return claims, nil
}