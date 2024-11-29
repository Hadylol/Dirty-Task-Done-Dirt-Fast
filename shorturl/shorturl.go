package shorturl

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	db "telegram-bot/Db"
	"time"

	"github.com/google/uuid"
)

type Url struct {
	URL     string `json:"url"`
	URL_Key string `json:"url_key"`
}

type UrlResult struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt string    `json:"created_at"`
	URLKey    string    `json:"url_key"`
	URL       string    `json:"url"`
}
var initDb = db.InitDB()
func Shorturl() {

	http.HandleFunc("/", redirectURLHandler)

	//start the server
	log.Println("the server is running on PORT 4000...")

	log.Fatal(http.ListenAndServe(":4000", nil))
}

func ShortenURLHandler(URL string) (string, error) {
	//Generate unique key for the URL
	key := generateURLKey()
	dataStruct := Url{
		URL:     URL,
		URL_Key: "http://localhost:4000/" + key,
	}
	mapData, err := structToMap(dataStruct)
	if err != nil {
		log.Println("something went wrong while converting")
		return "", err
	}
	fmt.Println(mapData)
	exe := initDb.From("URLS").Insert([]map[string]interface{}{mapData}, false, "", "", "")
	response, _, err := exe.Execute()
	if err != nil {
		log.Println("Something went wrong while inserting")
		return "", err
	}
	var urls []UrlResult
	err = json.Unmarshal(response, &urls)
	if err != nil {
		fmt.Println("error while unmarshling structy")
		return "", nil
	}
	var URL_key_res string
	for _, url := range urls {
		URL_key_res = url.URLKey
	}

	return URL_key_res, nil
}

func generateURLKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.New(rand.NewSource(time.Now().UnixNano()))
	key := make([]byte, 6)
	for i := range key {
		key[i] = charset[rand.Intn(len(charset))]

	}
	//log.Println(string(key))
	return string(key)

}

 func redirectURLHandler(w http.ResponseWriter, r *http.Request) {
	//reterving the key for the URL
	URLKey :=

	orignalURL, found := getURL(URLKey)
	if !found {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, orignalURL, http.StatusFound)
}

func getURL(key string) (string, bool) {
	originalURL, found := KeyStore[key]
return originalURL, found
}

func structToMap(data Url) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
