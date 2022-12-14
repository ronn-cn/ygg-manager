/*
 Navicat Premium Data Transfer

 Source Server         : 172.16.7.100_5432
 Source Server Type    : PostgreSQL
 Source Server Version : 100019
 Source Host           : 172.16.7.100:5432
 Source Catalog        : manager
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 100019
 File Encoding         : 65001

 Date: 28/03/2022 10:37:45
*/


-- ----------------------------
-- Sequence structure for m_accounts_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."m_accounts_id_seq";
CREATE SEQUENCE "public"."m_accounts_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for m_application_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."m_application_id_seq";
CREATE SEQUENCE "public"."m_application_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for m_assemble_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."m_assemble_id_seq";
CREATE SEQUENCE "public"."m_assemble_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for m_company_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."m_company_id_seq";
CREATE SEQUENCE "public"."m_company_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for m_device_status_ida_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."m_device_status_ida_seq";
CREATE SEQUENCE "public"."m_device_status_ida_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 32767
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for m_devices_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."m_devices_id_seq";
CREATE SEQUENCE "public"."m_devices_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for m_license_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."m_license_id_seq";
CREATE SEQUENCE "public"."m_license_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for m_menu_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."m_menu_id_seq";
CREATE SEQUENCE "public"."m_menu_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for m_permission_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."m_permission_id_seq";
CREATE SEQUENCE "public"."m_permission_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for m_record_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."m_record_id_seq";
CREATE SEQUENCE "public"."m_record_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for m_setting_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."m_setting_id_seq";
CREATE SEQUENCE "public"."m_setting_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 32767
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for m_version_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."m_version_id_seq";
CREATE SEQUENCE "public"."m_version_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Table structure for m_accounts
-- ----------------------------
DROP TABLE IF EXISTS "public"."m_accounts";
CREATE TABLE "public"."m_accounts" (
  "id" int8 NOT NULL DEFAULT nextval('m_accounts_id_seq'::regclass),
  "ouid" varchar(64) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(255) COLLATE "pg_catalog"."default",
  "type" int2 NOT NULL,
  "contact" varchar(255) COLLATE "pg_catalog"."default",
  "detail" text COLLATE "pg_catalog"."default",
  "status" int2,
  "create_at" int8 NOT NULL,
  "update_at" int8,
  "account" varchar(32) COLLATE "pg_catalog"."default" NOT NULL,
  "passwd" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;
COMMENT ON COLUMN "public"."m_accounts"."id" IS '??????id';
COMMENT ON COLUMN "public"."m_accounts"."ouid" IS '??????ouid';
COMMENT ON COLUMN "public"."m_accounts"."name" IS '????????????';
COMMENT ON COLUMN "public"."m_accounts"."type" IS '??????????????????';
COMMENT ON COLUMN "public"."m_accounts"."contact" IS '????????????';
COMMENT ON COLUMN "public"."m_accounts"."detail" IS '????????????';
COMMENT ON COLUMN "public"."m_accounts"."status" IS '?????????????????????1??????0';
COMMENT ON COLUMN "public"."m_accounts"."create_at" IS '????????????';
COMMENT ON COLUMN "public"."m_accounts"."update_at" IS '????????????';
COMMENT ON COLUMN "public"."m_accounts"."account" IS '??????';
COMMENT ON COLUMN "public"."m_accounts"."passwd" IS '??????';

-- ----------------------------
-- Table structure for m_application
-- ----------------------------
DROP TABLE IF EXISTS "public"."m_application";
CREATE TABLE "public"."m_application" (
  "id" int8 NOT NULL DEFAULT nextval('m_application_id_seq'::regclass),
  "appid" varchar(32) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(255) COLLATE "pg_catalog"."default",
  "type" int2,
  "latest" varchar(32) COLLATE "pg_catalog"."default",
  "creator" varchar(64) COLLATE "pg_catalog"."default",
  "create_at" int8,
  "update_at" int8
)
;
COMMENT ON COLUMN "public"."m_application"."id" IS '??????id';
COMMENT ON COLUMN "public"."m_application"."appid" IS '??????id';
COMMENT ON COLUMN "public"."m_application"."name" IS '????????????';
COMMENT ON COLUMN "public"."m_application"."type" IS '????????????';
COMMENT ON COLUMN "public"."m_application"."latest" IS '??????????????????';
COMMENT ON COLUMN "public"."m_application"."creator" IS '?????????';
COMMENT ON COLUMN "public"."m_application"."create_at" IS '????????????';
COMMENT ON COLUMN "public"."m_application"."update_at" IS '????????????';

-- ----------------------------
-- Table structure for m_assembly
-- ----------------------------
DROP TABLE IF EXISTS "public"."m_assembly";
CREATE TABLE "public"."m_assembly" (
  "id" int8 NOT NULL DEFAULT nextval('m_assemble_id_seq'::regclass),
  "ouid" varchar(64) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(255) COLLATE "pg_catalog"."default",
  "relation" json,
  "creator" varchar(64) COLLATE "pg_catalog"."default",
  "create_at" int8,
  "update_at" int8
)
;
COMMENT ON COLUMN "public"."m_assembly"."id" IS '??????id';
COMMENT ON COLUMN "public"."m_assembly"."ouid" IS '??????ouid';
COMMENT ON COLUMN "public"."m_assembly"."name" IS '????????????';
COMMENT ON COLUMN "public"."m_assembly"."relation" IS '??????';
COMMENT ON COLUMN "public"."m_assembly"."creator" IS '?????????';
COMMENT ON COLUMN "public"."m_assembly"."create_at" IS '????????????';
COMMENT ON COLUMN "public"."m_assembly"."update_at" IS '????????????';

-- ----------------------------
-- Table structure for m_company
-- ----------------------------
DROP TABLE IF EXISTS "public"."m_company";
CREATE TABLE "public"."m_company" (
  "id" int8 NOT NULL DEFAULT nextval('m_company_id_seq'::regclass),
  "name" varchar(255) COLLATE "pg_catalog"."default",
  "type" varchar(255) COLLATE "pg_catalog"."default",
  "address" varchar(255) COLLATE "pg_catalog"."default",
  "telephone" varchar(30) COLLATE "pg_catalog"."default",
  "email" varchar(50) COLLATE "pg_catalog"."default",
  "webside" varchar(255) COLLATE "pg_catalog"."default",
  "description" text COLLATE "pg_catalog"."default",
  "create_at" int8,
  "update_at" int8
)
;
COMMENT ON COLUMN "public"."m_company"."id" IS '??????id';
COMMENT ON COLUMN "public"."m_company"."name" IS '????????????';
COMMENT ON COLUMN "public"."m_company"."type" IS '????????????';
COMMENT ON COLUMN "public"."m_company"."address" IS '????????????';
COMMENT ON COLUMN "public"."m_company"."telephone" IS '????????????';
COMMENT ON COLUMN "public"."m_company"."email" IS '????????????';
COMMENT ON COLUMN "public"."m_company"."webside" IS '????????????';
COMMENT ON COLUMN "public"."m_company"."description" IS '????????????';
COMMENT ON COLUMN "public"."m_company"."create_at" IS '????????????';
COMMENT ON COLUMN "public"."m_company"."update_at" IS '????????????';

-- ----------------------------
-- Table structure for m_device_status
-- ----------------------------
DROP TABLE IF EXISTS "public"."m_device_status";
CREATE TABLE "public"."m_device_status" (
  "id" int2 NOT NULL DEFAULT nextval('m_device_status_ida_seq'::regclass),
  "name" varchar(16) COLLATE "pg_catalog"."default" NOT NULL
)
;
COMMENT ON COLUMN "public"."m_device_status"."id" IS '??????id';
COMMENT ON COLUMN "public"."m_device_status"."name" IS '??????????????????';

-- ----------------------------
-- Table structure for m_devices
-- ----------------------------
DROP TABLE IF EXISTS "public"."m_devices";
CREATE TABLE "public"."m_devices" (
  "id" int8 NOT NULL DEFAULT nextval('m_devices_id_seq'::regclass),
  "ouid" varchar(64) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(255) COLLATE "pg_catalog"."default",
  "pin" char(6) COLLATE "pg_catalog"."default" NOT NULL,
  "assembly" varchar(64) COLLATE "pg_catalog"."default",
  "status" int2 NOT NULL,
  "license" varchar(32) COLLATE "pg_catalog"."default",
  "product_json" text COLLATE "pg_catalog"."default",
  "install_json" text COLLATE "pg_catalog"."default",
  "create_at" int8 NOT NULL,
  "update_at" int8,
  "sell_time" int8,
  "install_at" int8,
  "last_time" int8,
  "owner" int8,
  "remark" varchar(255) COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."m_devices"."id" IS '??????id';
COMMENT ON COLUMN "public"."m_devices"."ouid" IS '??????ouid';
COMMENT ON COLUMN "public"."m_devices"."name" IS '????????????';
COMMENT ON COLUMN "public"."m_devices"."pin" IS '??????PIN???';
COMMENT ON COLUMN "public"."m_devices"."assembly" IS '????????????';
COMMENT ON COLUMN "public"."m_devices"."status" IS '????????????';
COMMENT ON COLUMN "public"."m_devices"."license" IS '?????????????????????';
COMMENT ON COLUMN "public"."m_devices"."product_json" IS '????????????';
COMMENT ON COLUMN "public"."m_devices"."install_json" IS '????????????';
COMMENT ON COLUMN "public"."m_devices"."create_at" IS '?????????????????????????????????';
COMMENT ON COLUMN "public"."m_devices"."update_at" IS '?????????';
COMMENT ON COLUMN "public"."m_devices"."sell_time" IS '????????????';
COMMENT ON COLUMN "public"."m_devices"."install_at" IS '????????????';
COMMENT ON COLUMN "public"."m_devices"."last_time" IS '??????????????????';
COMMENT ON COLUMN "public"."m_devices"."owner" IS '?????????';
COMMENT ON COLUMN "public"."m_devices"."remark" IS '???????????????255???????????????120???';

-- ----------------------------
-- Table structure for m_license
-- ----------------------------
DROP TABLE IF EXISTS "public"."m_license";
CREATE TABLE "public"."m_license" (
  "id" int8 NOT NULL DEFAULT nextval('m_license_id_seq'::regclass),
  "code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL,
  "count" int4,
  "permit" text COLLATE "pg_catalog"."default",
  "vendor_json" json,
  "expires_at" int8,
  "create_at" int8,
  "creator" varchar(64) COLLATE "pg_catalog"."default",
  "remarks" varchar(255) COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."m_license"."id" IS '??????id';
COMMENT ON COLUMN "public"."m_license"."code" IS '?????????';
COMMENT ON COLUMN "public"."m_license"."count" IS '????????????';
COMMENT ON COLUMN "public"."m_license"."permit" IS '?????????????????????';
COMMENT ON COLUMN "public"."m_license"."vendor_json" IS '????????????';
COMMENT ON COLUMN "public"."m_license"."expires_at" IS '????????????';
COMMENT ON COLUMN "public"."m_license"."create_at" IS '????????????';
COMMENT ON COLUMN "public"."m_license"."creator" IS '?????????';
COMMENT ON COLUMN "public"."m_license"."remarks" IS '????????????';

-- ----------------------------
-- Table structure for m_menu
-- ----------------------------
DROP TABLE IF EXISTS "public"."m_menu";
CREATE TABLE "public"."m_menu" (
  "id" int4 NOT NULL DEFAULT nextval('m_menu_id_seq'::regclass),
  "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "path" varchar(255) COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."m_menu"."id" IS '??????id';
COMMENT ON COLUMN "public"."m_menu"."name" IS '????????????';
COMMENT ON COLUMN "public"."m_menu"."path" IS '????????????';

-- ----------------------------
-- Table structure for m_permission
-- ----------------------------
DROP TABLE IF EXISTS "public"."m_permission";
CREATE TABLE "public"."m_permission" (
  "id" int2 NOT NULL DEFAULT nextval('m_permission_id_seq'::regclass),
  "name" varchar(32) COLLATE "pg_catalog"."default" NOT NULL,
  "allow_menu" varchar(255) COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."m_permission"."id" IS '??????id';
COMMENT ON COLUMN "public"."m_permission"."name" IS '???????????????';
COMMENT ON COLUMN "public"."m_permission"."allow_menu" IS '????????????????????????????????????,??????';

-- ----------------------------
-- Table structure for m_record
-- ----------------------------
DROP TABLE IF EXISTS "public"."m_record";
CREATE TABLE "public"."m_record" (
  "id" int8 NOT NULL DEFAULT nextval('m_record_id_seq'::regclass),
  "type" int2,
  "record" json,
  "create_at" int4
)
;
COMMENT ON COLUMN "public"."m_record"."id" IS '??????id';
COMMENT ON COLUMN "public"."m_record"."type" IS '????????????';
COMMENT ON COLUMN "public"."m_record"."record" IS '???????????????json??????';
COMMENT ON COLUMN "public"."m_record"."create_at" IS '????????????';

-- ----------------------------
-- Table structure for m_setting
-- ----------------------------
DROP TABLE IF EXISTS "public"."m_setting";
CREATE TABLE "public"."m_setting" (
  "id" int2 NOT NULL DEFAULT nextval('m_setting_id_seq'::regclass),
  "value" text COLLATE "pg_catalog"."default",
  "key" varchar(32) COLLATE "pg_catalog"."default" NOT NULL
)
;
COMMENT ON COLUMN "public"."m_setting"."id" IS '??????id';
COMMENT ON COLUMN "public"."m_setting"."value" IS 'value???';
COMMENT ON COLUMN "public"."m_setting"."key" IS 'key???';

-- ----------------------------
-- Table structure for m_version
-- ----------------------------
DROP TABLE IF EXISTS "public"."m_version";
CREATE TABLE "public"."m_version" (
  "id" int8 NOT NULL DEFAULT nextval('m_version_id_seq'::regclass),
  "version" varchar(32) COLLATE "pg_catalog"."default",
  "appid" varchar(32) COLLATE "pg_catalog"."default",
  "method" int2,
  "description" varchar(255) COLLATE "pg_catalog"."default",
  "create_at" int8,
  "from" varchar(255) COLLATE "pg_catalog"."default",
  "craetor" varchar(64) COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."m_version"."id" IS '??????id';
COMMENT ON COLUMN "public"."m_version"."version" IS '?????????';
COMMENT ON COLUMN "public"."m_version"."appid" IS '??????id';
COMMENT ON COLUMN "public"."m_version"."method" IS '????????????';
COMMENT ON COLUMN "public"."m_version"."description" IS '????????????';
COMMENT ON COLUMN "public"."m_version"."create_at" IS '????????????';
COMMENT ON COLUMN "public"."m_version"."from" IS '??????????????????';
COMMENT ON COLUMN "public"."m_version"."craetor" IS '?????????';

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."m_accounts_id_seq"
OWNED BY "public"."m_accounts"."id";
SELECT setval('"public"."m_accounts_id_seq"', 2, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."m_application_id_seq"
OWNED BY "public"."m_application"."id";
SELECT setval('"public"."m_application_id_seq"', 2, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."m_assemble_id_seq"
OWNED BY "public"."m_assembly"."id";
SELECT setval('"public"."m_assemble_id_seq"', 2, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."m_company_id_seq"
OWNED BY "public"."m_company"."id";
SELECT setval('"public"."m_company_id_seq"', 2, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."m_device_status_ida_seq"
OWNED BY "public"."m_device_status"."id";
SELECT setval('"public"."m_device_status_ida_seq"', 14, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."m_devices_id_seq"
OWNED BY "public"."m_devices"."id";
SELECT setval('"public"."m_devices_id_seq"', 2, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."m_license_id_seq"
OWNED BY "public"."m_license"."id";
SELECT setval('"public"."m_license_id_seq"', 2, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."m_menu_id_seq"
OWNED BY "public"."m_menu"."id";
SELECT setval('"public"."m_menu_id_seq"', 15, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."m_permission_id_seq"
OWNED BY "public"."m_permission"."id";
SELECT setval('"public"."m_permission_id_seq"', 8, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."m_record_id_seq"
OWNED BY "public"."m_record"."id";
SELECT setval('"public"."m_record_id_seq"', 2, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."m_setting_id_seq"
OWNED BY "public"."m_setting"."key";
SELECT setval('"public"."m_setting_id_seq"', 2, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."m_version_id_seq"
OWNED BY "public"."m_version"."id";
SELECT setval('"public"."m_version_id_seq"', 2, false);

-- ----------------------------
-- Primary Key structure for table m_accounts
-- ----------------------------
ALTER TABLE "public"."m_accounts" ADD CONSTRAINT "m_accounts_pkey" PRIMARY KEY ("ouid");

-- ----------------------------
-- Primary Key structure for table m_application
-- ----------------------------
ALTER TABLE "public"."m_application" ADD CONSTRAINT "m_application_pkey" PRIMARY KEY ("appid");

-- ----------------------------
-- Primary Key structure for table m_assembly
-- ----------------------------
ALTER TABLE "public"."m_assembly" ADD CONSTRAINT "m_assemble_pkey" PRIMARY KEY ("ouid");

-- ----------------------------
-- Primary Key structure for table m_company
-- ----------------------------
ALTER TABLE "public"."m_company" ADD CONSTRAINT "m_company_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table m_device_status
-- ----------------------------
ALTER TABLE "public"."m_device_status" ADD CONSTRAINT "m_device_status_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table m_devices
-- ----------------------------
ALTER TABLE "public"."m_devices" ADD CONSTRAINT "m_devices_pkey" PRIMARY KEY ("ouid");

-- ----------------------------
-- Primary Key structure for table m_license
-- ----------------------------
ALTER TABLE "public"."m_license" ADD CONSTRAINT "m_license_pkey" PRIMARY KEY ("code");

-- ----------------------------
-- Primary Key structure for table m_menu
-- ----------------------------
ALTER TABLE "public"."m_menu" ADD CONSTRAINT "m_menu_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table m_permission
-- ----------------------------
ALTER TABLE "public"."m_permission" ADD CONSTRAINT "m_permission_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table m_record
-- ----------------------------
ALTER TABLE "public"."m_record" ADD CONSTRAINT "m_record_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table m_setting
-- ----------------------------
ALTER TABLE "public"."m_setting" ADD CONSTRAINT "m_setting_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table m_version
-- ----------------------------
ALTER TABLE "public"."m_version" ADD CONSTRAINT "m_version_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Foreign Keys structure for table m_accounts
-- ----------------------------
ALTER TABLE "public"."m_accounts" ADD CONSTRAINT "type" FOREIGN KEY ("type") REFERENCES "public"."m_permission" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- ----------------------------
-- Foreign Keys structure for table m_application
-- ----------------------------
ALTER TABLE "public"."m_application" ADD CONSTRAINT "creator" FOREIGN KEY ("creator") REFERENCES "public"."m_accounts" ("ouid") ON DELETE CASCADE ON UPDATE CASCADE;

-- ----------------------------
-- Foreign Keys structure for table m_assembly
-- ----------------------------
ALTER TABLE "public"."m_assembly" ADD CONSTRAINT "creator" FOREIGN KEY ("creator") REFERENCES "public"."m_accounts" ("ouid") ON DELETE CASCADE ON UPDATE CASCADE;

-- ----------------------------
-- Foreign Keys structure for table m_devices
-- ----------------------------
ALTER TABLE "public"."m_devices" ADD CONSTRAINT "assembly" FOREIGN KEY ("assembly") REFERENCES "public"."m_assembly" ("ouid") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "public"."m_devices" ADD CONSTRAINT "license" FOREIGN KEY ("license") REFERENCES "public"."m_license" ("code") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "public"."m_devices" ADD CONSTRAINT "owner" FOREIGN KEY ("owner") REFERENCES "public"."m_company" ("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "public"."m_devices" ADD CONSTRAINT "status" FOREIGN KEY ("status") REFERENCES "public"."m_device_status" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- ----------------------------
-- Foreign Keys structure for table m_license
-- ----------------------------
ALTER TABLE "public"."m_license" ADD CONSTRAINT "creator" FOREIGN KEY ("creator") REFERENCES "public"."m_accounts" ("ouid") ON DELETE CASCADE ON UPDATE CASCADE;

-- ----------------------------
-- Foreign Keys structure for table m_version
-- ----------------------------
ALTER TABLE "public"."m_version" ADD CONSTRAINT "craetor" FOREIGN KEY ("craetor") REFERENCES "public"."m_accounts" ("ouid") ON DELETE CASCADE ON UPDATE CASCADE;
