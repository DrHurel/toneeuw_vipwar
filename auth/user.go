package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"utils"
)

type RETURN_VALUE int

const (
	SUCCESS              RETURN_VALUE = 0
	FAILED_TO_GET_TOKEN  RETURN_VALUE = 1
	FAILED_TO_GET_REWARD RETURN_VALUE = 2
	FAILED_TO_UPDATE     RETURN_VALUE = 3
	FAILED               RETURN_VALUE = 4
)

type User struct {
	ClientId     string `json:"client_id"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    string `json:"expires_in"`
	Secret       string `json:"secret"`
}

func (this *User) Credentials() RETURN_VALUE {
	if this.ExpiresIn == "" {
		shouldReturn, returnValue := this.FetchToken()
		if shouldReturn {
			return returnValue
		}

	}

	if this.isValidate() {
		println("\n=== Token is valid ===\n")

		return SUCCESS
	} else {
		if this.Refresh() == FAILED {
			this.FetchToken()
		}
	}

	return SUCCESS
}

func (this *User) Refresh() RETURN_VALUE {
	data := url.Values{}
	data.Set("client_id", this.ClientId)   // replace with your actual client ID
	data.Set("client_secret", this.Secret) // replace with your actual client secret
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", this.RefreshToken) // assuming this.RefreshToken is your refresh token

	req, err := http.NewRequest("POST", "https://id.twitch.tv/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return FAILED // replace with your actual error value
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return FAILED // replace with your actual error value
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error refreshing token:", resp.StatusCode)
		return FAILED // replace with your actual error value
	}

	// Parse the response body to get the new access token
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	this.Token = result["access_token"].(string)

	for key, value := range result {
		if key == "access_token" {
			this.Token = value.(string)
		}
		if key == "refresh_token" {
			this.RefreshToken = value.(string)
		}
	}

	return SUCCESS // replace with your actual success value
}

func (this *User) isValidate() bool {
	println("Checking if token is valid")
	baseURL, err := url.Parse("https://id.twitch.tv/oauth2/validate")
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return false
	}

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return false
	}

	req.Header.Set(utils.AuthHeader, "OAuth "+this.Token) // remove the extra space after "OAuth"

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return false
	}

	println("Response code:", resp.StatusCode)
	body := make([]byte, 1000)
	resp.Body.Read(body)

	return resp.StatusCode != 401
}

func (this *User) FetchToken() (bool, RETURN_VALUE) {
	cmd := exec.Command("twitch", "token", "-u", "-s", "channel:manage:redemptions channel:manage:vips")

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(errb.String(), "\n")

	var token string
	for i, line := range lines {
		println(line)
		temp := strings.Split(line, ": ")
		if i == 2 {
			token = temp[len(temp)-1]
		}
		if i == 3 {
			this.RefreshToken = temp[len(temp)-1]
		}
		if i == 4 {
			this.ExpiresIn = temp[len(temp)-1]
		}
	}

	if token == "" {
		return true, FAILED_TO_GET_TOKEN
	} else {
		fmt.Println("Token:", token)
	}
	this.Token = token
	return false, SUCCESS
}
