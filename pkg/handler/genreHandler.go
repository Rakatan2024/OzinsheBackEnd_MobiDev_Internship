package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ozinshe/pkg/entity"
	"strconv"
)

func (h *Handler) GetAllGenres(c *gin.Context) {
	genres, err := h.svc.GetAllGenres()
	if err != nil {
		h.WriteHTTPResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, genres)
}

func (h *Handler) GetMovieMainsByGenre(c *gin.Context) {
	params := c.Request.URL.Query()
	userId := c.Value("decodedClaims").(*entity.Claims).Sub
	genreId, err := strconv.Atoi(params.Get("genreId"))
	if err != nil {
		h.log.Print("genreId is not a number")
		h.WriteHTTPResponse(c, http.StatusBadRequest, "genreId is not a number")
		return
	}
	movieMains, err := h.svc.GetMovieMainsByGenre(userId, genreId)
	if err != nil {
		h.log.Print("error in GetMovieMainsByCategory(handler)")
		h.WriteHTTPResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, movieMains)
}
