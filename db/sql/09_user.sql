-- DDL を定義

-- ユーザTodo管理テーブル
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  age INT,

  -- 共通項目
  delete_flag INT NOT NULL,
  created_datetime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_datetime TIMESTAMP,
  delete_date DATE
);

INSERT INTO users (id, name, age, delete_flag, created_datetime, updated_datetime, delete_date) VALUES
(1, 'サンプルユーザ', 20, 0, now(), NULL, NULL),
(2, 'テストユーザ', 30, 0, now(), NULL, NULL),
(3, 'フリーレン', 1000, 0, now(), NULL, NULL),
(4, 'マルシル', 2000, 0, now(), NULL, NULL);