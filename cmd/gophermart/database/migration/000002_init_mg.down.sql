DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS cfg;
CREATE TABLE cfg (id SERIAL PRIMARY KEY, flagRunAddr VARCHAR(100), flnm VARCHAR(100));
CREATE TABLE users (id SERIAL PRIMARY KEY, lgn VARCHAR(100), psw VARCHAR(100), token VARCHAR(300), wtdh INTEGER DEFAULT 0);
CREATE TABLE orders (id SERIAL PRIMARY KEY, nmb VARCHAR(100), sts VARCHAR(100), token VARCHAR(300), accural INTEGER DEFAULT -1, ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP, sumbals INTEGER DEFAULT 0);