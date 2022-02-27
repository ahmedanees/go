package helper

import (
	config "abwaab/config"
	"strconv"
	"strings"
	"time"

	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// use godot package to load/read the .env file and
// return the value of the key
func GoDotEnvVariable(key string) string {

  // load .env file
  err := godotenv.Load(".env")

  if err != nil {
    log.Fatalf("Error loading .env file")
  }

  return os.Getenv(key)
}

func ValidMailAddress(address string) (bool) {
	_, err := mail.ParseAddress(address)
	if err != nil {
			return false
	}
	return  true
}

// sending reponse function added 
func ResponseSend(message string, code string, response http.ResponseWriter){
	payload := map[string]string{"code": code, "message": message}
	json.NewEncoder(response).Encode(payload)
}


//HashPassword is used to encrypt the password before it is stored in the DB
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
			log.Panic(err)
	}

	return string(bytes)
}

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	email string `json:"email"`
	jwt.StandardClaims
}


//Generate JWT
func GenerateJWT(email string )(string,error){
	atClaims := jwt.MapClaims{}
  atClaims["authorized"] = true
  atClaims["email"] = email
  atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
  at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	// Declare the token with the algorithm used for signing, and the claims
	token, err := at.SignedString([]byte(config.SECRET_KEY))
	// Create the JWT string
	if err !=nil{
		log.Println("Error in JWT token generation")
		return "",err
	}
	return token, nil
}

//Extract claims from the JWT token

func ExtractClaims(tokenStr string) (jwt.MapClaims, bool) {
	hmacSecretString := config.SECRET_KEY
	hmacSecret := []byte(hmacSecretString)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			 // check token signing method etc
			 return hmacSecret, nil
	})

	if err != nil {
			log.Printf("Invalid JWT Token")
			return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			return claims, true
	} else {
			log.Printf("Invalid JWT Token", tokenStr)
			return nil, false
	}
}


// Credentials stores all of our access/consumer tokens
// and secret keys needed for authentication against
// the twitter REST API.
type Credentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

// getClient is a helper function that will return a twitter client
// that we can subsequently use to send tweets, or to stream new tweets
// this will take in a pointer to a Credential struct which will contain
// everything needed to authenticate and return a pointer to a twitter Client
// or an error
func GetClient(creds *Credentials) (*twitter.Client, error) {
	// Pass in your consumer key (API Key) and your Consumer Secret (API Secret)
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	// Pass in your Access Token and your Access Token Secret
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// Verify Credentials
	verifyParams := &twitter.AccountVerifyParams{
			SkipStatus:   twitter.Bool(true),
			IncludeEmail: twitter.Bool(true),
	}

	// we can retrieve the user and verify if the credentials
	// we have used successfully allow us to log in!
	_, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
			return nil, err
	}

	//log.Printf("User's ACCOUNT:\n%+v\n", user)
	return client, nil
}


// ReplaceSQL replaces the instance occurrence of any string pattern with an increasing $n based sequence
func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		 old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}