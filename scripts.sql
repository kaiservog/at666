create table place (
  lat  NUMERIC(14, 11) NOT NULL,
  lon NUMERIC(14, 11) NOT NULL,
  name VARCHAR(20) NOT NULL
);

ALTER TABLE place ADD CONSTRAINT no_duplicate_coord UNIQUE (lat, lon);