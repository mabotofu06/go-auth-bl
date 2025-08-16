--DDL を定義
\c go-auth-db

--認証元クライアントマスタ
CREATE TABLE mng_client_mst(
  client_id     char(10)  NOT NULL primary key,
  client_name   varchar(50)  NOT NULL,
  client_host   varchar(255) NOT NULL,
  --共通項目
  delete_flag      int       NOT NULL DEFAULT 0,
  created_datetime timestamp NOT NULL DEFAULT current_timestamp,
  updated_datetime timestamp,
  delete_date      date
);

--インデックスを定義
CREATE INDEX idx_mng_client_mst_client_id ON mng_client_mst(client_id);
CREATE INDEX idx_mng_client_mst_client_host ON mng_client_mst(client_host);

INSERT INTO
  mng_client_mst(client_id, client_name, client_host)
VALUES
('MYS0000000', 'Global Auth(Local)', 'localhost'),
('MYS0000001', 'System One(Local)' , 'localhost:8080'),
('MYS0000002', 'System Two(Local)' , 'localhost:3000');
