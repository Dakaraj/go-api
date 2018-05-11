package main

import (
	"database/sql"
	"log"
	"net/http"
	"net/url"

	"github.com/dakaraj/go-api/shortenURL/dbutils"
	"github.com/emicklei/go-restful"
	"github.com/mattheath/base62"
	_ "github.com/mattn/go-sqlite3"
)

// DB is a driver reference for main database
var DB *sql.DB

const addend = 12345

// ShortenURLResource is a struct for working with table shorten_url
type ShortenURLResource struct {
	OriginalURL string
	ShortToken  string
}

// Register is used for binding paths and methods to ShortenURLResource
func (su *ShortenURLResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/v1/shorten")
	ws.Route(ws.GET("/{token}").To(su.redirrectToOriginalURL).Produces(restful.MIME_JSON))
	ws.Route(ws.POST("").To(su.shortenURL).Produces(restful.MIME_JSON))
	container.Add(ws)
}

// GET /v1/shorten/{token}
func (su *ShortenURLResource) redirrectToOriginalURL(request *restful.Request, response *restful.Response) {
	token := request.PathParameter("token")
	log.Println("Token:", token)
	decodedToken := base62.DecodeToInt64(token) - addend
	log.Println("Decoded token:", decodedToken)
	err := DB.QueryRow(`
SELECT original_url, shorten_token
FROM shorten_url
WHERE id = ?;
`, decodedToken).Scan(&su.OriginalURL, &su.ShortToken)
	if err != nil {
		log.Println(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, "Token is invalid, no original URL found!")
	} else {
		response.WriteEntity(su)
	}
}

// POST /v1/shorten?url={url}
func (su *ShortenURLResource) shortenURL(request *restful.Request, response *restful.Response) {
	originalURL := request.QueryParameter("url")
	log.Println("Original url:", originalURL)
	_, err := url.ParseRequestURI(originalURL)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, "Provided URL is not valid")
	} else {
		statement, err := DB.Prepare(`
INSERT INTO shorten_url
(original_url)
VALUES (?);
`)
		result, err := statement.Exec(originalURL)
		if err != nil {
			response.AddHeader("Content-Type", "text/plain")
			response.WriteErrorString(http.StatusBadRequest, err.Error())
		} else {
			lastRowID, _ := result.LastInsertId()
			log.Println("Last row ID:", lastRowID)
			token := base62.EncodeInt64(lastRowID + addend)
			log.Println(token)
			DB.Exec("UPDATE shorten_url SET shorten_token = ?", token)
			su.OriginalURL = originalURL
			su.ShortToken = token
			response.WriteHeaderAndEntity(http.StatusCreated, su)
		}
	}
}

func main() {
	var err error
	DB, err = sql.Open("sqlite3", "./shortenURL/shortenURL.db")
	if err != nil {
		log.Fatal("DB Initialization failed:", err.Error())
	}
	dbutils.Initialize(DB)

	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})

	su := ShortenURLResource{}
	su.Register(wsContainer)

	server := &http.Server{
		Addr:    ":80",
		Handler: wsContainer,
	}
	log.Fatal(server.ListenAndServe())
}
