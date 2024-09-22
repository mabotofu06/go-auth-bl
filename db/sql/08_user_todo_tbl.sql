--DDL を定義
\c go-auth-db

--ユーザTodo管理テーブル
CREATE TABLE mng_user_todo_tbl(
  todo_id   varchar(50)  NOT NULL primary key,
  user_id   varchar(50)  NOT NULL,
  color_code varchar(6) NOT NULL default('000000'),
  favorite  int NOT NULL default(0),
  fix       int NOT NULL default(0),

  --共通項目
  delete_flag      int       NOT NULL,
  created_datetime timestamp NOT NULL,
  updated_datetime timestamp,
  delete_date      date,

  foreign key(todo_id)    references mng_todo_tbl     (todo_id),
--  foreign key(user_id)    references mng_user_info_tbl(user_id),
  foreign key(color_code) references color_mst        (color_code)
);