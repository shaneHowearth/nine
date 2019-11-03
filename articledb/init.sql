CREATE USER articlewriter;
CREATE DATABASE article_postgres_db owner articlewriter;
GRANT ALL PRIVILEGES ON DATABASE article_postgres_db TO articlewriter;
\connect article_postgres_db articlewriter
CREATE TABLE IF NOT EXISTS article
(
	id SERIAL PRIMARY KEY,
	title TEXT,
	pub_date DATE,
	body TEXT,
	tags TEXT[]
);
ALTER USER articlewriter WITH PASSWORD 'hu8jmn3';
