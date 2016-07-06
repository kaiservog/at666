create table comment (
  id integer PRIMARY KEY,
  lat NUMERIC(14, 11) NOT NULL,
  lon NUMERIC(14, 11) NOT NULL,
  comment_time timestamp NOT NULL,
  text VARCHAR(255) NOT NULL
);

CREATE SEQUENCE comment_id START 1000;
CREATE INDEX lat_lon ON comment (lat, lon);

--http://localhost:9002/at/comment/1.0000000000/1.0000000000/estou%20aqui
--http://localhost:9002/at/comment/1.0000000001/1.0000000001/segundo