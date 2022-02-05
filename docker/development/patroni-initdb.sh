#!/bin/sh

psql -d "$1" -c "
CREATE EXTENSION pgstattuple;

CREATE TABLE users (
    id int CONSTRAINT users_id PRIMARY KEY,
    firstName varchar (25),
    lastName varchar (25)
);

CREATE TABLE films (
    code        char(5) CONSTRAINT films_code PRIMARY KEY,
    title       varchar(40) NOT NULL,
    did         integer NOT NULL,
    date_prod   date,
    kind        varchar(10),
    len         interval hour to minute
);
"
