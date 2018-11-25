/*
* @Author: Bartuccio Antoine
* @Date:   2018-09-18 22:51:17
* @Last Modified by:   Bartuccio Antoine
* @Last Modified time: 2018-11-25 00:19:33
 */

package goproviders

import (
	"strconv"
	"time"

	tmdb "github.com/ryanbradynd05/go-tmdb"

	"gitlab.com/GT-RIMi/RIMo/media"
	"gitlab.com/GT-RIMi/RIMo/query"
	"gitlab.com/GT-RIMi/RIMo/settings"
)

// TmdbProvider query.Provider interface for TMDB
type TmdbProvider struct {
	conf      tmdb.Config
	api       *tmdb.TMDb
	language  string
	imgOrigin string
}

func (db *TmdbProvider) movieShortToMedia(movie tmdb.MovieShort) media.Media {
	release, err := time.Parse("2006-01-02", movie.ReleaseDate)
	if err != nil {
		release = time.Now()
	}
	return media.Media{
		ID:          strconv.Itoa(movie.ID),
		Title:       movie.Title,
		Image:       db.imgOrigin + movie.PosterPath,
		ReleaseDate: media.JSONTime{release},
		Adult:       movie.Adult,
	}
}

func (db *TmdbProvider) tvShortToMedia(serie tmdb.TvShort) media.Media {
	release, err := time.Parse("2006-01-02", serie.FirstAirDate)
	if err != nil {
		release = time.Now()
	}
	return media.Media{
		ID:          strconv.Itoa(serie.ID),
		Title:       serie.Name,
		Image:       db.imgOrigin + serie.PosterPath,
		ReleaseDate: media.JSONTime{release},
		Adult:       false,
	}
}

// Init initialize the backend
func (db *TmdbProvider) Init() {
	db.conf = tmdb.Config{
		ApiKey:   settings.SettingsValue("tmdb_api_v3_key").(string),
		UseProxy: false,
		Proxies:  []tmdb.Proxy{},
	}
	db.api = tmdb.Init(db.conf)
	db.language = "fr"
	db.imgOrigin = "https://image.tmdb.org/t/p/w300_and_h450_bestv2"
}

// IsCompatible check compatibility of the backend with a given media
func (db *TmdbProvider) IsCompatible(mediaType media.MediaType) bool {
	return mediaType == media.Movie || mediaType == media.Serie
}

// Search makes a search on TMD backend with a QueryStruct
func (db *TmdbProvider) Search(q query.QueryStruct) query.QueryResponse {
	var resp query.QueryResponse
	var supp []media.Media
	options := db.options(q)

	if q.MediaType == media.Movie {
		resp = db.searchMovie(q, options)
		supp = db.similarMovie(q, options)
	} else {
		resp = db.searchTv(q, options)
		supp = db.similarSerie(q, options)
	}

	for _, media := range supp {
		resp.Medias = append(resp.Medias, media)
	}

	resp.RemoveDuplicates()
	return resp
}

func (db *TmdbProvider) discoverMovie(id int) query.QueryResponse {

	resp, err := db.api.GetMovieSimilar(int(id), nil)
	if err != nil {
		return query.QueryResponse{
			Status: query.NotFound,
		}
	}

	medias := []media.Media{}
	for _, r := range resp.Results {
		medias = append(medias, db.movieShortToMedia(r))
	}
	return query.QueryResponse{
		Medias: medias,
		Status: query.Found,
	}
}

func (db *TmdbProvider) discoverTv(id int) query.QueryResponse {

	resp, err := db.api.GetTvSimilar(id, nil)
	if err != nil {
		return query.QueryResponse{
			Status: query.NotFound,
		}
	}

	medias := []media.Media{}
	for _, r := range resp.Results {
		medias = append(medias, db.tvShortToMedia(r))
	}
	return query.QueryResponse{
		Medias: medias,
		Status: query.Found,
	}
}

// Discover makes a search on TMD backend to find related content
func (db *TmdbProvider) Discover(q query.QueryDiscoverStruct) query.QueryResponse {

	id, err := strconv.ParseInt(q.Media.ID, 10, 0)
	if err != nil {
		return query.QueryResponse{
			Status: query.Err,
		}
	}
	if q.MediaType == media.Movie {
		return db.discoverMovie(int(id))
	}
	return db.discoverTv(int(id))
}

func (db *TmdbProvider) similarMovie(q query.QueryStruct, options map[string]string) []media.Media {
	response := []media.Media{}

	_, okYear := options["year"]
	_, okFirstAirDate := options["first_air_date_year"]

	if !okYear && !okFirstAirDate {
		return response
	}

	results, err := db.api.DiscoverMovie(options)

	if err != nil || len(results.Results) == 0 {
		return response
	}

	for _, movie := range results.Results {
		m := db.movieShortToMedia(movie)
		m.Suggested = true
		response = append(response, m)
	}

	return response
}

func (db *TmdbProvider) similarSerie(q query.QueryStruct, options map[string]string) []media.Media {
	response := []media.Media{}

	_, okYear := options["year"]
	_, okFirstAirDate := options["first_air_date_year"]

	if !okYear && !okFirstAirDate {
		return response
	}

	results, err := db.api.DiscoverTV(options)

	if err != nil || len(results.Results) == 0 {
		return response
	}

	for _, serie := range results.Results {
		m := db.tvShortToMedia(serie)
		m.Suggested = true
		response = append(response, m)
	}

	return response
}

func (db *TmdbProvider) options(q query.QueryStruct) map[string]string {
	options := make(map[string]string)

	options["include_adult"] = strconv.FormatBool(q.Adult)

	if q.Year != 0 {
		year := strconv.Itoa(q.Year)
		options["first_air_date_year"] = year
		options["year"] = year
	}

	if q.Page > 0 {
		options["page"] = strconv.FormatUint(uint64(q.Page), 10)
	}

	return options
}

func (db *TmdbProvider) searchMovie(q query.QueryStruct, options map[string]string) query.QueryResponse {
	response := query.QueryResponse{}
	movies, err := db.api.SearchMovie(q.Title, options)
	if err != nil {
		response.Status = query.Err
		return response
	}

	if len(movies.Results) == 0 {
		response.Status = query.NotFound
		return response
	}

	response.Status = query.Found
	medias := []media.Media{}

	for _, movie := range movies.Results {
		medias = append(medias, db.movieShortToMedia(movie))
	}

	response.Medias = medias

	return response
}

func (db *TmdbProvider) searchTv(q query.QueryStruct, options map[string]string) query.QueryResponse {
	response := query.QueryResponse{}
	series, err := db.api.SearchTv(q.Title, options)
	if err != nil {
		response.Status = query.Err
		return response
	}

	if len(series.Results) == 0 {
		response.Status = query.NotFound
		return response
	}

	response.Status = query.Found
	medias := []media.Media{}

	for _, serie := range series.Results {
		release, _ := time.Parse("2006-01-02", serie.FirstAirDate)
		if err != nil {
			release = time.Now()
		}
		medias = append(medias, media.Media{
			ID:          strconv.Itoa(serie.ID),
			Title:       serie.Name,
			Image:       db.imgOrigin + serie.PosterPath,
			ReleaseDate: media.JSONTime{release},
			Adult:       false,
		})
	}

	response.Medias = medias

	return response
}
