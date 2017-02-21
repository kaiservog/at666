create table comment (
  id integer PRIMARY KEY,
  lat NUMERIC(14, 11) NOT NULL,
  lon NUMERIC(14, 11) NOT NULL,
  comment_time timestamp NOT NULL,
  text VARCHAR(255) NOT NULL
);

CREATE SEQUENCE comment_id START 1000;
CREATE INDEX lat_lon ON comment (lat, lon);

--DB version 2

alter table comment add column nick VARCHAR(30) default 'anonymous';


--http://localhost:9002/at/comment/1.0000000000/1.0000000000/estou%20aqui
--http://localhost:9002/at/comment/1.0000000001/1.0000000001/segundo

#user=qdalnwznrczphp
#password=fTjPtcDNO-aEHXnI1xdsL0d_Xs
#host=ec2-54-243-249-154.compute-1.amazonaws.com
#database=d5vtbktupb1t70
#port=5432

--user=postgres
--password=admin
--host=localhost
--database=atdb
--port=5432
