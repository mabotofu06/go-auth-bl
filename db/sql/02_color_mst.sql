--DDL を定義
\c go-auth-db

CREATE TABLE color_mst(
  color_code  varchar(6)  NOT NULL primary key,
  color_label varchar(20) NOT NULL,
  --共通項目
  delete_flag       int       NOT NULL,
  created_datetime  timestamp NOT NULL,
  updated_datetime  timestamp,
  delete_date       date
);

INSERT INTO color_mst VALUES
( '000000', 'BLACK' , 0, now(), NULL, NULL),
( 'ffffff', 'WHITE', 0, now(), NULL, NULL),
( 'ff0000', 'RED' , 0, now(), NULL, NULL),
( '00ff00', 'GREEN' , 0, now(), NULL, NULL),
( '0000ff', 'BLUE' , 0, now(), NULL, NULL);