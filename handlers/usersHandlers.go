package handlers

import (
	"filmservice/models"
	"filmservice/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UsersHandlers struct {
	repo *repositories.UsersRepository
}

type createUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type updateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type changeUserPasswordRequest struct {
	Password string `json:"password"`
}

type userResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewUsersHandlers(usersRepo *repositories.UsersRepository) *UsersHandlers {
	return &UsersHandlers{repo: usersRepo}
}

// FindAll   	 godoc
// @Summary      Get all users
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {array}   userResponse "OK"
// @Failure      500  {object}  models.ApiError
// @Router       /users [get]
// @Security Bearer
func (h *UsersHandlers) FindAll(c *gin.Context) {
	users, err := h.repo.FindAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("could not load users"))
		return
	}

	dtos := make([]userResponse, 0, len(users))

	for _, u := range users {
		r := userResponse{
			Id:    u.Id,
			Name:  u.Name,
			Email: u.Email,
		}

		dtos = append(dtos, r)
	}

	c.JSON(http.StatusOK, dtos)
}

// FindById   	 godoc
// @Summary      Get user by id
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  userResponse "OK"
// @Failure      400  {object}  models.ApiError "Invalid User Id"
// @Failure      500  {object}  models.ApiError
// @Router       /users/{id} [get]
// @Security Bearer
func (h *UsersHandlers) FindById(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			models.NewApiError("Invalid User Id"),
		)
		return
	}

	user, err := h.repo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	r := userResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}

	c.JSON(http.StatusOK, r)
}

// Create   	 godoc
// @Summary      Create new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      createUserRequest  true  "Create user payload"
// @Success      200      {object}  object{id=int}     "OK"
// @Failure      400      {object}  models.ApiError   "Invalid payload"
// @Failure      500      {object}  models.ApiError
// @Router       /users [post]
// @Security Bearer
func (h *UsersHandlers) Create(c *gin.Context) {
	var request createUserRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid payload"))
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("Failed hash password"))
		return
	}

	user := models.User{
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: string(passwordHash),
	}

	id, err := h.repo.Create(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("Could not create user"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

// Update   	 godoc
// @Summary      Update user by id
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id       path      int               true  "User ID"
// @Param        request  body      updateUserRequest  true  "Update user payload"
// @Success      200      {object}  nil               "OK"
// @Failure      400      {object}  models.ApiError   "Invalid User Id / Could not update user"
// @Failure      500      {object}  models.ApiError
// @Router       /users/{id} [put]
// @Security Bearer
func (h *UsersHandlers) Update(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid User Id"))
		return
	}

	_, err = h.repo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	var request updateUserRequest
	err = c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Could not update user"))
		return
	}

	user := models.User{
		Id:    id,
		Name:  request.Name,
		Email: request.Email,
	}

	err = h.repo.Update(c, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Could not update user"))
		return
	}

	c.Status(http.StatusOK)
}

// ChangePassword   	 godoc
// @Summary      Change user password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id       path      int                        true  "User ID"
// @Param        request  body      changeUserPasswordRequest  true  "New password payload"
// @Success      200      {object}  nil                        "OK"
// @Failure      400      {object}  models.ApiError           "Invalid User Id / Invalid payload"
// @Failure      500      {object}  models.ApiError
// @Router       /users/{id}/changePassword [patch]
// @Security Bearer
func (h *UsersHandlers) ChangePassword(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid User Id"))
		return
	}

	_, err = h.repo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	var request changeUserPasswordRequest
	err = c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid payload"))
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("Failed hash password"))
		return
	}

	user := models.User{
		Id:           id,
		PasswordHash: string(passwordHash),
	}

	err = h.repo.ChangePassword(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("Could not change user password"))
		return
	}

	c.Status(http.StatusOK)
}

// Delete   	 godoc
// @Summary      Delete user by id
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  nil "OK"
// @Failure      400  {object}  models.ApiError "Invalid User Id"
// @Failure      500  {object}  models.ApiError
// @Router       /users/{id} [delete]
// @Security Bearer
func (h *UsersHandlers) Delete(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid User Id"))
		return
	}

	_, err = h.repo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	err = h.repo.Delete(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

// GetUserInfo   godoc
// @Summary      Get user info by userId
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  userResponse "OK"
// @Success      401  {object}  models.ApiError
// @Failure      500  {object}  models.ApiError
// @Router       /users/userInfo [get]
// @Security Bearer
func (h *UsersHandlers) GetUserInfo(c *gin.Context) {
	userId := c.GetInt("userId")
	user, err := h.repo.FindById(c, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, userResponse{
		Id:    user.Id,
		Email: user.Email,
		Name:  user.Name,
	})
}
