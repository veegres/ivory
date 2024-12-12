#!/bin/sh

psql -d "$1" -c "
CREATE EXTENSION pgstattuple;

CREATE ROLE trust WITH LOGIN SUPERUSER;
CREATE ROLE admin WITH LOGIN SUPERUSER PASSWORD 'admin';
CREATE ROLE sslca WITH LOGIN SUPERUSER PASSWORD 'sslca';
CREATE ROLE sslfull WITH LOGIN SUPERUSER PASSWORD 'sslfull';

CREATE TABLE users (
    id int CONSTRAINT users_id PRIMARY KEY,
    firstName varchar (25),
    lastName varchar (25)
);

INSERT INTO users(id, firstname, lastname)
SELECT id, id::text, id::text
FROM generate_series(1,1000) id;

CREATE TABLE films (
    code        char(5) CONSTRAINT films_code PRIMARY KEY,
    title       varchar(40) NOT NULL,
    did         integer NOT NULL,
    date_prod   date,
    kind        varchar(10),
    len         interval hour to minute
);

INSERT INTO films(code, title, did)
SELECT id::text, id::text, id
FROM generate_series(1,10000) id;
"
