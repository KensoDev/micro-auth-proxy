package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type GithubAuthRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
}

type GithubAuthResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
}

func main() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		client := &http.Client{}

		githubRequest := &GithubAuthRequest{
			ClientID:     os.Getenv("CLIENT_ID"),
			ClientSecret: os.Getenv("CLIENT_SECRET"),
			Code:         code,
		}

		jsonResponse, _ := json.Marshal(githubRequest)

		req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(jsonResponse))

		if err != nil {
			fmt.Println(err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		resp, err := client.Do(req)

		responseBody, _ := ioutil.ReadAll(resp.Body)

		githubResponse := &GithubAuthResponse{}

		err = json.Unmarshal(responseBody, githubResponse)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Requesting user")

		url := fmt.Sprintf("https://api.github.com/user?access_token=%s", githubResponse.AccessToken)
		fmt.Println(url)
		userResponse, err := http.Get(url)

		if err != nil {
			log.Fatal(err)
		}

		responseBody, _ = ioutil.ReadAll(userResponse.Body)
		fmt.Println(string(responseBody))
	})

	log.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}
