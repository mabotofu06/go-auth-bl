--DDL を定義
\c go-auth-db

--クライアント情報管理テーブル
CREATE TABLE mng_client_info_tbl(
  client_id   varchar(10)  NOT NULL primary key,
  client_name varchar(50) NOT NULL,

  --共通項目
  delete_flag      int       DEFAULT 0 NOT NULL,
  created_datetime timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_datetime timestamp,
  delete_date      date,
);

INSERT INTO mng_client_info_tbl VALUES
( 'CS10000000', 'API', 0, now(), NULL, NULL);