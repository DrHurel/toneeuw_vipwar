package auth

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"time"
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
}

func (this *User) Credentials() RETURN_VALUE {
	timeVal := strings.Split(this.ExpiresIn, " ")
	var tDate []int
	temp := strings.Split(timeVal[0], "-")
	for _, t := range temp {
		tInt, _ := strconv.Atoi(t)
		tDate = append(tDate, tInt)
	}
	var tTime []int

	for _, t := range strings.Split(timeVal[1], ":") {
		tInt, _ := strconv.Atoi(t)
		tTime = append(tTime, tInt)
	}

	t := time.Date(tDate[0], time.Month(tDate[1]), tDate[2], tTime[0], tTime[1], tTime[2], 0, time.Local)

	if t.Before(time.Now()) {

		shouldReturn, returnValue := this.FetchToken()
		if shouldReturn {
			return returnValue
		}
	} else {
		// Create the HTTP request
		// refresh token
		if this.isValidate() {
			return SUCCESS
		} else {
			return this.Refresh()
		}
	}

	return SUCCESS
}

func (this *User) Refresh() RETURN_VALUE {
	return SUCCESS
}

func (this *User) isValidate() bool {
	baseURL, err := url.Parse("https://api.twitch.tv/helix/channel_points/custom_rewards")
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return false
	}

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return false
	}

	req.Header.Set("Authorization", "OAuth "+this.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return false
	}

	return resp.StatusCode != 401
}

func (this *User) FetchToken() (bool, RETURN_VALUE) {
	cmd := exec.Command("twitch", "token", "-u", "-s", "channel:read:redemptions")

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
