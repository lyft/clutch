CREATE TABLE authn_tokens(
  user_id text,
  provider text,

  access_token bytea,
  refresh_token bytea,
  id_token bytea,
  expiry timestamptz,

  PRIMARY KEY(user_id, provider)
);
