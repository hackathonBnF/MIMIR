/*
* @Author: Bartuccio Antoine
* @Date:   2018-09-15 22:17:31
* @Last Modified by:   Bartuccio Antoine
* @Last Modified time: 2018-11-24 23:23:18
 */

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"

	"gitlab.com/GT-RIMi/RIMo/goproviders"
	"gitlab.com/GT-RIMi/RIMo/query"
	"gitlab.com/GT-RIMi/RIMo/settings"
)

func main() {

	// Loads settings
	settings.InitSettings("settings.json", "settings_custom.json")

	// Register Providers
	if err := query.RegisterProvider(&goproviders.TmdbProvider{}, "tmdb"); err != nil {
		log.Printf("Could not connect backend tmdb %s\n", err)
	}

	// Configure Gin
	router := gin.Default()
	router.Delims("{[{", "}]}")
	router.LoadHTMLGlob("templates/*")
	router.StaticFS("/static", http.Dir("./static"))

	// Configure sessions
	store := memstore.NewStore([]byte(settings.SettingsValue("session_secret").(string)))
	router.Use(sessions.Sessions("memory", store))

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.gohtml", gin.H{
			"settings": settings.SettingsValue("debug").(bool),
		})
	})

	router.POST("/discover", func(c *gin.Context) {
		var q query.QueryDiscoverStruct

		if err := c.ShouldBindJSON(&q); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, query.Discover(&q))

	})

	router.POST("/search", func(c *gin.Context) {
		var search query.QueryStruct
		var lastSearch *query.QueryStore

		if err := c.ShouldBindJSON(&search); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		session := sessions.Default(c)

		sessionKey := fmt.Sprintf("last_search_%d", search.MediaType)
		lastSearchRow := session.Get(sessionKey)
		if lastSearchRow != nil {
			if err := json.Unmarshal(lastSearchRow.([]byte), &lastSearch); err != nil {
				log.Printf("could not retrieve lastSearch %s\n", err)
				lastSearch = nil
			}
		}

		response := query.QueryStore{
			Query:    search,
			Response: query.Search(&search, lastSearch),
		}
		c.JSON(http.StatusOK, response.Response)

		// Save response in session
		exportedSearch, err := json.Marshal(response)
		if err != nil {
			log.Printf("error marshaling response for session %s\n", err)
			return
		}

		session.Set(sessionKey, exportedSearch)
		if err := session.Save(); err != nil {
			log.Printf("error saving session %s\n", err)
		}

	})
	if err := router.Run(); err != nil {
		log.Fatal(err)
	}
}
