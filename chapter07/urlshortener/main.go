package main

import (
	"database/sql"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/dakaraj/go-api/chapter07/urlshortener/base62"
	"github.com/dakaraj/go-api/chapter07/urlshortener/models"
	"github.com/gin-gonic/gin"
)

var encStringPattern *regexp.Regexp

// DBClient stores database session informaion. Needs to be initialized once
type DBClient struct {
	db *sql.DB
}

// Model is the record struct
type Model struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

// getOriginalURL returns original URL to the user
func (driver *DBClient) getOriginalURL(c *gin.Context) {
	var url string
	encodedPath := c.Param("encoded-string")
	if !encStringPattern.MatchString(encodedPath) {
		c.JSON(400, gin.H{
			"error":        true,
			"errorMessage": "Shortened URL is invalid!",
		})
		return
	}

	id := base62.FromBase62(encodedPath)
	err := driver.db.QueryRow(`SELECT url FROM web_url WHERE id = $1`, id).Scan(&url)
	if err != nil {
		c.JSON(500, gin.H{
			"error":        true,
			"errorMessage": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"url": url,
	})

}

// generateShortURL saves it to DB and returns shortened string
func (driver *DBClient) generateShortURL(c *gin.Context) {
	var body Model
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"error":        true,
			"errorMessage": err.Error(),
		})
		return
	}

	stripFragment := strings.Split(body.URL, "#")
	body.URL = stripFragment[0]

	if _, err := url.ParseRequestURI(body.URL); err != nil {
		c.JSON(400, gin.H{
			"error":        true,
			"errorMessage": "Provided URL is invalid!",
		})
	}

	var lastRowID int
	err := driver.db.QueryRow(`INSERT INTO web_url (url) VALUES ($1) RETURNING id`, body.URL).Scan(&lastRowID)
	if err != nil {
		c.JSON(500, gin.H{
			"error":        true,
			"errorMessage": err.Error(),
		})
		return
	}

	encodedID := base62.ToBase62(int(lastRowID))

	c.JSON(200, gin.H{
		"encodedString": encodedID,
	})
}

func main() {
	db, err := models.InitDB()
	if err != nil {
		log.Fatal("Error encountered on DB initialization: ", err.Error())
	}
	dbClient := &DBClient{db: db}
	defer db.Close()

	encStringPattern, _ = regexp.Compile("^[a-zA-Z0-9]+$")

	router := gin.Default()
	router.GET("/v1/short/:encoded-string", dbClient.getOriginalURL)
	router.POST("/v1/short", dbClient.generateShortURL)

	router.Run(":8008")
}
