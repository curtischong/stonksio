SET sql_safe_updates = FALSE;

USE defaultdb;
DROP DATABASE IF EXISTS stonkst CASCADE;
CREATE DATABASE IF NOT EXISTS stonkst;

USE stonkst;

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
  tradeprice DECIMAL,
  timestamp TIMESTAMP
);

DROP TABLE IF EXISTS wallet;
CREATE TABLE wallet(
   id UUID PRIMARY KEY,
   username TEXT,
   asset TEXT,
   balance DECIMAL
);

DROP TABLE IF EXISTS position;
CREATE TABLE position(
    id UUID PRIMARY KEY,
    username TEXT,
    size DECIMAL,
    price DECIMAL,
    timestamp TIMESTAMP
);