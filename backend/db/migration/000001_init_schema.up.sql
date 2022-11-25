CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "firstname" varchar NOT NULL,
  "lastname" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "mobile" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz
);

CREATE TABLE "accountBalance" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "balance" float NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz
);

CREATE TABLE "expenses" (
  "id" bigserial PRIMARY KEY,
  "amount" float NOT NULL,
  "description" varchar NOT NULL,
  "expenditure_date" timestamp NOT NULL,
  "category" varchar NOT NULL,
  "user_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz
);

CREATE TABLE "income" (
  "id" bigserial PRIMARY KEY,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz,
  "amount" float NOT NULL,
  "date" timestamptz NOT NULL,
  "user_id" bigint NOT NULL
);

CREATE TABLE "budget" (
  "id" bigserial PRIMARY KEY,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz,
  "budget_name" varchar NOT NULL,
  "initial_amount" float NOT NULL,
  "current_amount" float NOT NULL,
  "description" varchar NOT NULL,
  "start_date" timestamp NOT NULL,
  "end_date" timestamp NOT NULL,
  "user_id" bigint NOT NULL
);

CREATE TABLE "transactions" (
  "id" bigserial PRIMARY KEY,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz,
  "category" varchar NOT NULL,
  "amount" float NOT NULL,
  "transaction_date" timestamptz NOT NULL,
  "user_id" bigint NOT NULL
);

CREATE TABLE "savings" (
  "id" bigserial PRIMARY KEY,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz,
  "amount" float NOT NULL,
  "user_id" bigint NOT NULL
);

CREATE INDEX ON "users" ("firstname");

CREATE INDEX ON "users" ("lastname");

CREATE INDEX ON "users" ("email");

CREATE INDEX "fullname" ON "users" ("firstname", "lastname");

CREATE INDEX ON "expenses" ("category");

CREATE INDEX ON "expenses" ("user_id");

CREATE INDEX ON "expenses" ("description");

CREATE INDEX ON "expenses" ("expenditure_date");

CREATE INDEX ON "budget" ("budget_name");

CREATE INDEX ON "budget" ("user_id");

CREATE INDEX ON "transactions" ("category");

CREATE INDEX ON "transactions" ("user_id");

CREATE INDEX ON "transactions" ("amount");

COMMENT ON COLUMN "budget"."initial_amount" IS 'initial amount must not be greater than account balance';

ALTER TABLE "accountBalance" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "expenses" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "income" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "budget" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "savings" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
