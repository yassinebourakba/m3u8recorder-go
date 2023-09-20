CREATE TABLE IF NOT EXISTS accounts(
   id serial NOT NULL,

   slug VARCHAR(50) NOT NULL,
   thumbnail_url VARCHAR,
   created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

   -- Constraints
   PRIMARY KEY (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS account_slug_index ON accounts (slug);