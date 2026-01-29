package handlers

import (
	"filmservice/models"
	"filmservice/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WatchListHandlers struct {
	watchListRepo *repositories.WatchListRepository
}

func NewWatchListHandlers(watchListRepo *repositories.WatchListRepository) *WatchListHandlers {
	return &WatchListHandlers{
		watchListRepo: watchListRepo,
	}
}

// GetAll   	 godoc
// @Summary      Get all movies from watchlist
// @Tags         watchlist
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Movie "OK"
// @Failure      500  {object}  models.ApiError
// @Router       /watchlist [get]
// @Security Bearer
func (h *WatchListHandlers) GetAll(c *gin.Context) {
	movies, err := h.watchListRepo.GetAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, movies)
}

func (h *WatchListHandlers) parseMovieId(c *gin.Context) (int, bool) {
	idStr := c.Param("movieId")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid Movie Id"))
		return 0, false
	}

	return id, true
}

// Toggle   	 godoc
// @Summary      Add or remove movie from watchlist (toggle)
// @Tags         watchlist
// @Accept       json
// @Produce      json
// @Param        movieId   path      int  true  "Movie ID"
// @Success      204  "No Content"
// @Failure      400  {object}  models.ApiError "Invalid Movie Id"
// @Failure      500  {object}  models.ApiError
// @Router       /watchlist/{movieId} [post]
// @Security Bearer
func (h *WatchListHandlers) Toggle(c *gin.Context) {
	id, ok := h.parseMovieId(c)
	if !ok {
		return
	}

	exists, err := h.watchListRepo.Exists(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	if exists {
		err = h.watchListRepo.Delete(c, id)
	} else {
		err = h.watchListRepo.Add(c, id)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusNoContent)
}

// Delete   	 godoc
// @Summary      Remove movie from watchlist
// @Tags         watchlist
// @Accept       json
// @Produce      json
// @Param        movieId   path      int  true  "Movie ID"
// @Success      204  "No Content"
// @Failure      400  {object}  models.ApiError "Invalid Movie Id"
// @Failure      500  {object}  models.ApiError
// @Router       /watchlist/{movieId} [delete]
// @Security Bearer
func (h *WatchListHandlers) Delete(c *gin.Context) {
	id, ok := h.parseMovieId(c)
	if !ok {
		return
	}

	err := h.watchListRepo.Delete(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusNoContent)
}
