CREATE TABLE IF NOT EXISTS orders (
  "id"              SERIAL          PRIMARY KEY,
  "user_id"         INT8            NOT NULL,
  "total_quantity"  INT8            NOT NULL,
  "total_price"     INT8            NOT NULL,  
  "created_at"      TIMESTAMP(6)    NOT NULL DEFAULT (TIMEZONE('UTC', NOW())),
  "updated_at"      TIMESTAMP(6),
  "is_deleted"      BOOLEAN         NOT NULL DEFAULT false,
  CONSTRAINT fk_user_id
        FOREIGN KEY (user_id)
            REFERENCES users(id)
            ON UPDATE CASCADE
            ON DELETE CASCADE
);
