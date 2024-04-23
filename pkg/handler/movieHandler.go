package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"ozinshe/pkg/entity"
	"strconv"
	"strings"
)

func (h *Handler) CreateMovie(c *gin.Context) {
	err := c.Request.ParseMultipartForm(10 << 20) // Parse up to 10 MB of data
	if err != nil {
		h.log.Printf("error during ParseMultipartForm in CreateMovie(handler): %v", err)
		//http.Error(w, "Error parsing multipart form", http.StatusBadRequest)
		h.WriteHTTPResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	formData := c.Request.MultipartForm

	//jsonString := formData.Value["json"][0]
	//ScreenshotFileHeaders := formData.File["screenshots"]
	//PosterFileHeaders := formData.File["poster"]

	movie := entity.Movie{}
	if err := json.NewDecoder(strings.NewReader(formData.Value["json"][0])).Decode(&movie); err != nil {
		h.WriteHTTPResponse(c, http.StatusBadRequest, "Invalid input body")
		return
	}
	if err := h.svc.CreateMovie(&movie, formData); err != nil {
		h.WriteHTTPResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	h.WriteHTTPResponse(c, http.StatusCreated, "movie created!")
}

func (h *Handler) GetAllMovies(c *gin.Context) {
	params := c.Request.URL.Query()
	limit, err := strconv.Atoi(params.Get("limit"))
	if err != nil {
		h.log.Print("limit is not a number")
		h.WriteHTTPResponse(c, http.StatusBadRequest, "limit is not a number")
		return
	}
	offset, err := strconv.Atoi(params.Get("offset"))
	if err != nil {
		h.log.Print("offset is not a number")
		h.WriteHTTPResponse(c, http.StatusBadRequest, "offset is not a number")
		return
	}
	movies, err := h.svc.GetAllMovies(limit, offset)
	if err != nil {
		h.WriteHTTPResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, movies)
}

func (h *Handler) GetMovieById(c *gin.Context) {
	params := c.Param("id")
	movieID, err := strconv.Atoi(params)
	if err != nil {
		h.log.Print("movieId is not a number")
		h.WriteHTTPResponse(c, http.StatusBadRequest, "movieId is not a number")
		return
	}
	movie, err := h.svc.GetMovieById(movieID)
	if err != nil {
		h.WriteHTTPResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, movie)
}

func (h *Handler) UpdateMovieById(c *gin.Context) {
	params := c.Param("id")
	movieID, err := strconv.Atoi(params)
	if err != nil {
		h.log.Print("movieId is not a number")
		h.WriteHTTPResponse(c, http.StatusBadRequest, "movieId is not a number")
		return
	}
	movie := entity.Movie{}
	if err := c.BindJSON(&movie); err != nil {
		h.WriteHTTPResponse(c, http.StatusBadRequest, "Invalid input body")
		return
	}
	movie.Id = movieID
	err = h.svc.UpdateMovieById(&movie)
	if err != nil {
		h.WriteHTTPResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	h.WriteHTTPResponse(c, http.StatusOK, "movie updated!")
}

func (h *Handler) DeleteMovieById(c *gin.Context) {
	params := c.Param("id")
	movieID, err := strconv.Atoi(params)
	if err != nil {
		h.log.Print("movieId is not a number")
		h.WriteHTTPResponse(c, http.StatusBadRequest, "movieId is not a number")
		return
	}
	err = h.svc.DeleteMovieById(movieID)
	if err != nil {
		h.WriteHTTPResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	h.WriteHTTPResponse(c, http.StatusOK, "movie deleted!")
}

func (h *Handler) GetMovieSeasonById(c *gin.Context) {
	seasonId := c.Param("seasonId")
	movieId := c.Param("id")
	intSeasonId, err := strconv.Atoi(seasonId)
	if err != nil {
		h.log.Print("seasonId is not a number")
		h.WriteHTTPResponse(c, http.StatusBadRequest, "seasonId is not a number")
		return
	}
	intMovieId, err := strconv.Atoi(movieId)
	if err != nil {
		h.log.Print("seasonId is not a number")
		h.WriteHTTPResponse(c, http.StatusBadRequest, "movieId is not a number")
		return
	}
	seasonSeriesLinks, err := h.svc.GetMovieSeason(intMovieId, intSeasonId)
	if err != nil {
		h.WriteHTTPResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, seasonSeriesLinks)
}

func (h *Handler) GetMovieSeriesById(c *gin.Context) {
	seasonId := c.Param("seasonId")
	movieId := c.Param("movieId")
	episodeId := c.Param("seriesId")
	intSeasonId, err := strconv.Atoi(seasonId)
	if err != nil {
		h.log.Print("seasonId is not a number")
		h.WriteHTTPResponse(c, http.StatusBadRequest, "seasonId is not a number")
		return
	}
	intMovieId, err := strconv.Atoi(movieId)
	if err != nil {
		h.log.Print("seasonId is not a number")
		h.WriteHTTPResponse(c, http.StatusBadRequest, "movieId is not a number")
		return
	}
	intEpisodeId, err := strconv.Atoi(episodeId)
	if err != nil {
		h.log.Print("seasonId is not a number")
		h.WriteHTTPResponse(c, http.StatusBadRequest, "episodeId is not a number")
		return
	}
	seasonSeriesLinks, err := h.svc.GetMovieSeries(intMovieId, intEpisodeId, intSeasonId)
	if err != nil {
		h.WriteHTTPResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, seasonSeriesLinks)
}

func (h *Handler) GetMovieMainsByTitle(c *gin.Context) {
	params := c.Request.URL.Query()
	title := params.Get("title")
	movieMains, err := h.svc.GetMovieMainsByTitle(title)
	if err != nil {
		h.log.Print("error in GetMovieMainsByTitle(handler)")
		h.WriteHTTPResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, movieMains)
}
