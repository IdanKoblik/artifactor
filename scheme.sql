CREATE TABLE "product"(
    "id" BIGINT NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    PRIMARY KEY("id"),
    CONSTRAINT "product_name_unique" UNIQUE("name")
);

CREATE TABLE "users"(
    "username" VARCHAR(255) NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "mail" VARCHAR(255),
    "password" VARCHAR(255) NOT NULL,
    "admin" BOOLEAN NOT NULL,
    "upload" BOOLEAN NOT NULL,
    "delete" BOOLEAN NOT NULL,
    PRIMARY KEY("username")
);

CREATE TABLE "version"(
    "id" BIGINT NOT NULL,
    "product_id" BIGINT NOT NULL,
    "version" VARCHAR(255) NOT NULL,
    "file_path" VARCHAR(255) NOT NULL,
    "checksum" VARCHAR(255) NOT NULL,
    PRIMARY KEY("id"),
    CONSTRAINT "version_product_fk"
        FOREIGN KEY("product_id")
        REFERENCES "product"("id")
        ON DELETE CASCADE
);

CREATE TABLE "user_product"(
    "username" VARCHAR(255) NOT NULL,
    "product_id" BIGINT NOT NULL,
    PRIMARY KEY ("username", "product_id"),
    CONSTRAINT "user_product_username_fk"
        FOREIGN KEY("username")
        REFERENCES "users"("username")
        ON DELETE CASCADE,
    CONSTRAINT "user_product_product_fk"
        FOREIGN KEY("product_id")
        REFERENCES "product"("id")
        ON DELETE CASCADE
);

CREATE INDEX "idx_user_product_product_id"
ON "user_product"("product_id");
