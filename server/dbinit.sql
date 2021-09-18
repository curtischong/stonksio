SET sql_safe_updates = FALSE;

USE defaultdb;
DROP DATABASE IF EXISTS stonksio CASCADE;
CREATE DATABASE IF NOT EXISTS stonksio;

USE stonksio;

CREATE TABLE post (
  id UUID PRIMARY KEY,
  message TEXT
);