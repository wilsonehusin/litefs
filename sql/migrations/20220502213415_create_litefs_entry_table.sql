-- migrate:up
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

-- migrate:down
DROP TABLE litefs_entry;
