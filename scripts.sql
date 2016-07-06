create table comment (
  id integer PRIMARY KEY,
  place_lat NUMERIC(14, 11) NOT NULL,
  place_lon NUMERIC(14, 11) NOT NULL,
  comment_time timestamp NOT NULL,
  text VARCHAR(255) NOT NULL
);

CREATE SEQUENCE comment_id START 1000;