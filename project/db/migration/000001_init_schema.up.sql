CREATE TABLE IF NOT EXISTS "users" (
  "id" bigserial PRIMARY KEY,
  "email" varchar UNIQUE NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "hashad_password" varchar NOT NULL,
  "active" bigint NOT NULL,
  "roll_id" bigint NOT NULL DEFAULT 3,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "password_changed_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "roll" (
  "id" bigserial PRIMARY KEY,
  "roll" varchar NOT NULL
);

ALTER TABLE "users" ADD FOREIGN KEY ("roll_id") REFERENCES "roll" ("id");

INSERT INTO "roll" ("id","roll") VALUES (1, 'admin');
INSERT INTO "roll" ("id","roll") VALUES (2, 'shop_manager');
INSERT INTO "roll" ("id","roll") VALUES (3, 'customer');