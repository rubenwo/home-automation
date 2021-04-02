CREATE USER user IF NOT EXISTS;

CREATE DATABASE inventory_database;
GRANT ALL PRIVILEGES ON DATABASE inventory_database TO user;

CREATE DATABASE registry_database;
GRANT ALL PRIVILEGES ON DATABASE registry_database TO user;

CREATE DATABASE food_database;
GRANT ALL PRIVILEGES ON DATABASE food_database TO user;

CREATE DATABASE camera_database;
GRANT ALL PRIVILEGES ON DATABASE camera_database TO user;