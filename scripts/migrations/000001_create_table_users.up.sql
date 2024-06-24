CREATE TABLE IF NOT EXISTS users (
  "id"            SERIAL          PRIMARY KEY,
  "email"         VARCHAR(255)    NOT NULL UNIQUE,
  "password"      TEXT            NOT NULL,
  "name"          VARCHAR(255)    NOT NULL,
  "created_at"    TIMESTAMP(6)    NOT NULL DEFAULT (TIMEZONE('UTC', NOW())),
  "updated_at"    TIMESTAMP(6),
  "is_deleted"    BOOLEAN         NOT NULL DEFAULT false 
);
