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
    // "sort"
    "strings"
)

type tagNum struct{
    text string
    count int
    // image_count int
}
type tagNums []tagNum

type tagImg struct{
    text string
    image_url []string
    image_count int
}

type tagImgs []tagImg

type tagStr struct{
    text string
    tag_count int
    // image_url []anaconda.EntityMedia
    image_count int
    image_url []string
}

type tagStrs []tagImg

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

func Sort(mstr []tagNum) []tagNum { 
    i := 1
    for {
        if mstr[i - 1].count < mstr[i].count {
            mstr[i - 1], mstr[i] = mstr[i], mstr[i -1]
            i = 0
        }
        i = i + 1
        if i == len(mstr) {
            break
        }
    }
    return mstr
}

func arrayToStruct(array []string) []tagNum {
    m := map[string]int{}
    var newTagNum tagNum
    var m_struct []tagNum
    for _, arr := range array {
        _, ok := m[arr]
        if ok {
            m[arr] = m[arr] + 1
        }else{
            m[arr] = 1
        }
    }
    for tex, cou := range m {
        newTagNum.count = cou
        newTagNum.text = tex
        m_struct = append(m_struct, newTagNum)
    }
    // sort.Sort(m_struct)
    return Sort(m_struct)
}

func getImgBySlice(slice []string) []string {
    var images []string
    for _, link := range slice {
        if strings.Contains(link, "http://pbs.twimg.com/media") {
            images = append(images, link)
        }
    }
    return images
}

func main() {
    loadEnv()

    api := getTwitterApi()

    v := url.Values{}
    v.Set("count", "100")
    v.Set("lang", "ja")
    var tags_form tagNums
    var tag_form tagImg
    i := 1
    go func() {
        for {searchResult, _ := api.GetSearch("%23", v)
            for _, tweet := range searchResult.Statuses {
                tags := tweet.Entities.Hashtags
                for _, tag := range tags { 
                    tag_form.image_count = 0
                    if tag.Text != ""  {
                        tag_form.text = tag.Text
                        // tag_form.image_url = append(tag_form.image_url, tweet.Entities.Media)
                        // for _, urls := range tweet.Entities.Media {
                        //     tag_form.image_url = append(tag_form.image_url, tweet.Entities.Media)
                        // }
                        tag_form.image_count = tag_form.image_count + len(getImgBySlice(tweet.Entities.Media))
                        tags_form = append(tags_form, tag_form)
                    }
                }
                // fmt.Println(tweet.Entities.Media)
            }
            fmt.Println(i)
            i = i + 1
            time.Sleep(10 * time.Second)
        }
    }()

    // シグナル用のチャネル定義
    quit := make(chan os.Signal)
    // 受け取るシグナルを設定
    signal.Notify(quit, os.Interrupt)
    <-quit // ここでシグナルを受け取るまで以降の処理はされない

    // シグナルを受け取った後にしたい処理を書く

    fmt.Println(arrayToStruct(tags_form))

}

