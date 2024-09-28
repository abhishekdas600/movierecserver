package movies

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/abhishekdas600/movierecserver/models"
	"os"
	"strconv"

	"net/http"
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type GenreListResponse struct {
	Genres []Genre `json:"genres"`
}
type MoviesFromId struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Overview    string `json:"overview"`
	PosterPath  string `json:"poster_path"`
	Genres      []Genre `json:"genres"`
	
}

type Trailer struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Key         string `json:"key"`
	Type        string `json:"type"`
	Official    bool   `json:"official"`
	PublishedAt string `json:"published_at"`
}

type Response struct {
	ID      int       `json:"id"`
	Results []Trailer `json:"results"`
}


 

const baseURL = "https://api.themoviedb.org/3"


func GetGenres(c *gin.Context) {
	apiKey := os.Getenv("TMDB_API_KEY")
	url := fmt.Sprintf("%s/genre/movie/list?api_key=%s", baseURL, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch genres"})
		return
	}
	defer resp.Body.Close()

	var genreList GenreListResponse
	if err := json.NewDecoder(resp.Body).Decode(&genreList); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode genres"})
		return
	}

	c.JSON(http.StatusOK, genreList)
}

func GetMoviesByFilter(c *gin.Context) {
	apiKey := os.Getenv("TMDB_API_KEY")
    filter := c.Query("filter") 
    genres := c.Query("genres") 
    page := c.DefaultQuery("page", "1")

    var url string

    baseURL := fmt.Sprintf("%s/discover/movie?api_key=%s&page=%s", baseURL, apiKey, page)

    switch filter {
    case "top_rated":
        url = fmt.Sprintf("%s&sort_by=vote_average.desc", baseURL)
    case "recent":
        url = fmt.Sprintf("%s&sort_by=release_date.desc", baseURL)
    case "popular":
        url = fmt.Sprintf("%s&sort_by=popularity.desc", baseURL)
    default:
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter type"})
        return
    }

    if genres != "" {
        url = fmt.Sprintf("%s&with_genres=%s", url, genres)
    }

    resp, err := http.Get(url)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
        return
    }
    defer resp.Body.Close()

    var movieList models.MovieListResponse
    if err := json.NewDecoder(resp.Body).Decode(&movieList); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode movie list"})
        return
    }

    c.JSON(http.StatusOK, movieList)
}

func GetMoviesById(c *gin.Context){
	apiKey := os.Getenv("TMDB_API_KEY")
   movieId:= c.Param("id")

   if movieId == ""{
	    c.JSON(http.StatusBadRequest, gin.H{"error": "No movie ID found"})
	         return
   }
   url := fmt.Sprintf("%s/movie/%s?api_key=%s&language=en-US", baseURL,movieId, apiKey)

   resp,err := http.Get(url)
        if err != nil {
	        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movie details"})
		return
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
	    c.JSON(resp.StatusCode, gin.H{"error": "Failed to fetch movie details from TMDB"})
	return
   }
   if resp.StatusCode != http.StatusOK {
	c.JSON(resp.StatusCode, gin.H{"error": "Failed to fetch movie details from TMDB"})
	return
   }

   var movieDetails MoviesFromId
	if err := json.NewDecoder(resp.Body).Decode(&movieDetails); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode movie details"})
		return
	}
	c.JSON(http.StatusOK, movieDetails)
	
}

func GetTrailersByID(c *gin.Context) {
    apiKey := os.Getenv("TMDB_API_KEY")

    movieId := c.Param("id")
    if movieId == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No movie ID found"})
        return
    }

    id, err := strconv.Atoi(movieId)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
        return
    }

    url := fmt.Sprintf("%s/movie/%d/videos?api_key=%s&language=en-US", baseURL, id, apiKey)

    resp, err := http.Get(url)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching trailers"})
        return
    }
    defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
        return
    }


    if resp.StatusCode != http.StatusOK {
        c.JSON(http.StatusNotFound, gin.H{"error": "Trailers not found"})
        return
    }

    var response Response
    if err := json.Unmarshal(body, &response); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding response"})
        return
    }

    var mainTrailer *Trailer
    for _, trailer := range response.Results {
        if trailer.Type == "Trailer" && trailer.Official {
            if mainTrailer == nil || trailer.PublishedAt > mainTrailer.PublishedAt {
                mainTrailer = &trailer
            }
        }
    }

    if mainTrailer == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "No official trailer found"})
        return
    }

    c.JSON(http.StatusOK, mainTrailer)
}

func SearchMovies(c *gin.Context) {
	apiKey := os.Getenv("TMDB_API_KEY")
	searchQuery := c.Query("query")
	page := c.DefaultQuery("page", "1") 

	if searchQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	url := fmt.Sprintf("%s/search/movie?api_key=%s&query=%s&page=%s", baseURL, apiKey, url.QueryEscape(searchQuery), page)
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch search results"})
		return
	}
	defer resp.Body.Close()


	var searchResults models.MovieListResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResults); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode search results"})
		return
	}

	var movieIDs []int
	for _, movie := range searchResults.Results {
		movieIDs = append(movieIDs, movie.ID)
	}

	session := sessions.Default(c)
	session.Set("searched_movie_ids", movieIDs) 
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"page":          searchResults.Page,
		"total_pages":   searchResults.TotalPages,
		"results":       searchResults.Results,
	})
}

func GetRecommendations(c *gin.Context) {
	session := sessions.Default(c)
	movieIDs := session.Get("searched_movie_ids")
	if movieIDs == nil {
		c.JSON(http.StatusOK, gin.H{"recommendations": nil}) 
		return
	}

	ids, ok := movieIDs.([]int)
	if !ok || len(ids) == 0 {
		c.JSON(http.StatusOK, gin.H{"recommendations": nil}) 
		return
	}

	apiKey := os.Getenv("TMDB_API_KEY")
	recommendations := []models.Movie{} 

	for _, id := range ids {
		url := fmt.Sprintf("%s/movie/%d/recommendations?api_key=%s", baseURL, id, apiKey)
		resp, err := http.Get(url)
		if err != nil {
			continue 
		}
		defer resp.Body.Close()

		var recResponse models.MovieListResponse
		if err := json.NewDecoder(resp.Body).Decode(&recResponse); err != nil {
			continue 
		}

		recommendations = append(recommendations, recResponse.Results...)
		if len(recommendations) >= 20 {
			break 
		}
	}

	if len(recommendations) > 20 {
		recommendations = recommendations[:20]
	}

	c.JSON(http.StatusOK, gin.H{"recommendations": recommendations})
}