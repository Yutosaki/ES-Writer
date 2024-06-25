package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt string
}

var db *sql.DB

var cognitoRegion string
var clientId string
var jwksURL string
var region string
var accessID string
var secretAccessKey string
var sessionToken string

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handler called...")

	var user User
	err := db.QueryRow("SELECT id, username, email, created_at FROM users WHERE username = $1", "test").Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		fmt.Printf("error in query: %s", err)
		return
	}

	fmt.Fprintf(w, "ID: %s, Name: %s, Email: %s, CreatedAt: %s\n", user.ID, user.Name, user.Email, user.CreatedAt)
}

func main() {
	fmt.Println("server started...")
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var dbHost string = os.Getenv("DB_HOST")
	var dbUser string = os.Getenv("DB_USER")
	var dbPassword string = os.Getenv("DB_PASSWORD")
	var dbName string = os.Getenv("DB_NAME")
	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		log.Fatalf("環境変数が設定されていません: dbHost=%s, dbUser=%s, dbPassword=%s, dbName=%s", dbHost, dbUser, dbPassword, dbName)
	}

	cognitoRegion = os.Getenv("COGNITO_REGION")
	clientId = os.Getenv("COGNITO_CLIENT_ID")
	jwksURL = os.Getenv("TOKEN_KEY_URL")
	if cognitoRegion == "" || clientId == "" || jwksURL == "" {
		log.Fatalf("環境変数が設定されていません: cognitoRegion=%s, clientId=%s, jwksURL=%s", cognitoRegion, clientId, jwksURL)
	}

	region = os.Getenv("AWS_DEFAULT_REGION")
	accessID = os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	sessionToken = os.Getenv("AWS_SESSION_TOKEN")
	if region == "" || accessID == "" || secretAccessKey == "" {
		log.Fatalf("環境変数が設定されていません: region=%s, accessID=%s, secretAccessKey=%s", region, accessID, secretAccessKey)
	}

	// 環境変数からDB接続情報を取得
	db, err = sql.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName))
	if err != nil {
		fmt.Println("error in db connection")
		log.Fatal(err)
	}
	defer db.Close()

	// 直接DB接続情報を記述
	// db, err = sql.Open("postgres", "host=db user=postgres password=postgres dbname=testdb sslmode=disable")
	// if err != nil {
	// 	fmt.Println("error in db connection")
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	http.HandleFunc("/", handler)
	http.HandleFunc("/signin", signin)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/checkEmail", checkEmail)
	http.HandleFunc("/resendEmail", resendEmail)
	http.HandleFunc("/welcome", welcome)
	http.HandleFunc("/saveprofile", saveProfile)
	http.HandleFunc("/getAnswers", processQuestionsWithAI)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
