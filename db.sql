DROP TABLE IF EXISTS users;

CREATE TABLE users (
 id serial PRIMARY KEY,
 email VARCHAR ( 255 ) UNIQUE NOT NULL,
 password Text NOT NULL,
 created_on TIMESTAMP NOT NULL,
 updated_on TIMESTAMP NOT NULL
        
);

DROP TABLE IF EXISTS user_search_tweets;
CREATE TABLE user_search_tweets (
 id serial PRIMARY KEY,
 user_id INTEGER REFERENCES users (id),
 tweet Text NOT NULL,
 created_on TIMESTAMP NOT NULL,
 updated_on TIMESTAMP NOT NULL

);
