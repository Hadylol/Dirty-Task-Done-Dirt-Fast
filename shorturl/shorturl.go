package shorturl

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
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
		URL_Key: key,
	}

	mapData, err := structToMap(dataStruct)
	if err != nil {
		log.Printf("Error while converting Struct %+v To Map", dataStruct)
		return "", fmt.Errorf("failed to Convert Struct To Map %w", err)
	}

	exe := initDb.From("URLS").Insert([]map[string]interface{}{mapData}, false, "", "", "")
	response, _, err := exe.Execute()
	if err != nil {
		log.Printf("Error Inserting date into URL's Table : %v. Data: %+v", err, mapData)
		return "", fmt.Errorf("failed to inert data into the DB %w", err)
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

	url := strings.TrimPrefix(r.URL.Path, "/")
	fmt.Println("this is the  url from the Req ", url)

	redirectedURL, found, err := getURL(url)
	fmt.Println(found)
	if err != nil {
		fmt.Printf("Error occured while getting the URL: %s", redirectedURL)
		return
	}
	if !found {
		log.Printf("The URL %s Doesnt exist in the DB", redirectedURL)
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, redirectedURL, http.StatusFound)
}

func getURL(key string) (string, bool, error) {

	unparsed_URl, _, err := initDb.From("URLS").Select("url", "", false).Eq("url_key", key).Execute()
	if err != nil {
		log.Printf("Error occurred while fetching URL ")
		return "", false, fmt.Errorf("failed to fetch the data from URL's Table %w", err)
	}
	fmt.Println(len(unparsed_URl))

	if len(unparsed_URl) <= 2 {
		fmt.Printf("Parsed Response is empty for the Url_key %s", key)
		return "", false, nil
	}

	var urlRes []Url
	err = json.Unmarshal(unparsed_URl, &urlRes)
	if err != nil {
		log.Printf("Error occured while parssing the URL: %+v", unparsed_URl)
		return "", false, fmt.Errorf("failed parsing the data %w ", err)
	}

	var redirectedURL string
	for _, originalURL := range urlRes {
		redirectedURL = originalURL.URL
	}

	return redirectedURL, true, nil
}

func structToMap(data Url) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshling data of type Url : %+v. erorr: %v", data, err)
		return nil, fmt.Errorf("failed to marshal Url struct %w ", err)
	}

	var result map[string]interface{}

	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		log.Printf("Error Unmarshalling JSON data : %s. Error %v ", string(jsonData), err)
		return nil, fmt.Errorf("failed to unmarshal JSON data into map: %w ", err)
	}

	return result, nil
}
