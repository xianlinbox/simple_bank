CREATE TABLE "user" (
  "username" varchar NOT NULL PRIMARY KEY,
  "email" varchar NOT NULL UNIQUE,
  "password" varchar NOT NULL, -- hashed password
  "full_name" varchar NOT NULL,
  "password_expired_at" timestamptz NOT NULL DEFAULT (now() + interval '90 day'),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "account" ADD FOREIGN KEY ("owner") REFERENCES "user"("username");
