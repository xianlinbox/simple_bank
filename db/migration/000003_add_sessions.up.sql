CREATE TABLE "sessions" (
  "id" uuid NOT NULL PRIMARY KEY,
  "username" varchar NOT NULL UNIQUE,
  "refresh_token" varchar NOT NULL, -- hashed password
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "expired_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users"("username");
