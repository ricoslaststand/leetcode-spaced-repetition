package controllers

import (
	"leetcode-spaced-repetition/internal/utils"
	"leetcode-spaced-repetition/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService services.AuthService
}

type loginRequestBody struct {
	email    string `json:"email" binding:"required" validation:"email"`
	password string `json:"password" binding:"required"`
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Login godoc
// @Summary      User login
// @Description  Validates credentials and returns an authentication token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  loginRequestBody  true  "Login credentials"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /auth/login [post]
func (c AuthController) Login(ctx *gin.Context) {
	var body loginRequestBody

	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.FormatErrorBody(ctx, http.StatusBadRequest, err.Error())
		return
	}

	isValid, err := c.authService.Login(ctx, body.email, body.password)
	if err != nil {
		utils.FormatErrorBody(ctx, http.StatusInternalServerError, "An internal server error occurred.")
		return
	}

	if !isValid {
		utils.FormatErrorBody(ctx, http.StatusUnauthorized, "Invalid login credentials")
		return
	}
}
