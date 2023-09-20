CREATE TABLE IF NOT EXISTS records(
   id serial NOT NULL,
   account_id serial NOT NULL,
   url VARCHAR NOT NULL,
   is_hidden BOOLEAN DEFAULT false,

   name VARCHAR(100),
   duration INT,
   thumbnail_url VARCHAR,
   views INT DEFAULT 0,
   likes INT DEFAULT 0,
   published_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
   res_width INT,
   res_height INT,

   source VARCHAR(50),

   -- Constraints
   PRIMARY KEY (id),
   FOREIGN KEY (account_id) REFERENCES accounts(id)
);