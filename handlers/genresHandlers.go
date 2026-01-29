package handlers

import (
	"filmservice/models"
	"filmservice/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GenreHandler struct {
	genresRepo *repositories.GenresRepository
}

type createGenreRequest struct {
	Title string
}

type updateGenreRequest struct {
	Title string
}

func NewGenreHandler(genreRepo *repositories.GenresRepository) *GenreHandler {
	return &GenreHandler{
		genresRepo: genreRepo,
	}
}

// FindAll   	 godoc
// @Summary      Get all genres
// @Tags         genres
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Genre "OK"
// @Failure      500  {object}  models.ApiError
// @Router       /genres [get]
// @Security Bearer
func (h *GenreHandler) FindAll(c *gin.Context) {
	genres, err := h.genresRepo.FindAll(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, genres)
}

// FindById   	 godoc
// @Summary      Get genre by id
// @Tags         genres
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Genre ID"
// @Success      200  {object}  models.Genre "OK"
// @Failure      400  {object}  models.ApiError "Invalid Genre Id"
// @Failure      404  {object}  models.ApiError "Genre not found"
// @Failure      500  {object}  models.ApiError
// @Router       /genres/{id} [get]
// @Security Bearer
func (h *GenreHandler) FindById(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid Genre Id"))
		return
	}

	genre, err := h.genresRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, genre)
}

// Create   	 godoc
// @Summary      Create genre
// @Tags         genres
// @Accept       json
// @Produce      json
// @Param        request body createGenreRequest true "Genre payload"
// @Success      200  {object}  object{id=int} "OK"
// @Failure      400  {object}  models.ApiError "Invalid payload"
// @Failure      500  {object}  models.ApiError
// @Router       /genres [post]
// @Security Bearer
func (h *GenreHandler) Create(c *gin.Context) {
	var request createGenreRequest

	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Could not bind JSON"))
		return
	}

	genre := models.Genre{
		Title: request.Title,
	}

	id, err := h.genresRepo.Create(c, genre)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})

}

// Update   	 godoc
// @Summary      Update genre
// @Tags         genres
// @Accept       json
// @Produce      json
// @Param        id      path      int                true  "Genre ID"
// @Param        request body      updateGenreRequest true  "Genre payload"
// @Success      200  "OK"
// @Failure      400  {object}  models.ApiError "Invalid payload"
// @Failure      404  {object}  models.ApiError "Genre not found"
// @Failure      500  {object}  models.ApiError
// @Router       /genres/{id} [put]
// @Security Bearer
func (h *GenreHandler) Update(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid Genre Id"))
		return
	}

	_, err = h.genresRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	var request updateGenreRequest
	err = c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Could not update genre"))
		return
	}

	genre := models.Genre{
		Title: request.Title,
	}

	h.genresRepo.Update(c, id, genre)

	c.Status(http.StatusOK)
}

// Delete   	 godoc
// @Summary      Delete genre
// @Tags         genres
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Genre ID"
// @Success      200  "OK"
// @Failure      400  {object}  models.ApiError "Invalid Genre Id"
// @Failure      404  {object}  models.ApiError "Genre not found"
// @Failure      500  {object}  models.ApiError
// @Router       /genres/{id} [delete]
// @Security Bearer
func (h *GenreHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid Genre Id"))
		return
	}

	_, err = h.genresRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	err = h.genresRepo.Delete(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}
