package main

import (
	"log"
	"os"
	"path/filepath"

	"cache"
)

func init() {
	//create the cache file
	localUserCache, err := os.UserCacheDir()
	if err != nil {
		log.Fatalln("Error getting user cache directory:", err)
		return
	}
	cachePath := filepath.Join(localUserCache, "toneeuw")

	if _, err := os.Stat(cachePath); os.IsNotExist(err) {

		res := os.MkdirAll(cachePath, os.ModePerm)
		if res != nil {
			log.Fatalln("Error creating cache directory:", res)
			return
		}
		log.Println("Creating cache directory")
	}

	cachePath = filepath.Join(cachePath, "cache.json")

	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		_, err := os.Create(cachePath)
		if err != nil {

			log.Fatalln("Error creating cache file:", err)
			panic(err)

		}

		log.Println("Creating cache file")
	}
}

func main() {

	localUserCache, err := os.UserCacheDir()
	if err != nil {
		log.Fatalln("Error getting user cache directory:", err)
		return
	}

	cachePath := filepath.Join(localUserCache, "toneeuw", "cache.json")

	cache, err := cache.New(cachePath)

	cache.User.Credentials()

	cache.VIP.Steal(&cache.User, nil)

	cache.Save(cachePath)

}
