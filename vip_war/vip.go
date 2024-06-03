package vipwar

import (
	"auth"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"utils"
)

type VIP struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

const errURLMsg = "Error parsing URL:"

func (this *VIP) Steal(mediator *auth.User, thief interface{}) {
	//get broadcaster id

	baseURL, err := url.Parse("https://api.twitch.tv/helix/users")
	if err != nil {
		fmt.Println(errURLMsg, err)
		return
	}

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		fmt.Println(errURLMsg, err)
		return
	}
	req.Header.Set(utils.AuthHeader, utils.FormatToken(mediator.Token))
	req.Header.Set(utils.ClientIdHeader, mediator.ClientId)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	// Parse the response body to get the new access token
	var result map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)

	broadcasterId := result["data"].([]interface{})[0].(map[string]interface{})["id"].(string)

	vidURL, err := url.Parse("https://api.twitch.tv/helix/channels/vips")
	if err != nil {
		log.Fatalln(errURLMsg, err)
		return
	}

	params := url.Values{}
	params.Add("broadcaster_id", broadcasterId)
	vidURL.RawQuery = params.Encode()

	req, err = http.NewRequest("GET", vidURL.String(), nil)
	if err != nil {
		log.Fatalln("Error creating request:", err)
		return
	}
	req.Header.Set(utils.AuthHeader, utils.FormatToken(mediator.Token))
	req.Header.Set(utils.ClientIdHeader, mediator.ClientId)

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	// Parse the response body to get the new access token
	var result2 map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result2)

	for k, v := range result2 {
		fmt.Println(k, ": ", v)
	}

	rewardURL, err := url.Parse("https://api.twitch.tv/helix/channel_points/custom_rewards")
	if err != nil {
		log.Fatalln(errURLMsg, err)
		return
	}

	params = url.Values{}
	params.Add("broadcaster_id", broadcasterId)
	rewardURL.RawQuery = params.Encode()

	req, err = http.NewRequest("GET", rewardURL.String(), nil)
	if err != nil {
		log.Fatalln("Error creating request:", err)
		return
	}

	req.Header.Set(utils.AuthHeader, "Bearer "+mediator.Token)
	req.Header.Set(utils.ClientIdHeader, mediator.ClientId)

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	// Parse the response body to get the new access token
	var result3 map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result3)

	fmt.Println("\n === Reward === \n")
	for k, v := range result3 {
		fmt.Println(k, ": ", v)
	}

	//remove the vip granted to victim

	//add the vip to the thief

}
