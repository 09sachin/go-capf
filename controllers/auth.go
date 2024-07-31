package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}



func generateToken() (string, error) {
	form := url.Values{}
	form.Add("grant_type", "password")
	form.Add("client_id", "775863c3-38aa-451e-90da-0d145e2ce4fe")
	form.Add("client_secret", "cXhVcfpUJpKOri_EG4XqKlmOj6ZvmPMa8AIEnMsPU7gkw3ET5purRwKOOfF0qIHbB7HXRhrQ2jnU8wA07542Dg")
	form.Add("username", "oauthclient")
	form.Add("password", "Password@1234")

	req, err := http.NewRequest("POST", "https://apis.pmjay.gov.in/idmtoken", bytes.NewBufferString(form.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get token, status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResponse TokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %s", tokenResponse.TokenType, tokenResponse.AccessToken), nil
}