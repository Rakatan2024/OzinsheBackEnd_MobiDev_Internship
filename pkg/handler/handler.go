package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"ozinshe/pkg/entity"
	"ozinshe/pkg/service"
)

type Handler struct {
	log *log.Logger
	svc service.SvcInterface
}

func CreateHandler(service service.SvcInterface, log *log.Logger) Handler {
	return Handler{svc: service, log: log}
}

func (h *Handler) InitRoutes() *gin.Engine {
	ginServer := gin.Default()
	ginServer.MaxMultipartMemory = 8 << 20
	ginServer.Static("/assets", "./assets")
	core := ginServer.Group("/core", h.AuthMiddleware())
	{
		core.POST("/movie", h.AdminRoleMiddleware(), h.CreateMovie)
		core.POST("/favorites", h.CreateFavoriteMovie)
		core.POST("/movie/:id/season/:seasonId", h.AdminRoleMiddleware(), h.AddNewSeries)
		core.POST("/movie/:id/season", h.AdminRoleMiddleware(), h.AddNewSeason)
		core.GET("/home", h.HomePageHandler)
		core.GET("/movies/page", h.GetAllMovies)
		core.GET("/movie/genres", h.GetAllGenres)
		core.GET("/movie/:id", h.GetMovieById)
		core.GET("/movie/:id/season/:seasonId", h.GetMovieSeasonById)
		core.GET("/movie/:id/season/:seasonId/series/:seriesId")
		core.GET("/categories", h.GetCategories)
		core.GET("/mainPage/", h.GetMovieMainsByAllCategory)
		core.GET("/user/profile", h.GetUserProfile)
		core.GET("/movieMain/category", h.GetMovieMainsByCategory)
		core.GET("/movieMain/search", h.GetMovieMainsByTitle)
		core.GET("/movieMain/search/genre", h.GetMovieMainsByGenre)
		core.GET("/favorites/", h.GetFavoriteMovies)
		core.PUT("/user/profile", h.UpdateUserProfile)
		core.PUT("/user/profile/password", h.ChangePassword)
		core.PUT("/movie/:id", h.AdminRoleMiddleware(), h.UpdateMovieById)
		core.PUT("/movie/:id/season/:seasonId", h.AdminRoleMiddleware(), h.UpdateSeason)
		core.PUT("/movie/:id/season/:seasonId/series/:seriesId", h.AdminRoleMiddleware(), h.UpdateSeries)
		core.DELETE("/movie/:id/season/:seasonId", h.AdminRoleMiddleware(), h.DeleteMovieSeason)
		core.DELETE("/movie/:id/season/:seasonId/series/:seriesId", h.AdminRoleMiddleware(), h.DeleteMovieSeries)
		core.DELETE("/movie/:id", h.AdminRoleMiddleware(), h.DeleteMovieById)
		core.DELETE("/favorites/", h.DeleteFavoriteMovies)
	}
	auth := ginServer.Group("/auth")
	{
		auth.POST("/sign-up", h.SignUp)
		auth.GET("/verifyAccount", h.VerifyAccount)
		auth.POST("/sign-in", h.SignIn)
		auth.POST("/passwordRecover", h.PasswordRecover)
	}
	ginServer.GET("/swagger", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return ginServer
}

func (h *Handler) WriteHTTPResponse(c *gin.Context, statusCode int, msg string) {
	c.AbortWithStatusJSON(statusCode, entity.ErrorJSONResponse{Message: msg})
}
