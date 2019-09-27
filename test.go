package main

import (
    "fmt"
    "log"
    "net/url"
    "os"

    "github.com/ChimeraCoder/anaconda"
    "github.com/joho/godotenv"
)

func loadEnv() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
}

func getTwitterApi() *anaconda.TwitterApi {
    anaconda.SetConsumerKey(os.Getenv("CONSUMER_KEY"))
    anaconda.SetConsumerSecret(os.Getenv("CONSUMER_SECRET"))
    return anaconda.NewTwitterApi(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_TOKEN_SECRET"))
}

func main() {
    loadEnv()

    api := getTwitterApi()

    v := url.Values{}
    v.Set("count", "100")

    searchResult, _ := api.GetSearch("golang", v)
    i := 1
    j := 1
    // fmt.Println(i)
    for _, tweet := range searchResult.Statuses {
        // fmt.Println(i)
        // fmt.Println(tweet.Entities.Media)
        if len(tweet.Entities.Media) > 0 {
            for _, slice := range tweet.Entities.Media {
                fmt.Println(slice)
                fmt.Println(j)
                j = j + 1
            }
        }
        i = i + 1

    }
}
