package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"cache"
)

func init() {
	//create the cache file
	localUserCache, err := os.UserCacheDir()
	if err != nil {
		fmt.Println("Error getting user cache directory:", err)
		return
	}
	cachePath := filepath.Join(localUserCache, "toneeuw")

	if _, err := os.Stat(cachePath); os.IsNotExist(err) {

		res := os.MkdirAll(cachePath, os.ModePerm)
		if res != nil {
			fmt.Println("Error creating cache directory:", res)
			return
		}
		fmt.Println("Creating cache directory")
	}

	cachePath = filepath.Join(cachePath, "cache.json")

	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		_, err := os.Create(cachePath)
		if err != nil {

			fmt.Println("Error creating cache file:", err)
			panic(err)

		}

		fmt.Println("Creating cache file")
	}
}

func main() {

	localUserCache, err := os.UserCacheDir()
	if err != nil {
		fmt.Println("Error getting user cache directory:", err)
		return
	}

	cachePath := filepath.Join(localUserCache, "toneeuw", "cache.json")

	cache, err := cache.New(cachePath)



	cache.User.Credentials()

	baseURL, err := url.Parse("https://api.twitch.tv/helix/channel_points/custom_rewards")
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	// Set the query parameters
	params := url.Values{}
	params.Add("broadcaster_id", "drhurell")
	//params.Add("id", "your_reward_id")
	//params.Add("only_manageable_rewards", "true")
	baseURL.RawQuery = params.Encode()

	// Create the HTTP request
	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the headers
	req.Header.Set("Authorization", "Bearer "+cache.User.Token)
	req.Header.Set("Client-Id", "96v4gygg6e8ljx4ah9s32rw7bz2p55")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error kinsending request:", err)
		return
	}
	test2 := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&test2)
	if test2 == nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	println(cache.User.Token)
	println(cache.User.RefreshToken)
	//fmt.Println(test2)

	//replace the file content with the new vip variable
	cache.Save(cachePath)

}
