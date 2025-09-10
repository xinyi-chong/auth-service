package auth

import (
	authconsts "auth-service/internal/shared/consts"
	token "auth-service/pkg/jwt"
	"github.com/google/uuid"
	"github.com/xinyi-chong/common-lib/consts"
	apperrors "github.com/xinyi-chong/common-lib/errors"
	"github.com/xinyi-chong/common-lib/response"
	"github.com/xinyi-chong/common-lib/success"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service Service
	logger  *zap.Logger
}

func NewController(service Service, logger *zap.Logger) *Controller {
	return &Controller{service: service, logger: logger}
}

// Register godoc
// @Summary Register a new user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param account body RegisterParam true "User Account"
// @Success 201 {object} response.Response "Created"
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /auth/register [post]
func (ctrl *Controller) Register(c *gin.Context) {
	var param RegisterParam
	if err := c.ShouldBindJSON(&param); err != nil {
		ctrl.logger.Debug("Invalid request payload", zap.Error(err))
		response.Error(c, apperrors.ErrBadRequest)
		return
	}

	ctx := c.Request.Context()
	err := ctrl.service.Register(ctx, param)
	if err != nil {
		ctrl.logger.Error("Register error", zap.String("username", *param.Username), zap.String("email", param.Email), zap.Error(err))
		response.Error(c, err, apperrors.ErrRegistrationFailed)
		return
	}

	response.Success(c, success.Registered, nil)
}

// Login godoc
// @Summary Login
// @Description Login
// @Tags Authentication
// @Accept json
// @Produce json
// @Param account body LoginParam true "Credentials"
// @Success 200 {object} response.Response{data=LoginResponse} "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "User not found"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /auth/login [post]
func (ctrl *Controller) Login(c *gin.Context) {
	var param LoginParam
	if err := c.ShouldBindJSON(&param); err != nil {
		ctrl.logger.Debug("Invalid request payload", zap.Error(err))
		response.Error(c, apperrors.ErrIncorrectX.WithField(consts.EmailOrPasswordField))
		return
	}

	ctx := c.Request.Context()
	resp, err := ctrl.service.Login(ctx, param.Email, param.Password)
	if err != nil {
		ctrl.logger.Error("Login error", zap.String("email", param.Email), zap.Error(err))
		response.Error(c, err, apperrors.ErrLoginFailed)
		return
	}

	setRefreshTokenCookie(c, resp.RefreshToken)

	response.Success(c, success.LoggedIn, resp)
}

// ChangePassword godoc
// @Summary Change Password
// @Description Change Password
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param body body ChangePasswordParam true "Change Password"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /auth/change-password [patch]
func (ctrl *Controller) ChangePassword(c *gin.Context) {
	userIDValue, exists := c.Get(consts.CtxUserID)
	if !exists {
		ctrl.logger.Debug("Missing user ID in context")
		response.Error(c, apperrors.ErrUnauthorized)
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		ctrl.logger.Error("Invalid user ID type in context", zap.Any("user_id", userIDValue))
		response.Error(c, apperrors.ErrInternalServerError)
		return
	}

	var req ChangePasswordParam
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.logger.Debug("Invalid request payload", zap.Error(err))
		response.Error(c, apperrors.ErrBadRequest)
		return
	}

	ctx := c.Request.Context()
	err := ctrl.service.ChangePassword(ctx, userID, req.OldPassword, req.NewPassword)
	if err != nil {
		ctrl.logger.Error("ChangePassword error", zap.Error(err))
		response.Error(c, err)
		return
	}

	response.Success(c, success.XChanged.WithField(consts.PasswordField), nil)
}

// RefreshToken godoc
// @Summary Refresh Token
// @Description Refresh Token
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=Tokens} "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /auth/refresh [post]
func (ctrl *Controller) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if refreshToken == "" {
		if err != nil {
			ctrl.logger.Debug("Cookie missing", zap.Error(err))
		}
		response.Error(c, apperrors.ErrUnauthorized)
		return
	}

	ctx := c.Request.Context()
	tokens, err := ctrl.service.RefreshToken(ctx, refreshToken)
	if err != nil {
		ctrl.logger.Error("RefreshToken error", zap.Error(err))
		response.Error(c, err, apperrors.ErrSessionExpired)
		return
	}

	setRefreshTokenCookie(c, tokens.RefreshToken)

	response.Success(c, success.SessionRefreshed, Tokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// Logout godoc
// @Summary Logout
// @Description Logout
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /auth/logout [post]
func (ctrl *Controller) Logout(c *gin.Context) {
	clearRefreshTokenCookie(c)

	ctx := c.Request.Context()

	refreshToken, _ := c.Cookie(authconsts.CookieRefreshToken)
	err := token.InvalidateToken(ctx, refreshToken)
	if err != nil {
		ctrl.logger.Error("Invalidate Refresh Token error", zap.Error(err))
		response.Error(c, apperrors.ErrSessionExpired)
		return
	}

	accessToken := c.GetString(consts.CtxAccessToken)
	err = token.InvalidateToken(ctx, accessToken)
	if err != nil {
		ctrl.logger.Error("Invalidate Access Token error", zap.Error(err))
		response.Error(c, apperrors.ErrSessionExpired)
		return
	}

	response.Success(c, success.LoggedOut, nil)
}
