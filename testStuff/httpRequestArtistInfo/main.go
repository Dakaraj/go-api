package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Album contains basic information about album
type Album struct {
	Name      string
	Playcount float64
	MBID      string
}

func main() {
	response, err := http.Get("http://ws.audioscrobbler.com/2.0/?method=artist.gettopalbums&artist=dark%20tranquillity&api_key=2d454ded2fc3780b003469a3d823cdac&limit=8984&format=json")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	start := time.Now()

	var jsonBody map[string]map[string]interface{}
	json.Unmarshal(body, &jsonBody)
	albumList := jsonBody["topalbums"]["album"].([]interface{})

	var validAlbumsList []Album

	for _, item := range albumList {
		itemMap := item.(map[string]interface{})
		if mbid, ok := itemMap["mbid"]; ok {
			album := Album{
				MBID:      mbid.(string),
				Playcount: itemMap["playcount"].(float64),
				Name:      itemMap["name"].(string),
			}
			validAlbumsList = append(validAlbumsList, album)
			// fmt.Printf("Name: %s, Playcounts: %.0f, MBID: %s\n", album.Name, album.Playcount, album.MBID)
		}
	}

	total := time.Now().UnixNano() - start.UnixNano()

	fmt.Printf("Time elapsed: %.3fs\n", float64(total)/1000000000)
	fmt.Printf("Total albums: %d\n", len(validAlbumsList))
}
