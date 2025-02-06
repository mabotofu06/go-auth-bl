--DDL を定義
\c go-auth-db

--認可情報管理テーブル
CREATE TABLE mng_permission_inf_tbl(
  permission_id   varchar(10)  NOT NULL primary key,
  token           varchar(50),
  expire_datetime timestamp,

  --共通項目
  delete_flag      int       DEFAULT 0 NOT NULL,
  created_datetime timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_datetime timestamp,
  delete_date      date,
);

INSERT INTO mng_client_info_tbl VALUES
( 'CS10000000', 'API', 0, now(), NULL, NULL);