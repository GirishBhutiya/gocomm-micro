CREATE TABLE "users" (
	"user_id" bigserial PRIMARY KEY,
	"first_name" VARCHAR(50) NOT NULL,
	"last_name" VARCHAR(50) NOT NULL,
	"email" VARCHAR(50) NOT NULL UNIQUE,
	"password" VARCHAR(100) NOT NULL,
	"status" VARCHAR(50) NOT NULL,
	"activation_code" VARCHAR(200) NOT NULL,
	"forgot_code" VARCHAR(200) NOT NULL,
	"forgot_code_senttime" TIMESTAMP,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL
) WITH (
  OIDS=FALSE
);



CREATE TABLE "customers" (
	"cust_id" bigserial PRIMARY KEY,
	"first_name" VARCHAR(50) NOT NULL,
	"last_name" VARCHAR(50) NOT NULL,
	"email" VARCHAR(50) NOT NULL UNIQUE,
	"password" VARCHAR(100) NOT NULL,
	"status" VARCHAR(50) NOT NULL,
	"activation_code" VARCHAR(200) NOT NULL,
	"forgot_code" VARCHAR(200) NOT NULL,
	"forgot_code_senttime" TIMESTAMP NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL
) WITH (
  OIDS=FALSE
);



CREATE TABLE "cust_addresses" (
	"cust_addr_id" bigserial PRIMARY KEY,
	"cust_id" bigserial NOT NULL,
	"street_address" VARCHAR(100) NOT NULL,
	"city" VARCHAR(100) NOT NULL,
	"state" VARCHAR(50) NOT NULL,
	"pincode" VARCHAR(100) NOT NULL,
	"phonenumber" VARCHAR(20) NOT NULL
) WITH (
  OIDS=FALSE
);



CREATE TABLE "product" (
	"product_id" bigserial PRIMARY KEY,
	"product_name" VARCHAR(300) NOT NULL,
	"product_slug" VARCHAR(350) NOT NULL UNIQUE,
	"cat_id" bigserial NOT NULL,
	"sub_cat_id" bigserial NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL
) WITH (
  OIDS=FALSE
);



CREATE TABLE "categories" (
	"cat_id" bigserial PRIMARY KEY,
	"cat_name" VARCHAR(100) NOT NULL,
	"cat_icon" VARCHAR(255) NOT NULL
) WITH (
  OIDS=FALSE
);



CREATE TABLE "sub_categories" (
	"sub_cat_id" bigserial PRIMARY KEY,
	"category_id" bigserial NOT NULL,
	"sub_cat_name" VARCHAR(100) NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL
) WITH (
  OIDS=FALSE
);



CREATE TABLE "variations" (
	"variation_id" bigserial PRIMARY KEY,
	"variation_name" VARCHAR(200) NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL
) WITH (
  OIDS=FALSE
);



CREATE TABLE "variation_options" (
	"var_opt_id" bigserial PRIMARY KEY,
	"option_name" VARCHAR(200) NOT NULL,
	"variation_id" bigserial NOT NULL
) WITH (
  OIDS=FALSE
);



CREATE TABLE "product_variations" (
	"id" bigserial PRIMARY KEY,
	"product_id" bigserial NOT NULL,
	"variation_id" bigserial NOT NULL
) WITH (
  OIDS=FALSE
);



CREATE TABLE "product_variation_options" (
	"id" bigserial PRIMARY KEY,
	"product_variaton_id" bigserial NOT NULL,
	"variationa_name" VARCHAR(50) NOT NULL,
	"variation_image" VARCHAR(200) NOT NULL,
	"sku" VARCHAR(50) NOT NULL,
	"price" FLOAT(50) NOT NULL,
	"product_stock_id" bigserial NOT NULL
) WITH (
  OIDS=FALSE
);



CREATE TABLE "products_stocks" (
	"id" bigserial PRIMARY KEY,
	"total_stock" bigserial NOT NULL,
	"unit_price" FLOAT NOT NULL,
	"total_price" FLOAT NOT NULL
) WITH (
  OIDS=FALSE
);


ALTER TABLE "cust_addresses" ADD CONSTRAINT "cust_addresses_fk0" FOREIGN KEY ("cust_id") REFERENCES "customers"("cust_id");

ALTER TABLE "product" ADD CONSTRAINT "product_fk0" FOREIGN KEY ("cat_id") REFERENCES "categories"("cat_id");
ALTER TABLE "product" ADD CONSTRAINT "product_fk1" FOREIGN KEY ("sub_cat_id") REFERENCES "sub_categories"("sub_cat_id");


ALTER TABLE "sub_categories" ADD CONSTRAINT "sub_categories_fk0" FOREIGN KEY ("category_id") REFERENCES "categories"("cat_id");


ALTER TABLE "variation_options" ADD CONSTRAINT "variation_options_fk0" FOREIGN KEY ("variation_id") REFERENCES "variations"("variation_id");

ALTER TABLE "product_variations" ADD CONSTRAINT "product_variations_fk0" FOREIGN KEY ("product_id") REFERENCES "product"("product_id");
ALTER TABLE "product_variations" ADD CONSTRAINT "product_variations_fk1" FOREIGN KEY ("variation_id") REFERENCES "variations"("variation_id");

ALTER TABLE "product_variation_options" ADD CONSTRAINT "product_variation_options_fk0" FOREIGN KEY ("product_variaton_id") REFERENCES "product_variations"("id");
ALTER TABLE "product_variation_options" ADD CONSTRAINT "product_variation_options_fk1" FOREIGN KEY ("product_stock_id") REFERENCES "products_stocks"("id");



INSERT INTO public.users (first_name,last_name,email,password,status,activation_code,forgot_code,forgot_code_senttime,created_at,updated_at) VALUES
	 ('Girish','Bhutiya','girish@bhutiya.com','$2a$12$77/o9XrlwNZJu7P5Pslp8u3zLB2eot8TKmAcxljxuxfcNtuFZ3Qne','1','','',NULL,'20223-01-02 00:00:00','2023-01-02 00:00:00');