package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Token struct {
	Access_Token string `json:"access_token"`
	Type         string `json:"token_type"`
	Duration     int    `json:"expires_in"`
}

var access_token string = GetToken()

func main() {

	http.ListenAndServe("localhost:8080", nil)
}

func GetToken() string {
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

	return decodeData.Access_Token
}
