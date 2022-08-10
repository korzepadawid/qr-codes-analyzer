CREATE TABLE "users"
(
    "username"   varchar PRIMARY KEY NOT NULL,
    "email"      varchar UNIQUE      NOT NULL,
    "full_name"  varchar             NOT NULL,
    "password"   varchar             NOT NULL,
    "created_at" timestamptz         NOT NULL DEFAULT (now())
);

CREATE TABLE "groups"
(
    "id"          bigserial PRIMARY KEY NOT NULL,
    "owner"       varchar               NOT NULL,
    "title"       varchar               NOT NULL,
    "description" varchar               NOT NULL default '',
    "created_at"  timestamptz           NOT NULL DEFAULT (now())
);

CREATE TABLE "qr_codes"
(
    "id"              bigserial PRIMARY KEY NOT NULL,
    "owner"           varchar               NOT NULL,
    "group_id"        bigint                NOT NULL,
    "usages_count"    bigint                NOT NULL DEFAULT 0,
    "redirection_url" varchar               NOT NULL,
    "title"           varchar               NOT NULL,
    "description"     varchar               NOT NULL default '',
    "storage_url"     varchar UNIQUE        NOT NULL,
    "created_at"      timestamptz           NOT NULL DEFAULT (now())
);

CREATE TABLE "redirects"
(
    "id"         bigserial PRIMARY KEY NOT NULL,
    "qr_code_id" bigint                NOT NULL,
    "ipv4"       varchar               NOT NULL default '',
    "isp"        varchar               NOT NULL default '',
    "state"      varchar               NOT NULL default '',
    "country"    varchar               NOT NULL default '',
    "created_at" timestamptz           NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "groups" ("owner");

CREATE INDEX ON "qr_codes" ("group_id");

ALTER TABLE "groups"
    ADD FOREIGN KEY ("owner") REFERENCES "users" ("username") ON DELETE CASCADE;

ALTER TABLE "qr_codes"
    ADD FOREIGN KEY ("owner") REFERENCES "users" ("username") ON DELETE CASCADE;

ALTER TABLE "qr_codes"
    ADD FOREIGN KEY ("group_id") REFERENCES "groups" ("id") ON DELETE CASCADE;

ALTER TABLE "redirects"
    ADD FOREIGN KEY ("qr_code_id") REFERENCES "qr_codes" ("id") ON DELETE CASCADE;
