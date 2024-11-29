package db

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

func InitDB() *supabase.Client {
	supabaseURL := "https://bmhsmlceplmnxgfmfrkp.supabase.co"
	err := godotenv.Load()
	if err != nil {
		panic("Error Loading .env file DB")
	}
	supabaseKey := os.Getenv("DB_KEY")
	client, err := supabase.NewClient(supabaseURL, supabaseKey, &supabase.ClientOptions{})
	if err != nil {
		log.Println("Something went Wrong init client : ", err)
	}
	return client
}
