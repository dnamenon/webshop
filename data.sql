

-- ----------------------------
--  Records of books
-- ----------------------------
BEGIN;
INSERT INTO "users" VALUES ('1', 'JerBear', 'Garnee@example.com', 'hipster');
INSERT INTO "users" VALUES ('2', 'Swarley', 'Barney@example.com', 'happiness');
INSERT INTO "users" VALUES ('3', 'Roundabound', 'Anakin@example.com', 'instructor');
INSERT INTO "users" VALUES ('4', 'Crossfire', 'Freddie@example.com ', 'Crossfire');
COMMIT;


-- ----------------------------
--  Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."books_id_seq" RESTART 4 OWNED BY "books"."id";

-- ----------------------------
--  Primary key structure for table books
-- ----------------------------
ALTER TABLE "public"."books" ADD PRIMARY KEY ("id") NOT DEFERRABLE INITIALLY IMMEDIATE;
