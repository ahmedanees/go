package database

import (
	helper "abwaab/helper"
	models "abwaab/models"
	"database/sql"
	"fmt"
)

var (
  host     = helper.GoDotEnvVariable("DB_HOST")
  port     = helper.GoDotEnvVariable("DB_PORT")
  user     = helper.GoDotEnvVariable("DB_USER")
  password = helper.GoDotEnvVariable("DB_PASSWORD")
  dbname   = helper.GoDotEnvVariable("DATABASE_NAME")
)


type User struct {  
	email      string
	password     string
}
func dsn() string{
	
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
	"password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)
	return psqlInfo
}
//DBinstance func

func DBinstance() (*sql.DB, error) {  
    // connection string
    psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
         
    // open database
    db, err := sql.Open("postgres", psqlconn)
    CheckError(err)
    // check db
    err = db.Ping()
    CheckError(err)
 
    fmt.Println("Connected!")

	return db, nil
}

func CheckUserAlreadyExist(email string) (bool, error) {
	var id int
	db, error := DBinstance()
	CheckError(error)
	// Query for a value based on a single row.
	
	row := db.QueryRow("SELECT id FROM users WHERE email = $1 ",email)
	err := row.Scan(&id)
	
	if err != nil && err != sql.ErrNoRows {
    // log the error
		return false, fmt.Errorf("error %s:", err)
	}
  if id > 0 {
	  return true, nil
  } else {
    return false, nil
  }  
}


func GetUserInfo(email string) (models.User, error) {
	var user models.User
	var id int16
	var password string
	db, error := DBinstance()
	CheckError(error)
	// Query for a value based on a single row.
	row := db.QueryRow("SELECT id, email,password FROM users WHERE email = $1 ",email)
	err := row.Scan(&id, &email, &password)
	fmt.Println(err)
	if err != nil && err != sql.ErrNoRows {
    // log the error
		return user, fmt.Errorf("error %s:", err)
	}
  user.Id = id
	user.Email = email
	user.Password = password
  return user, nil
}

func GetUserId(email string) (int16, error) {
	var id int16
	var password string
	db, error := DBinstance()
	CheckError(error)
	// Query for a value based on a single row.
	row := db.QueryRow("SELECT id, email,password FROM users WHERE email = $1 ",email)
	err := row.Scan(&id, &email, &password)
	fmt.Println(err)
	if err != nil && err != sql.ErrNoRows {
    // log the error
		return id, fmt.Errorf("error %s:", err)
	}
  return id, nil
}

func CheckError(err error) {
	if err != nil {
			panic(err)
	}
}