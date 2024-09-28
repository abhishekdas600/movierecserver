package router

import (
	"github.com/abhishekdas600/movierecserver/services/movies"

	"github.com/gin-gonic/gin"
)

func SetupMovieRouters(router *gin.Engine) {
    router.GET("/genres", movies.GetGenres)
	router.GET("/movie/filter", movies.GetMoviesByFilter)
	router.GET("/movies/:id", movies.GetMoviesById)
	router.GET("/trailer/:id", movies.GetTrailersByID)
	router.GET("/watchlist", movies.GetUserWatchlist)
	router.POST("/addwatchlist/:tmdb_id", movies.AddMovieToWatchlist)
	router.GET("/search", movies.SearchMovies)
	router.GET("/recommendations", movies.GetRecommendations)
	router.GET("/favourites", movies.GetUserFavourites)
	router.POST("/addtofavourites/:tmdb_id", movies.AddMovieToFavourites)
}