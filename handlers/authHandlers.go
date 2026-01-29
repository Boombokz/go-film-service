package handlers

import (
	"filmservice/config"
	"filmservice/models"
	"filmservice/repositories"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandlers struct {
	usersRepo *repositories.UsersRepository
}

func NewAuthHandlers(usersRepo *repositories.UsersRepository) *AuthHandlers {
	return &AuthHandlers{usersRepo: usersRepo}
}

type signInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signInResponse struct {
	Token string `json:"token"`
}

// SignIn   	 godoc
// @Summary      Sign in
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body signInRequest true "Sign in payload"
// @Success      200  {object}  signInResponse "OK"
// @Success      400  {object}  models.ApiError "Invalid request payload"
// @Failure      401  {object}  models.ApiError "Invalid credentials"
// @Failure      500  {object}  models.ApiError
// @Router       /auth/signIn [post]
func (h *AuthHandlers) SignIn(c *gin.Context) {
	var request signInRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid request payload"))
		return
	}

	user, err := h.usersRepo.FindByEmail(c, request.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.NewApiError("Invalid credentials"))
		return
	}

	claims := jwt.RegisteredClaims{
		Subject:   strconv.Itoa(user.Id),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.Config.JwtExpiresIn)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Config.JwtSecretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("Could not generate JWT token"))
		return
	}

	c.JSON(http.StatusOK, signInResponse{
		Token: tokenString,
	})
}

func (h *AuthHandlers) SignOut(c *gin.Context) {
	// TODO: Delete token
	c.Status(http.StatusOK)
}
