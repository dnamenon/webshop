

-- ----------------------------
--  Records of books
-- ----------------------------
BEGIN;
INSERT INTO items VALUES ('1', 'Brownie', 'July 21, 2015', 'JerBear', '10.00', 'images/brownie.jpg', 'this is a good brownie');
INSERT INTO items VALUES ('2', 'Pie', 'July 21, 2015', 'Swarley', '15.00', 'images/pie.jpg', 'good apple pie');
INSERT INTO items VALUES ('3', 'Cake', 'July 21, 2015', 'Roundabound', '12.00', 'images/bcake.jpg', 'decent enough cake');
INSERT INTO items VALUES ('4', 'Cheese', 'July 21, 2015', 'Crossfire', '3.00', 'images/cheese.jpg', 'serviceable cheese');
COMMIT;


-- ----------------------------
--  Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."books_id_seq" RESTART 4 OWNED BY "books"."id";

-- ----------------------------
--  Primary key structure for table books
-- ----------------------------
ALTER TABLE "public"."books" ADD PRIMARY KEY ("id") NOT DEFERRABLE INITIALLY IMMEDIATE;
