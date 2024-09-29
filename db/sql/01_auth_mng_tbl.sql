--DDL を定義
\c go-auth-db

--ユーザー認証管理テーブル
CREATE TABLE mng_user_auth_tbl(
  user_id           varchar(50) NOT NULL primary key,
  password          varchar(60) NOT NULL,
  password_histry1  varchar(60),
  password_histry2  varchar(60),
  password_histry3  varchar(60),
  password_fail_cnt int NOT NULL DEFAULT 0,
  password_lock     int NOT NULL DEFAULT 0,

--共通項目
  delete_flag      int       NOT NULL DEFAULT 0,
  created_datetime timestamp NOT NULL DEFAULT now(),
  updated_datetime timestamp,
  delete_date      date
);

--パスワードは"password"をbcryptでハッシュ化したものを格納
INSERT INTO mng_user_auth_tbl (user_id, password, password_histry1, password_histry2, password_histry3, password_fail_cnt, password_lock, delete_flag, created_datetime, updated_datetime, delete_date) VALUES
( 'test_user', '$2a$10$olHJMPAHQc63r79xZJQrO.wfuxcChHMCIh8rSjLgIniZ/Npfyw9aG', NULL, NULL, NULL, 0, 0, 0, now(), NULL, NULL),
( 'sampleUser', '$2a$10$z7GBWNWVVtbPjNmtR3uvA.HQ06mHw/.G/sFbSzuLxTfs7PfI8CyL2', NULL, NULL, NULL, 0, 0, 0, now(), NULL, NULL),
( 'elf_frieren', '$2a$10$Hm77dQ1nnoPIZIftKY2SNu7Ae60yf4HMKr8d3Y5zEz.86Mg8WPKFO', NULL, NULL, NULL, 0, 0, 0, now(), NULL, NULL),
( 'elf_marsel', '$2a$10$wrw0DEEiHjo5y1X3Bc0Dgeh/kSaASzMB0kXwhfddra5B4D3C2n8SW', NULL, NULL, NULL, 0, 0, 0, now(), NULL, NULL);