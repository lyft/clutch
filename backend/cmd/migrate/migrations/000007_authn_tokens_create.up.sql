CREATE TABLE authn_tokens(
  user_id text,
  provider text,

  token_type text,
  id_token bytea,
  access_token bytea,
  refresh_token bytea,

  PRIMARY KEY(user_id, provider)
);
