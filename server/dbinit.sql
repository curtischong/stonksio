SET sql_safe_updates = FALSE;

USE defaultdb;
DROP DATABASE IF EXISTS stonksio CASCADE;
CREATE DATABASE IF NOT EXISTS stonksio;

USE stonksio;

DROP TABLE IF EXISTS post;
CREATE TABLE post (
  id UUID PRIMARY KEY,
  username TEXT NOT NULL,
  userPicUrl TEXT NOT NULL DEFAULT '',
  body TEXT NOT NULL,
  timestamp TIMESTAMP NOT NULL
);

DROP TABLE IF EXISTS price;
CREATE TABLE price (
  id UUID PRIMARY KEY,
  asset TEXT,
  price DECIMAL,
  timestamp TIMESTAMP
);