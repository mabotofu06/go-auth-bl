--DDL を定義
\c go-auth-db

--ユーザー認証管理テーブル
CREATE TABLE mng_user_auth_tbl(
  id           varchar(256)  NOT NULL primary key,--認証コードとして使用
  token        varchar(256) NOT NULL,
  expire_datetime timestamp,
  
  --共通項目
  delete_flag      int       NOT NULL,
  created_datetime timestamp NOT NULL,
  updated_datetime timestamp,
  delete_date      date

--  foreign key(user_id) references mng_user_info_tbl(user_id)
);

--INSERT INTO mng_user_auth_tbl VALUES
--( 'test_user', 'password' , NULL, NULL, 0, now(), NULL, NULL),
--( 'sampleUser', 'password', NULL, NULL, 0, now(), NULL, NULL),
--( 'elf_hinmel', 'hinmel', NULL, NULL, 0, now(), NULL, NULL),
--( 'elf_falin', '_falin', NULL, NULL, 0, now(), NULL, NULL);