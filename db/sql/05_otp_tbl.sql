--DDL を定義
\c go-auth-db

--OTP認証管理テーブル
CREATE TABLE mng_otp_tbl(
  user_id        varchar(50)  NOT NULL primary key,
  otp            varchar(7),
  otp_valid_time timestamp,
  --共通項目
  delete_flag      int       NOT NULL,
  created_datetime timestamp NOT NULL,
  updated_datetime timestamp,
  delete_date      date,

  foreign key(user_id) references mng_user_auth_tbl(user_id)
);