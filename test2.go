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
)

type tagNum struct{
    text string
    count int
}
type tagNums []tagNum
type tagImg struct{
    text string
    img []string
}
type tagStr struct{
    text string
    count int
    images []string
    point int
}

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

func Sort(mstr []tagStr) []tagStr { 
    i := 1
    for {
        if mstr[i - 1].point < mstr[i].point {
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

func arrayToStruct(array []tagImg) []tagStr {
    m := map[string]int{}
    var newTagStr tagStr
    var m_struct []tagStr
    var m_structt []tagStr
    for _, arr := range array {
        _, ok := m[arr.text]
        if ok {
            m[arr.text] = m[arr.text] + 1
        }else{
            m[arr.text] = 1
        }
    }
    for tex, cou := range m {
        newTagStr.count = cou
        newTagStr.text = tex
        m_struct = append(m_struct, newTagStr)
    }
    for _, tagim := range m_struct {
        for _, artagim := range array {
            if tagim.text == artagim.text {
                // tagim.images = append(tagim.images, artagim.img)
                for _, img := range artagim.img {
                    tagim.images = append(tagim.images, img)
                }
            }
        }
        tagim.point = tagim.count * len(tagim.images)
        m_structt = append(m_structt, tagim)
    }
    m_structt = compress(m_structt)
    return Sort(m_structt)
}
func compress(array []tagStr) []tagStr {
    for i, _ := range array {
        for j := i+1; j < len(array); j++ {
            if array[i].text == array[j].text{
                array[j].point = 0
            }
        }
    }
    return array
} 

func main() {
    loadEnv()

    api := getTwitterApi()

    v := url.Values{}
    v.Set("count", "100")
    v.Set("lang", "ja")
    var tags_form []tagImg
    var tag_form tagImg
    var images []string
    i := 1
    go func() {
        // tags_form = nil
        // tag_form.text = ""
        // tag_form.img = nil
        // images = nil
        for {searchResult, _ := api.GetSearch("%23", v)
            for _, tweet := range searchResult.Statuses {
                images = nil
                medias := tweet.Entities.Media
                // fmt.Println(medias)
                for _, media := range medias {
                    images = append(images, media.Media_url)
                }
                tags := tweet.Entities.Hashtags
                for _, tag := range tags {
                    if tag.Text != ""  {
                        tag_form.text = tag.Text
                        // tag_form.img = append(tag_form.img, images)
                        // for _, _ = range images {
                        for _, img := range images{
                            tag_form.img = append(tag_form.img, img)
                            // tag_form.img = append(tag_form.img, "あああ")
                        }
                        tags_form = append(tags_form, tag_form)
                        tag_form.img = nil
                    }
                }
                // fmt.Println(images)
            }
            fmt.Println(i)
            if i % 20 == 0 {
                for k := 0; k < 10; k++{
                    fmt.Println(arrayToStruct(tags_form)[k])
                }
                tags_form = nil
            }
            i = i + 1
            time.Sleep(5 * time.Second)
        }
    }()

    // シグナル用のチャネル定義
    quit := make(chan os.Signal)
    // 受け取るシグナルを設定
    signal.Notify(quit, os.Interrupt)
    <-quit // ここでシグナルを受け取るまで以降の処理はされない

    // シグナルを受け取った後にしたい処理を書く
    fmt.Println("終了しました。")
    // for i = 0; i < 11; i++ {
    //     fmt.Println(arrayToHash(tags_form)[i])
    // }

}