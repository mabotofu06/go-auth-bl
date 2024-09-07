--DDL を定義
\c tone-db

--Todo管理テーブル
CREATE TABLE mng_todo_tbl(
  todo_id          varchar(50)  NOT NULL primary key,
  title            varchar(50) NOT NULL,
  note             text,
  todo_status_code int       NOT NULL,
  --共通項目
  delete_flag      int       NOT NULL,
  created_datetime timestamp NOT NULL,
  updated_datetime timestamp,
  delete_date      date,

  foreign key(todo_status_code) references todo_status_mst(todo_status_code)
);

INSERT INTO mng_todo_tbl VALUES 
('todo_sample_id_00000000000000', 'Todoタイトル１', 'てすとてすと', 0, 0, now(), NULL, NULL),
('todo_sample_id_00000000000001', 'Todoタイトル２', 'てすとてすと', 2, 0, now(), NULL, NULL),
('todo_sample_id_00000000000002', 'Todoタイトル３', 'てすとてすと', 1, 0, now(), NULL, NULL),
('todo_sample_id_00000000000003', 'Todoタイトル４', 'てすとてすと', 1, 0, now(), NULL, NULL),
('todo_sample_id_00000000000004', 'Todoタイトル５', 'てすとてすと', 2, 0, now(), NULL, NULL),
('todo_sample_id_00000000000005', 'Todoタイトル６', 'てすとてすと', 0, 0, now(), NULL, NULL);