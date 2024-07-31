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
	form.Add("client_id", PMJAY_CLIENT_ID)
	form.Add("client_secret", PMJAY_CLIENT_SECRET)
	form.Add("username", PMJAY_CLIENT_USERNAME)
	form.Add("password", PMJAY_CLIENT_PASSWORD)

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