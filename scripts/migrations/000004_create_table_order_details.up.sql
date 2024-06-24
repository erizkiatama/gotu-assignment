CREATE TABLE IF NOT EXISTS order_details (
  "id"              SERIAL          PRIMARY KEY,
  "order_id"        INT8            NOT NULL,
  "book_id"         INT8            NOT NULL,
  "quantity"        INT8            NOT NULL,
  "price"           INT8            NOT NULL,  
  "created_at"      TIMESTAMP(6)    NOT NULL DEFAULT (TIMEZONE('UTC', NOW())),
  "updated_at"      TIMESTAMP(6),
  "is_deleted"      BOOLEAN         NOT NULL DEFAULT false,
  CONSTRAINT fk_order_id
        FOREIGN KEY (order_id)
            REFERENCES orders(id)
            ON UPDATE CASCADE
            ON DELETE CASCADE,
    CONSTRAINT fk_book_id
        FOREIGN KEY (book_id)
            REFERENCES books(id)
            ON UPDATE CASCADE
            ON DELETE CASCADE
);
