CREATE TABLE shortlink (
  -- hash is the unique identifer for the shortlink
  slhash text PRIMARY KEY,
  -- path is the page path to rediect to
  page_path text,
  -- state is shortlink state
  state JSONB
);

CREATE INDEX IF NOT EXISTS state_json ON shortlink USING GIN (state jsonb_path_ops);
