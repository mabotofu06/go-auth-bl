--DDL を定義
\c go-auth-db

--マスタ
CREATE TABLE todo_status_mst(
  todo_status_code  int  NOT NULL primary key,
  todo_status_label varchar(20) NOT NULL,
  --共通項目
  delete_flag       int       NOT NULL,
  created_datetime  timestamp NOT NULL,
  updated_datetime  timestamp,
  delete_date       date
);

INSERT INTO todo_status_mst VALUES
( 0, 'TODO' , 0, now(), NULL, NULL),
( 1, 'DOING', 0, now(), NULL, NULL),
( 2, 'DONE' , 0, now(), NULL, NULL);