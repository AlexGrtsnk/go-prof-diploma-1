DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS cfg;
CREATE TABLE cfg (id INTEGER PRIMARY KEY, flagRunAddr VARCHAR(100), flnm VARCHAR(100));
CREATE TABLE users (id INTEGER PRIMARY KEY, lgn VARCHAR(100), psw VARCHAR(100), token VARCHAR(300), wtdh FLOAT DEFAULT 0, balance FLOAT DEFAULT 0);
CREATE TABLE orders (id INTEGER PRIMARY KEY, nmb VARCHAR(100), sts VARCHAR(100), token VARCHAR(300), accural FLOAT DEFAULT -1, ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP, sumbals FLOAT DEFAULT 0);