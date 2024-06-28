package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type UserInfo struct {
	Sub      string `json:"sub"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Username string `json:"preferred_username"`
	Profile  string `json:"profile"`
	Picture  string `json:"picture"`
}

func ValidateToken(token string) (string, error) {
	// Make a request to the Roblox OAuth userinfo endpoint
	req, err := http.NewRequest("GET", "https://apis.roblox.com/oauth/v1/userinfo", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		var errorResponse struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			return "", err
		}
		return "", errors.New(errorResponse.ErrorDescription)
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to validate token")
	}

	var userInfo UserInfo
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		return "", err
	}

	// Extract the user ID from the profile URL
	profileURL := userInfo.Profile
	parts := strings.Split(profileURL, "/")
	userIDStr := parts[len(parts)-2]
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(userID, 10), nil
}
