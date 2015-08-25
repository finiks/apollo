package main

import (
    "encoding/json"
    "encoding/xml"
    "net/http"
    "io/ioutil"
    "strings"
    "strconv"
)

func (t *EntriesTab) fetchTags() {
    title := strings.Replace(t.slice[t.cursor].Title, " ", "+", -1)
    urls := map[string]string{
        "omdb": "http://www.omdbapi.com/?s={TITLE}&type=movie&y=&plot=full&r=json",
        "hummingbird":"http://hummingbird.me/api/v1/search/anime?query={TITLE}",
        "gamesdb":"http://thegamesdb.net/api/GetGamesList.php?name={TITLE}",
        "googlebooks":"https://www.googleapis.com/books/v1/volumes?q={TITLE}&projection=lite&printType=books&maxResults=10",
    }
    url := strings.Replace(urls[t.taggingAPI], "{TITLE}", title, 1)
    t.a.logDebug(url)

    res, err := http.Get(url)
    if err != nil {
        t.a.logError(err.Error())
        return
    }
    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        t.a.logError(err.Error())
        return
    }

    switch t.taggingAPI {
    case "omdb":
        t.fetchOMDBTags(&body)
    case "hummingbird":
        t.fetchHummingbirdTags(&body)
    case "gamesdb":
        t.fetchGamesDBTags(&body)
    case "googlebooks":
        t.fetchGoogleBooksTags(&body)
    }

    if len(t.search) > 0 {
        t.view = "tag"
    }
}

func (t *EntriesTab) fetchOMDBTags(body *[]byte) {
    type OMDBEntry struct {
        Title string
        Year string
        ImdbID string
    }

    type OMDBData struct {
        Search []OMDBEntry
    }

    var data OMDBData
    err := json.Unmarshal(*body, &data)
    if err != nil {
        t.a.logError(err.Error())
        return
    }

    for i := 0; i < len(data.Search); i++ {
        t.search = append(t.search, Entry{
            Title: data.Search[i].Title,
            Year: data.Search[i].Year,
            TagID: data.Search[i].ImdbID,
        })
    }
}

func (t *EntriesTab) fetchHummingbirdTags(body *[]byte) {
    type HummingbirdEntry struct {
        Id int
        Title string
        Episode_count int
        Started_airing string
    }

    var data []HummingbirdEntry
    err := json.Unmarshal(*body, &data)
    if err != nil {
        t.a.logError(err.Error())
        return
    }

    for i := 0; i < len(data); i++ {
        if i < 10 {
            releaseDate := strings.Split(data[i].Started_airing, "-")
            t.search = append(t.search, Entry{
                Title: data[i].Title,
                TagID: strconv.Itoa(data[i].Id),
                Year: releaseDate[0],
                EpisodeTotal: data[i].Episode_count,
                EpisodeDone: data[i].Episode_count,
            })
        }
    }
}

func (t *EntriesTab) fetchGamesDBTags(body *[]byte) {
    type Game struct {
        id string
        GameTitle string
        ReleaseDate string
        Platform string
    }

    type GamesDBData struct {
        XMLName xml.Name `xml:"Data"`
        Game []Game
    }

    var data GamesDBData
    err := xml.Unmarshal(*body, &data)
    if err != nil {
        t.a.logError(err.Error())
        return
    }

    platforms := map[string]string{
        "Sony Playstation": "PS", "Sony Playstation 2": "PS2",
        "Sony Playstation 3": "PS3", "Sony Playstation 4": "PS4",
        "Sony PSP": "PSP", "Sony Playstation Vita": "VITA",
        "Microsoft Xbox": "XBOX", "Microsoft Xbox 360": "X360",
        "Microsoft Xbox One": "XONE",
        "Nintendo Entertainment System (NES)": "NES",
        "Super Nintendo (SNES)": "SNES", "Nintendo 64": "N64",
        "Nintendo GameCube": "NGC", "Nintendo DS": "NDS", "Nintendo 3DS": "3DS",
        "Nintendo Game Boy": "GB",
        "Nintendo Game Boy Color": "GBC", "Nintendo Game Boy Advance": "GBA",
        "Nintendo Wii": "WII", "Nintendo Wii U": "WIIU",
        "PC": "PC",
    }

    for i := 0; i < len(data.Game); i++ {
        if i < 10 {
            releaseDate := strings.Split(data.Game[i].ReleaseDate, "/")
            if len(releaseDate) > 2 {
                t.search = append(t.search, Entry{
                    Title: data.Game[i].GameTitle,
                    Year: releaseDate[2],
                    TagID: data.Game[i].id,
                    Info1: platforms[data.Game[i].Platform],
                })
            }
        }
    }
}

func (t *EntriesTab) fetchGoogleBooksTags(body *[]byte) {
    type GoogleBooksInfo struct {
        Title string
        Authors []string
        PublishedDate string
    }

    type GoogleBooksEntry struct {
        Id string
        VolumeInfo GoogleBooksInfo
    }

    type GoogleBooksData struct {
        Items []GoogleBooksEntry
    }

    var data GoogleBooksData
    err := json.Unmarshal(*body, &data)
    if err != nil {
        t.a.logError(err.Error())
        return
    }

    for i := 0; i < len(data.Items); i++ {
        if i < 10 {
            if len(data.Items[i].VolumeInfo.Authors) == 0 {
                data.Items[i].VolumeInfo.Authors = append(data.Items[i].VolumeInfo.Authors, "")
            }
            t.search = append(t.search, Entry{
                Title: data.Items[i].VolumeInfo.Title,
                TagID: data.Items[i].Id,
                Info1: data.Items[i].VolumeInfo.Authors[0],
            })
        }
    }
}