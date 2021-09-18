SET sql_safe_updates = FALSE;

USE defaultdb;
DROP DATABASE IF EXISTS stonksio CASCADE;
CREATE DATABASE IF NOT EXISTS stonksio;

USE stonksio;

DROP TABLE IF EXISTS post
CREATE TABLE post (
  id UUID PRIMARY KEY,
  username TEXT,
  userPicUrl TEXT,
  body TEXT
);

DROP TABLE IF EXISTS ohlc
CREATE TABLE ohlc(
  id UUID PRIMARY KEY,
  asset TEXT,
  open DECIMAL,
  high DECIMAL,
  low DECIMAL,
  close DECIMAL,
);