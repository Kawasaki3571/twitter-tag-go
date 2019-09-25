package main

import (
    "fmt"
    "log"
    "net/url"
    "os"
    "os/signal"

    "github.com/ChimeraCoder/anaconda"
    "github.com/joho/godotenv"
    "time"
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

func arrayToHash(array []string) map[string]int {
    m := map[string]int{}
    for _, arr := range array {
        _, ok := m[arr]
        if ok {
            m[arr] = m[arr] + 1
        }else{
            m[arr] = 1
        }
    }
    return m
}

func main() {
    loadEnv()

    api := getTwitterApi()

    v := url.Values{}
    v.Set("count", "100")
    v.Set("lang", "ja")
    var tags_form []string
    i := 1
    go func() {
        for {searchResult, _ := api.GetSearch("%23", v)
            for _, tweet := range searchResult.Statuses {
                // fmt.Println(tweet.FullText)
                // fmt.Println(tweet.User.ScreenName)
                tags := tweet.Entities.Hashtags
                for _, tag := range tags {
                    if tag.Text != ""  {
                        tags_form = append(tags_form, tag.Text)
                    }
                }
                // fmt.Println("https://twitter.com/" + tweet.User.ScreenName + "/status/" + tweet.IdStr )
                // fmt.Println("========================================")
            }
            fmt.Println(i)
            i = i + 1
            time.Sleep(80 * time.Second)
        }
    }()
    
    // シグナル用のチャネル定義
    quit := make(chan os.Signal)
    // 受け取るシグナルを設定
    signal.Notify(quit, os.Interrupt)
    <-quit // ここでシグナルを受け取るまで以降の処理はされない

    // シグナルを受け取った後にしたい処理を書く
    fmt.Println(arrayToHash(tags_form))
}

