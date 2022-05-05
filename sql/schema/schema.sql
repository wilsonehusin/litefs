CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(255) primary key);
CREATE TABLE litefs_entry (
  id TEXT NOT NULL,
  parent_id TEXT,
  name TEXT NOT NULL,
  modtime TEXT NOT NULL,
  content BLOB,
  PRIMARY KEY (id),
  FOREIGN KEY (parent_id) references litefs_entry(id) ON DELETE CASCADE,
  UNIQUE (parent_id, name)
);
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20220502213415');
