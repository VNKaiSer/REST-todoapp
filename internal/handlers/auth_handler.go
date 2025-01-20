package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"todo-app/bunapp"
	"todo-app/httputil/httperror"
	"todo-app/httputil/httpresponse"
	"todo-app/internal/constants"
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

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// CheckToken implements handlers.AuthHandlerService.
func (a *AuthHandler) CheckToken(w http.ResponseWriter, r *http.Request) {

	claims, ok := r.Context().Value("current_user").(*JwtPayload)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	fmt.Printf("Current user: %+v\n", claims)
	render.JSON(w, r, httpresponse.SingleResponse{
		Message: "success",
		Data:    claims,
		Status:  http.StatusOK,
	})
}

var _ handlers.AuthHandlerService = (*AuthHandler)(nil)

func NewAuthHandler(app *bunapp.App) *AuthHandler {
	return &AuthHandler{app: app}
}

// @Summary User login
// @Description Đăng nhập với username và password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dtos.AuthDTO true "Login request body"
// @Success 200
// @Failure 400 {object} httperror.ErrResponse
// @Router /api/auth/login [post]
func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var authDTO dtos.AuthDTO
	err := json.NewDecoder(r.Body).Decode(&authDTO)
	fmt.Println(authDTO)
	if err != nil {
		render.Render(w, r, httperror.ErrInvalidRequest(err))
		return
	}

	if authDTO.Username == "" || authDTO.Password == "" {
		render.Render(w, r, httperror.ErrInvalidRequest(errors.New("email or password is required")))
		return
	}

	var user db.User
	err = a.app.DB().NewSelect().Model(&user).Where("username = ?", authDTO.Username).Scan(r.Context())
	if err != nil {
		render.Render(w, r, httperror.ErrForbidden(errors.New(err.Error())))
		return
	}

	IsMatch, err := utils.ComparePassword(user.PasswordHash, authDTO.Password)
	if err != nil {
		render.Render(w, r, httperror.ErrForbidden(errors.New("invalid password")))
		return
	}

	if !IsMatch {
		render.Render(w, r, httperror.ErrForbidden(errors.New("invalid password")))
		return
	}

	Token, RefreshToken, err := a.NewJWT().GenerateTokenPair(authDTO.Username, user.ID, false)
	if err != nil {
		render.Render(w, r, httperror.ErrInternalError(err))
		return
	}

	render.JSON(w, r, httpresponse.SingleResponse{
		Message: "success",
		Data: TokenResponse{
			AccessToken:  Token,
			RefreshToken: RefreshToken,
		},
		Status: http.StatusOK,
	})

}

// RefreshToken implements handlers.AuthHandlerService.
// @Summary Refresh token
// @Description Refresh token if access token is expired and generate new access token and refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token request body"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} httperror.ErrResponse
// @Router /api/auth/refresh-token [post]
func (a *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {

	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Render(w, r, httperror.ErrInvalidRequest(err))
		return
	}

	// Kiểm tra nếu Refresh Token không tồn tại
	if req.RefreshToken == "" {
		render.Render(w, r, httperror.ErrInvalidRequest(fmt.Errorf("refresh token is required")))
		return
	}

	// Giải mã và xác thực Refresh Token
	claims, err := a.NewJWT().VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		render.Render(w, r, httperror.ErrUnAuthorized(fmt.Errorf("invalid or expired refresh token")))
		return
	}

	// Kiểm tra userId từ Refresh Token
	userId := claims.Sub
	if userId == 0 {
		render.Render(w, r, httperror.ErrUnAuthorized(fmt.Errorf("invalid user ID in refresh token")))
		return
	}

	// Kiểm tra Refresh Token từ database (nếu bạn lưu refresh token)
	var sessions db.Session
	err = a.app.DB().NewSelect().Model(&sessions).Where("refresh_token = ?", req.RefreshToken).Scan(r.Context())
	if err != nil {
		render.Render(w, r, httperror.ErrInternalError(errors.New("refresh token is not valid")))
		return
	}

	// Tạo Access Token mới
	AccessToken, RefreshToken, err := a.NewJWT().GenerateTokenPair(claims.Username, userId, false)
	if err != nil {
		render.Render(w, r, httperror.ErrInternalError(err))
		return
	}

	// Cập nhật db token mới
	_, err = a.app.DB().NewUpdate().Model(&sessions).Where("refresh_token = ?", req.RefreshToken).Set("access_token = ?", AccessToken).Set("refresh_token = ?", RefreshToken).Exec(r.Context())
	if err != nil {
		render.Render(w, r, httperror.ErrInternalError(err))
		return
	}

	// Trả về token mới
	response := map[string]string{
		"access_token":  AccessToken,
		"refresh_token": RefreshToken,
	}
	render.JSON(w, r, response)
}

// Register implements handlers.AuthHandlerService.
// @Summary User register
// @Description User register todo app
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dtos.AuthDTO true "User register request body"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} httperror.ErrResponse
// @Router /api/auth/register [post]
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
		AccessToken:  Token,
		RefreshToken: RefreshToken,
		UserID:       user.ID,
	}).Exec(r.Context())

	if err != nil {
		render.Render(w, r, httperror.ErrInternalError(err))
		return
	}

	render.JSON(w, r, httpresponse.SingleResponse{
		Data: TokenResponse{
			AccessToken:  Token,
			RefreshToken: RefreshToken,
		},
	})
}

type JwtPayload struct {
	Username    string `json:"username"`
	Sub         int64  `json:"sub"`
	Exp         int64  `json:"exp"`
	Iat         int64  `json:"iat"`
	IsAnonymous bool   `json:"is_anonymous"`
	jwt.RegisteredClaims
}

type JWT struct {
	secretKey        string
	refreshSecretKey string
	accessDuration   time.Duration
	refreshDuration  time.Duration
}

func (a *AuthHandler) NewJWT() *JWT {
	return &JWT{
		secretKey:        a.app.Config().Jwt.Secret,
		refreshSecretKey: a.app.Config().Jwt.Secret,
		accessDuration:   time.Hour * 24,
		refreshDuration:  time.Hour * 24 * 7,
	}
}

func (j *JWT) GenerateTokenPair(username string, uid int64, isAnonymous bool) (accessToken string, refreshToken string, err error) {
	now := time.Now()
	// Create access token
	accessClaims := JwtPayload{
		Username:    username,
		Sub:         uid,
		Exp:         now.Add(time.Duration(j.accessDuration)).Unix(),
		Iat:         now.Unix(),
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
		Username:    username,
		Sub:         uid,
		Exp:         now.Add(time.Duration(j.refreshDuration)).Unix(),
		Iat:         now.Unix(),
		IsAnonymous: isAnonymous,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.refreshDuration))),
			ID:        uuid.NewV4().String(),
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

func (a *AuthHandler) Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		var err error
		if tokenString == "" {
			err = errors.New("no token provided")
			render.Render(w, r, httperror.ErrUnAuthorized(err))
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims, err := a.NewJWT().VerifyAccessToken(tokenString)
		if err != nil {
			render.Render(w, r, httperror.ErrUnAuthorized(err))
			return
		}

		ctx := context.WithValue(r.Context(), constants.CurrentUser, claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
