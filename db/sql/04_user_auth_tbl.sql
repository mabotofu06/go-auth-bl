--DDL を定義
\c go-auth-db

--ユーザー認証管理テーブル
CREATE TABLE mng_user_auth_tbl(
  user_id           varchar(50)  NOT NULL primary key,
  password          varchar(256) NOT NULL,
  session_token     varchar(256),
  last_session_time timestamp,

  --共通項目
  delete_flag      int       NOT NULL,
  created_datetime timestamp NOT NULL,
  updated_datetime timestamp,
  delete_date      date,

  foreign key(user_id) references mng_user_info_tbl(user_id)
);

INSERT INTO mng_user_auth_tbl VALUES
( 'test_user', 'password' , NULL, NULL, 0, now(), NULL, NULL),
( 'sampleUser', 'password', NULL, NULL, 0, now(), NULL, NULL),
( 'elf_hinmel', 'hinmel', NULL, NULL, 0, now(), NULL, NULL),
( 'elf_falin', '_falin', NULL, NULL, 0, now(), NULL, NULL);