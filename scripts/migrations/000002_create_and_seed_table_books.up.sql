CREATE TABLE IF NOT EXISTS books (
  "id"            SERIAL          PRIMARY KEY,
  "title"         VARCHAR(255)    NOT NULL,
  "author"        VARCHAR(255)    NOT NULL,
  "description"   TEXT,
  "price"         INT8            NOT NULL,  
  "created_at"    TIMESTAMP(6)    NOT NULL DEFAULT (TIMEZONE('UTC', NOW())),
  "updated_at"    TIMESTAMP(6),
  "is_deleted"    BOOLEAN         NOT NULL DEFAULT false 
);


  INSERT INTO books ("title", "author", "description", "price")
    VALUES
      ('The Great Gatsby', 'F. Scott Fitzgerald', 'Description for The Great Gatsby', 100000),
      ('It Ends with Us', 'Colleen Hoover', 'Description for It Ends with Us', 150000),
      ('To Kill a Mockingbird', 'Harper Lee', 'Description for To Kill a Mockingbird', 200000),
      ('1984', 'George Orwell', 'Description for 1984', 225000),
      ('The Catcher in the Rye', 'J.D. Salinger', 'Description for The Catcher in the Rye', 250000),
      ('Pride and Prejudice', 'Jane Austen', 'Description for Pride and Prejudice', 300000),
      ('The Hobbit', 'J.R.R. Tolkien', 'Description for The Hobbit', 375000),
      ('Moby-Dick', 'Herman Melville', 'Description for Moby-Dick', 400000),
      ('War and Peace', 'Leo Tolstoy', 'Description for War and Peace', 450000),
      ('Ulysses', 'James Joyce', 'Description for Ulysses', 500000);
