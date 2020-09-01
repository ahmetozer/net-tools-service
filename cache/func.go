package cache

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ahmetozer/net-tools-service/cache/memory"
)

var (
	storage       Storage
	cacheDuration = "10s"
)

func init() {
	storage = memory.NewStorage()

	cacheDuration, ok := os.LookupEnv("cache")
	if ok {
		if _, err := time.ParseDuration(cacheDuration); err != nil {
			log.Fatalln("\033[1;31ma" + fmt.Sprint(err) + "\033[0m")
		}
	} else {
		log.Println("Environment variable \"cache\" is not set. Default cache value is 10s .")
	}
}
func Set(storageHash string, cacheableString string) string {

	content := storage.Get(storageHash)
	if content != nil {
		return string(content)
	}
	content = []byte(cacheableString)

	if d, err := time.ParseDuration(cacheDuration); err == nil {
		storage.Set(storageHash, content, d)
		return cacheableString
	} else {
		fmt.Println(err)
		return cacheableString
	}

}

func IsCached(hash string) bool {
	return storage.Get(hash) != nil
}

func Get(hash string) string {
	return string(storage.Get(hash))
}
