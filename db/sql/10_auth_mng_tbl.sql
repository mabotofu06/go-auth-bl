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

--パスワードは"password"をbcryptでハッシュ化したものを格納
INSERT INTO mng_user_auth_tbl VALUES
( 'test_user', '$2a$10$olHJMPAHQc63r79xZJQrO.wfuxcChHMCIh8rSjLgIniZ/Npfyw9aG' , NULL, NULL, 0, now(), NULL, NULL),
( 'sampleUser', '$2a$10$z7GBWNWVVtbPjNmtR3uvA.HQ06mHw/.G/sFbSzuLxTfs7PfI8CyL2', NULL, NULL, 0, now(), NULL, NULL),
( 'elf_hinmel', '$2a$10$Hm77dQ1nnoPIZIftKY2SNu7Ae60yf4HMKr8d3Y5zEz.86Mg8WPKFO', NULL, NULL, 0, now(), NULL, NULL),
( 'elf_falin', '$2a$10$wrw0DEEiHjo5y1X3Bc0Dgeh/kSaASzMB0kXwhfddra5B4D3C2n8SW', NULL, NULL, 0, now(), NULL, NULL);