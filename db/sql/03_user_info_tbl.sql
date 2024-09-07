--DDL を定義
\c tone-db

--ユーザ情報管理テーブル
CREATE TABLE mng_user_info_tbl(
  user_id   varchar(50)  NOT NULL primary key,
  user_name varchar(100) NOT NULL,
  email     varchar(150),
  gender    int,
  age       int,

  --共通項目
  delete_flag      int       NOT NULL,
  created_datetime timestamp NOT NULL,
  updated_datetime timestamp,
  delete_date      date 
);

INSERT INTO mng_user_info_tbl VALUES
( 'test_user', 'テストユーザ', NULL, NULL,NULL, 0, now(), NULL, NULL),
( 'sampleUser', 'サンプルユーザ', NULL, NULL,NULL, 0, now(), NULL, NULL),
( 'elf_hinmel', 'フリーレン', NULL, NULL,NULL, 0, now(), NULL, NULL),
( 'elf_falin', 'マルシル', NULL, NULL,NULL, 0, now(), NULL, NULL);