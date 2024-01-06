package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"
)

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Duration    int    `json:"expires_in"`
}

type Albums struct {
	Href  string  `json:"href"`
	Items []Album `json:"items"`
}

type Album struct {
	Group     string   `json:"album_group"`
	AlbumType string   `json:"album_type"`
	Artists   []Artist `json:"Artists"`
	Markets   []string `json:"available_markets"`
	Urls      []struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Href        string `json:"href"`
	Id          string `json:"id"`
	Images      []Image
	Name        string `json:"name"`
	ReleaseDate string `json:"release_date"`
	Precision   string `json:"release_date_precision"`
	TracksNb    int    `json:"total_tracks"`
	Type        string `json:"type"`
	Uri         string `json:"uri"`
}

type Image struct {
	Height int    `json:"height"`
	Link   string `json:"url"`
	Width  int    `json:"width"`
}

type Artist struct {
	URLs []struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Href string `json:"href"`
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Uri  string `json:"uri"`
}

type Track struct {
	Album struct {
		AlbumType string   `json:"album_type"`
		Artists   []Artist `json:"artists"`
		Markets   []string `json:"available_markets"`
		Urls      []struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href      string  `json:"href"`
		Id        string  `json:"id"`
		Images    []Image `json:"images"`
		Name      string  `json:"name"`
		Release   string  `json:"release_date"`
		Precision string  `json:"release_date_precision"`
		TracksNb  int     `json:"total_tracks"`
		Type      string  `json:"type"`
		Uri       string  `json:"uri"`
	} `json:"album"`
	Artists  []Artist `json:"artists"`
	Markets  []string `json:"available_markets"`
	DiscNb   int      `json:"disc_number"`
	Duration int      `json:"duration_ms"`
	Explicit bool     `json:"explicit"`
	Ids      []struct {
		Irc string `json:"isrc"`
	} `json:"external_ids"`
	Urls []struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Href       string `json:"href"`
	Id         string `json:"id"`
	IsLocal    bool   `json:"is_local"`
	Name       string `json:"name"`
	Popularity int    `json:"popularity"`
	Preview    string `json:"preview_url"`
	TrackNb    int    `json:"track_number"`
	Type       string `json:"type"`
	Uri        string `json:"uri"`
}

func main() {
	temp, err := template.ParseGlob("templates/*.html")
	if err != nil {
		fmt.Println("Erreur dans la récupération des templates : ", err)
		return
	}

	RootDoc, _ := os.Getwd()
	fileserver := http.FileServer(http.Dir(RootDoc + "/web/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fileserver))

	http.HandleFunc("/album/jul", func(w http.ResponseWriter, r *http.Request) {
		api_url := "https://api.spotify.com/v1/artists/3IW7ScrzXmPvZhB27hmfgy/albums"
		token := GetToken()
		httpClient := http.Client{
			Timeout: time.Second * 2,
		}

		req, errReq := http.NewRequest(http.MethodGet, api_url, nil)
		if errReq != nil {
			fmt.Println("Problème dans la requête d'obtentions des albums : ", errReq.Error())
		}

		req.Header.Add("Authorization", token.TokenType+" "+token.AccessToken)

		res, errRes := httpClient.Do(req)
		if res.Body != nil {
			defer res.Body.Close()
		} else {
			fmt.Println("Erreur dans l'envoi de la requête d'album : ", errRes.Error())
		}

		body, errBody := io.ReadAll(res.Body)
		if errBody != nil {
			fmt.Println("Erreur dans la lecture de la réponse d'ablums : ", errBody.Error())
		}

		var decodeData Albums
		json.Unmarshal(body, &decodeData)

		temp.ExecuteTemplate(w, "albums", decodeData)
	})

	http.HandleFunc("/track/sdm", func(w http.ResponseWriter, r *http.Request) {
		api_url := "https://api.spotify.com/v1/tracks/7A1nhuX64Y2JB206h3FjBK"
		token := GetToken()
		httpClient := http.Client{
			Timeout: time.Second * 2,
		}

		req, errReq := http.NewRequest(http.MethodGet, api_url, nil)
		if errReq != nil {
			fmt.Println("Erreur dans la requête de track : ", errReq.Error())
		}

		req.Header.Add("Authorization", token.TokenType+" "+token.AccessToken)

		res, resErr := httpClient.Do(req)
		if res.Body != nil {
			defer res.Body.Close()
		} else {
			fmt.Println("Problème dans l'envoi de la requête de track : ", resErr.Error())
		}

		body, errBody := io.ReadAll(res.Body)
		if errBody != nil {
			fmt.Println("Problème dans la lecture du crops de la requête de track : ", errBody.Error())
		}

		var decodeData Track
		json.Unmarshal(body, &decodeData)

		temp.ExecuteTemplate(w, "bolide alemand", decodeData)
	})

	http.ListenAndServe("localhost:8080", nil)
}

func GetToken() Token {
	api_url := "https://accounts.spotify.com/api/token"
	httpClient := http.Client{
		Timeout: time.Second * 2,
	}

	req, errReq := http.NewRequest(http.MethodGet, api_url, nil)
	if errReq != nil {
		fmt.Println("Erreur dans la requête d'obtention de token : ", errReq.Error())
	}

	res, errRes := httpClient.Do(req)
	if res.Body != nil {
		defer res.Body.Close()
	} else {
		fmt.Println("Problème dans l'envoi de la requête d'obtention de token : ", errRes.Error())
	}

	body, errBody := io.ReadAll(res.Body)
	if errBody != nil {
		fmt.Println("Problème dans la lecture du corps de la requête d'obtention de token : ", errBody.Error())
	}

	var decodeData Token
	json.Unmarshal(body, &decodeData)

	return decodeData
}
