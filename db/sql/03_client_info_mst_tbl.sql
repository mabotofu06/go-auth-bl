--DDL を定義
\c go-auth-db

--クライアント情報管理テーブル
CREATE TABLE mng_user_info_tbl(
  client_id   varchar(10)  NOT NULL primary key,
  client_name varchar(50) NOT NULL,

  --共通項目
  delete_flag      int       DEFAULT 0 NOT NULL,
  created_datetime timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_datetime timestamp,
  delete_date      date,

  -- 外部キー制約
  FOREIGN KEY (user_id) REFERENCES mng_user_auth_tbl(user_id)
);

INSERT INTO mng_user_info_tbl VALUES
( 'test_user', 'テストユーザ', NULL, NULL, NULL, 0, now(), NULL, NULL),
( 'sampleUser', 'サンプルユーザ', NULL, NULL, NULL, 0, now(), NULL, NULL),
( 'elf_frieren', 'フリーレン', NULL, NULL, NULL, 0, now(), NULL, NULL),
( 'elf_marsel', 'マルシル', NULL, NULL, NULL, 0, now(), NULL, NULL);