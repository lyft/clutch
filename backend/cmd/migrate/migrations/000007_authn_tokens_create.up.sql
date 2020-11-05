CREATE TYPE token_type AS ENUM ('oidc');

CREATE TABLE authn_tokens(
  id text,
  provider text,
  type token_type,
  refresh_token bytea,

  PRIMARY KEY(id, provider)
);
