package userController

import (
	"abwaab/database"
	helper "abwaab/helper"
	models "abwaab/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// User Sign up function

func UserSignup(response http.ResponseWriter, request *http.Request){
	response.Header().Set("Content-Type","application/json")
	// define User model 
	var user models.User
	// check db
 	json.NewDecoder(request.Body).Decode(&user)
	email := user.Email
	password := user.Password

	if email == ""  || password == "" {
		helper.ResponseSend("201","Email and password are required.",response)
		return 	
	} 
	if !helper.ValidMailAddress(email) {
		helper.ResponseSend("201", "Please enter valid email address.",response)
		return 
	} 
	
	is_exist , er := database.CheckUserAlreadyExist(email)
	
	if is_exist && er == nil {
		helper.ResponseSend("201", "The user is already exist with "+email+ ". Please enter different email.",response)
		return 	
	} else {
		
		db, dbError := database.DBinstance()
        defer db.Close()
		database.CheckError(dbError)

		user.Password = helper.HashPassword(user.Password)

		insertStmt := `INSERT INTO "users"("email", "password", "created_on", "updated_on") values ($1, $2, $3, $4)`
		_, insertError := db.Exec(insertStmt,email, user.Password, time.Now(), time.Now())
    database.CheckError(insertError)
		if (insertError == nil){
			helper.ResponseSend("200","You are successfully sign up with email "+email+ ". Please check your email to continue.",response)
			return 	
		} else {
			helper.ResponseSend("400","Something went wrong please try again later.",response)
			return 	
		}

	}
}


// User Login function

func UserLogin(response http.ResponseWriter, request *http.Request){
  response.Header().Set("Content-Type","application/json")

  var user models.User
	// check db
	json.NewDecoder(request.Body).Decode(&user)
	email := user.Email
	password := user.Password

	if email == ""  || password == "" {
		helper.ResponseSend("201","Email and password are required.",response)
		return 	
	} 
	if !helper.ValidMailAddress(email) {
		helper.ResponseSend("201", "Please enter valid email address.",response)
		return 
	} 
	
	userData , er := database.GetUserInfo(email)	
	if er!=nil{
	  response.WriteHeader(http.StatusInternalServerError)
	  response.Write([]byte(`{"message":"`+er.Error()+`"}`))
	  return
  }
	 
		
	userPass:= []byte(user.Password)
	dbPass:= []byte(userData.Password)

	passErr:= bcrypt.CompareHashAndPassword(dbPass, userPass)

	if passErr != nil{
	  response.Write([]byte(`{"response":"Wrong Password!"}`))
	  return
  }
	
	
  jwtToken, err := helper.GenerateJWT(email)
	//jwtToken, err := GenerateJWT(email)
  
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"`+err.Error()+`"}`))
		return
  }
  response.Write([]byte(`{"token":"`+jwtToken+`"}`))
	return 	
  
}

// CreateTweet is used to create user tweet
func CreateTweet(response http.ResponseWriter, request *http.Request){
	response.Header().Set("Content-Type","application/json")
	// Authorize token
	clientToken :=	request.Header.Get("token")
	
	if clientToken == "" {
		response.WriteHeader(http.StatusUnauthorized)
		return
	}
	_, returnFlag := helper.ExtractClaims(clientToken)
	fmt.Println(returnFlag)
	if !returnFlag {
		response.WriteHeader(http.StatusUnauthorized)
		return
	}
	// End Authorize token
	
	var data models.TweetData
	// check db
	json.NewDecoder(request.Body).Decode(&data)
	tweetDescription := data.Description
	
	if tweetDescription == "" {
		helper.ResponseSend("201","Please enter description to post on tw.",response)
		return 	
	} 

	creds := helper.Credentials{
		AccessToken:       helper.GoDotEnvVariable("ACCESS_TOKEN"),
		AccessTokenSecret: helper.GoDotEnvVariable("ACCESS_TOKEN_SECRET"),
		ConsumerKey:       helper.GoDotEnvVariable("CONSUMER_KEY"),
		ConsumerSecret:    helper.GoDotEnvVariable("CONSUMER_SECRET"),
	}

	client, err := helper.GetClient(&creds)
	if err != nil {
			log.Println("Error getting Twitter Client")
			log.Println(err)
	}

	tweet, resp, err := client.Statuses.Update(tweetDescription, nil)

	
	if err != nil {
	    log.Println(err)
	}
	
	log.Printf("%+v\n He: ", resp)
	log.Printf("%+v\n", tweet)

    response.Write([]byte(`{"token":"sucess"}`))
	return 
}

// findAndBulkInsert is used to create user tweet

func findAndBulkInsert(hastag string, userId int16){

	creds := helper.Credentials{
		AccessToken:       helper.GoDotEnvVariable("ACCESS_TOKEN"),
		AccessTokenSecret: helper.GoDotEnvVariable("ACCESS_TOKEN_SECRET"),
		ConsumerKey:       helper.GoDotEnvVariable("CONSUMER_KEY"),
		ConsumerSecret:    helper.GoDotEnvVariable("CONSUMER_SECRET"),
	}

	client, err := helper.GetClient(&creds)
	if err != nil {
			log.Println("Error getting Twitter Client")
			log.Println(err)
	}
	searchParams := &twitter.SearchTweetParams{
		Query:      hastag,
		Count:      50,
		ResultType: "recent",
		Lang:       "en",
		TweetMode: "extended",

	}
	searchResult, _, err := client.Search.Tweets(searchParams)
	db, dbError := database.DBinstance()
    defer db.Close()
	database.CheckError(dbError)

	vals := []interface{}{}
	sqlStr := "INSERT INTO user_search_tweets(user_id, tweet, created_on, updated_on) VALUES "
	for _, tweet := range searchResult.Statuses {

		var tweetText string
		if tweet.FullText != "" {
			tweetText = tweet.FullText 
		}	else {
			tweetText = tweet.Text
		}
		sqlStr += "(?, ?, ?, ?),"
    vals = append(vals, userId, tweetText,  time.Now(), time.Now())
		
	}
	
	//trim the last ,
	sqlStr = strings.TrimSuffix(sqlStr, ",")

	//Replacing ? with $n for postgres
	sqlStr = helper.ReplaceSQL(sqlStr, "?")
	//prepare the statement
	stmt, _ := db.Prepare(sqlStr)
	//format all vals at once
	_, er := stmt.Exec(vals...)
	//fmt.Println("Data",res)
  database.CheckError(er)

}

// SearchAndSaveTweet is used to create user tweet
func SearchAndSaveTweet(response http.ResponseWriter, request *http.Request){
	response.Header().Set("Content-Type","application/json")
	// Authorize token
	clientToken :=	request.Header.Get("token")
	fmt.Println("token::: ",clientToken)
	if clientToken == "" {
		response.WriteHeader(http.StatusUnauthorized)
		return
	}
	userData, returnFlag := helper.ExtractClaims(clientToken)
	
	if !returnFlag {
		response.WriteHeader(http.StatusUnauthorized)
		return
	}
	
	// End Authorize token
	
	//var insertWait sync.WaitGroup
	var data models.TweetData
	// check db
	json.NewDecoder(request.Body).Decode(&data)
	hastag := data.Hashtag
	fmt.Println(data)
	if hastag == "" {
		helper.ResponseSend("201","Please enter description to post on tw.",response)
		return 	
	} 
	// Searching tweet 
	//insertWait.Add(1)
    email :=userData["email"].(string)
    userId , er := database.GetUserId(email)	
	if er!=nil{
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"`+er.Error()+`"}`))
		return
	}
	go findAndBulkInsert(hastag, userId)	// add in back groud 
	//insertWait.Wait() 
	//response.Write([]byte(`{"token":"sucess"}`))
	helper.ResponseSend("200","Success.",response)
	return 
}


func listRecords(page int) []interface{}{
    db, _ := database.DBinstance()
	limit := 10
	offset := limit * (page - 1)

    rows, err := db.Query("SELECT * FROM user_search_tweets ORDER BY id LIMIT $2 OFFSET $1",  offset, limit)
    if err != nil {
        // handle this error better than this
        panic(err)
    }
    defer rows.Close()
    
    record := []interface{}{}
    for rows.Next() {
        var id int16
        var tweet string
        var user_id string
        var created_at string
        var updated_at string
        
        err = rows.Scan(&id, &user_id, &tweet, &created_at , &updated_at)
        if err != nil {
        // handle this error
        panic(err)
        }
        record = append(record, id, user_id,  tweet)
    }
    // get any error encountered during iteration
    err = rows.Err()
    if err != nil {
        panic(err)
    }
    return record
}

// SearchAndSaveTweet is used to create user tweet
func ListTweet(response http.ResponseWriter, request *http.Request){
	response.Header().Set("Content-Type","application/json")
	// Authorize token
	clientToken :=	request.Header.Get("token")
	
	if clientToken == "" {
		response.WriteHeader(http.StatusUnauthorized)
		return
	}
	userData, returnFlag := helper.ExtractClaims(clientToken)
	
	if !returnFlag {
		response.WriteHeader(http.StatusUnauthorized)
		return
	}
	fmt.Println(userData["email"].(string))
	// End Authorize token
    email := userData["email"].(string)
	
	user , er := database.GetUserId(email)	
	if er!=nil{
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"`+er.Error()+`"}`))
		return
	}
			
	fmt.Println("Data", user)
    page := request.URL.Query().Get("page")
    limit := request.URL.Query().Get("limit")
	fmt.Println("Data", page)

    if page == "" {
        page = "1"
    }
    intPage, _ := strconv.Atoi(page)
    
    if limit == "" {
        limit = "10"
    }

    fmt.Println(listRecords(intPage))
    result := listRecords(intPage)
    s, _ := json.Marshal(result) 
	response.Write(s)
	return 
}
