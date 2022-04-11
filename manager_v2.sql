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

 Date: 24/03/2022 15:11:25
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
  "update_at" int8
)
;
COMMENT ON COLUMN "public"."m_accounts"."id" IS '自增id';
COMMENT ON COLUMN "public"."m_accounts"."ouid" IS '账户ouid';
COMMENT ON COLUMN "public"."m_accounts"."name" IS '账户名称';
COMMENT ON COLUMN "public"."m_accounts"."type" IS '账户权限类型';
COMMENT ON COLUMN "public"."m_accounts"."contact" IS '联系方式';
COMMENT ON COLUMN "public"."m_accounts"."detail" IS '详细信息';
COMMENT ON COLUMN "public"."m_accounts"."status" IS '账户状态，启用1禁用0';
COMMENT ON COLUMN "public"."m_accounts"."create_at" IS '创建时间';
COMMENT ON COLUMN "public"."m_accounts"."update_at" IS '更新时间';

-- ----------------------------
-- Records of m_accounts
-- ----------------------------

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
COMMENT ON COLUMN "public"."m_application"."id" IS '自增id';
COMMENT ON COLUMN "public"."m_application"."appid" IS '应用id';
COMMENT ON COLUMN "public"."m_application"."name" IS '应用名称';
COMMENT ON COLUMN "public"."m_application"."type" IS '应用类型';
COMMENT ON COLUMN "public"."m_application"."latest" IS '应用最新版本';
COMMENT ON COLUMN "public"."m_application"."creator" IS '创建者';
COMMENT ON COLUMN "public"."m_application"."create_at" IS '创建时间';
COMMENT ON COLUMN "public"."m_application"."update_at" IS '更新时间';

-- ----------------------------
-- Records of m_application
-- ----------------------------

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
COMMENT ON COLUMN "public"."m_assembly"."id" IS '自增id';
COMMENT ON COLUMN "public"."m_assembly"."ouid" IS '系统ouid';
COMMENT ON COLUMN "public"."m_assembly"."name" IS '系统名称';
COMMENT ON COLUMN "public"."m_assembly"."relation" IS '关系';
COMMENT ON COLUMN "public"."m_assembly"."creator" IS '创建者';
COMMENT ON COLUMN "public"."m_assembly"."create_at" IS '创建时间';
COMMENT ON COLUMN "public"."m_assembly"."update_at" IS '更新时间';

-- ----------------------------
-- Records of m_assembly
-- ----------------------------

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
COMMENT ON COLUMN "public"."m_company"."id" IS '自增id';
COMMENT ON COLUMN "public"."m_company"."name" IS '单位名称';
COMMENT ON COLUMN "public"."m_company"."type" IS '单位类型';
COMMENT ON COLUMN "public"."m_company"."address" IS '单位地址';
COMMENT ON COLUMN "public"."m_company"."telephone" IS '单位电话';
COMMENT ON COLUMN "public"."m_company"."email" IS '单位邮箱';
COMMENT ON COLUMN "public"."m_company"."webside" IS '单位网站';
COMMENT ON COLUMN "public"."m_company"."description" IS '单位简介';
COMMENT ON COLUMN "public"."m_company"."create_at" IS '创建时间';
COMMENT ON COLUMN "public"."m_company"."update_at" IS '更新时间';

-- ----------------------------
-- Records of m_company
-- ----------------------------

-- ----------------------------
-- Table structure for m_device_status
-- ----------------------------
DROP TABLE IF EXISTS "public"."m_device_status";
CREATE TABLE "public"."m_device_status" (
  "id" int2 NOT NULL DEFAULT nextval('m_device_status_ida_seq'::regclass),
  "name" varchar(16) COLLATE "pg_catalog"."default" NOT NULL
)
;
COMMENT ON COLUMN "public"."m_device_status"."id" IS '自增id';
COMMENT ON COLUMN "public"."m_device_status"."name" IS '设备状态名称';

-- ----------------------------
-- Records of m_device_status
-- ----------------------------
INSERT INTO "public"."m_device_status" VALUES (1, '注册');
INSERT INTO "public"."m_device_status" VALUES (2, '质检中');
INSERT INTO "public"."m_device_status" VALUES (3, '质检合格');
INSERT INTO "public"."m_device_status" VALUES (4, '质检不合格');
INSERT INTO "public"."m_device_status" VALUES (5, '销售中');
INSERT INTO "public"."m_device_status" VALUES (6, '已售出');
INSERT INTO "public"."m_device_status" VALUES (7, '待安装');
INSERT INTO "public"."m_device_status" VALUES (8, '安装中');
INSERT INTO "public"."m_device_status" VALUES (9, '安装完成');
INSERT INTO "public"."m_device_status" VALUES (10, '正常');
INSERT INTO "public"."m_device_status" VALUES (11, '停用');
INSERT INTO "public"."m_device_status" VALUES (12, '维护');
INSERT INTO "public"."m_device_status" VALUES (13, '开发');

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
COMMENT ON COLUMN "public"."m_devices"."id" IS '自增id';
COMMENT ON COLUMN "public"."m_devices"."ouid" IS '设备ouid';
COMMENT ON COLUMN "public"."m_devices"."name" IS '设备名称';
COMMENT ON COLUMN "public"."m_devices"."pin" IS '设备PIN码';
COMMENT ON COLUMN "public"."m_devices"."assembly" IS '设备系统';
COMMENT ON COLUMN "public"."m_devices"."status" IS '设备状态';
COMMENT ON COLUMN "public"."m_devices"."license" IS '来着哪个注册码';
COMMENT ON COLUMN "public"."m_devices"."product_json" IS '产品信息';
COMMENT ON COLUMN "public"."m_devices"."install_json" IS '安装信息';
COMMENT ON COLUMN "public"."m_devices"."create_at" IS '创建于（注册于）时间戳';
COMMENT ON COLUMN "public"."m_devices"."update_at" IS '更新于';
COMMENT ON COLUMN "public"."m_devices"."sell_time" IS '出售时间';
COMMENT ON COLUMN "public"."m_devices"."install_at" IS '安装时间';
COMMENT ON COLUMN "public"."m_devices"."last_time" IS '最后登录时间';
COMMENT ON COLUMN "public"."m_devices"."owner" IS '所有者';
COMMENT ON COLUMN "public"."m_devices"."remark" IS '备注，最多255字符，限制120字';

-- ----------------------------
-- Records of m_devices
-- ----------------------------

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
COMMENT ON COLUMN "public"."m_license"."id" IS '自增id';
COMMENT ON COLUMN "public"."m_license"."code" IS '授权码';
COMMENT ON COLUMN "public"."m_license"."count" IS '授权数量';
COMMENT ON COLUMN "public"."m_license"."permit" IS '允许注册的系统';
COMMENT ON COLUMN "public"."m_license"."vendor_json" IS '厂商信息';
COMMENT ON COLUMN "public"."m_license"."expires_at" IS '过期时间';
COMMENT ON COLUMN "public"."m_license"."create_at" IS '创建时间';
COMMENT ON COLUMN "public"."m_license"."creator" IS '创建者';
COMMENT ON COLUMN "public"."m_license"."remarks" IS '备注说明';

-- ----------------------------
-- Records of m_license
-- ----------------------------

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
COMMENT ON COLUMN "public"."m_menu"."id" IS '自增id';
COMMENT ON COLUMN "public"."m_menu"."name" IS '菜单名称';
COMMENT ON COLUMN "public"."m_menu"."path" IS '菜单路径';

-- ----------------------------
-- Records of m_menu
-- ----------------------------
INSERT INTO "public"."m_menu" VALUES (1, '首页', NULL);
INSERT INTO "public"."m_menu" VALUES (2, '账号管理', NULL);
INSERT INTO "public"."m_menu" VALUES (3, '设备管理', NULL);
INSERT INTO "public"."m_menu" VALUES (4, '系统管理', NULL);
INSERT INTO "public"."m_menu" VALUES (5, '应用管理', NULL);
INSERT INTO "public"."m_menu" VALUES (6, '密钥管理', NULL);
INSERT INTO "public"."m_menu" VALUES (7, '商户管理', NULL);
INSERT INTO "public"."m_menu" VALUES (8, '日志管理', NULL);
INSERT INTO "public"."m_menu" VALUES (9, '配置信息', NULL);

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
COMMENT ON COLUMN "public"."m_permission"."id" IS '自增id';
COMMENT ON COLUMN "public"."m_permission"."name" IS '权限类名称';
COMMENT ON COLUMN "public"."m_permission"."allow_menu" IS '允许使用的菜单列表，逗号,分割';

-- ----------------------------
-- Records of m_permission
-- ----------------------------
INSERT INTO "public"."m_permission" VALUES (1, '超级管理员', '0');
INSERT INTO "public"."m_permission" VALUES (2, '管理员', '0');
INSERT INTO "public"."m_permission" VALUES (4, '运维', '1,3,4,5,6,7,8');
INSERT INTO "public"."m_permission" VALUES (5, '生产', '1,3,4,5,6,8');
INSERT INTO "public"."m_permission" VALUES (3, '开发', '1,3,4,5,7,8');
INSERT INTO "public"."m_permission" VALUES (6, '销售', '1,3,4,5');
INSERT INTO "public"."m_permission" VALUES (7, '客户', '3,8');

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
COMMENT ON COLUMN "public"."m_record"."id" IS '自增id';
COMMENT ON COLUMN "public"."m_record"."type" IS '日志类型';
COMMENT ON COLUMN "public"."m_record"."record" IS '日志记录，json格式';
COMMENT ON COLUMN "public"."m_record"."create_at" IS '创建时间';

-- ----------------------------
-- Records of m_record
-- ----------------------------

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
COMMENT ON COLUMN "public"."m_setting"."id" IS '自增id';
COMMENT ON COLUMN "public"."m_setting"."value" IS 'value值';
COMMENT ON COLUMN "public"."m_setting"."key" IS 'key值';

-- ----------------------------
-- Records of m_setting
-- ----------------------------

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
COMMENT ON COLUMN "public"."m_version"."id" IS '自增id';
COMMENT ON COLUMN "public"."m_version"."version" IS '版本号';
COMMENT ON COLUMN "public"."m_version"."appid" IS '应用id';
COMMENT ON COLUMN "public"."m_version"."method" IS '更新方式';
COMMENT ON COLUMN "public"."m_version"."description" IS '版本描述';
COMMENT ON COLUMN "public"."m_version"."create_at" IS '创建时间';
COMMENT ON COLUMN "public"."m_version"."from" IS '来自哪台主机';
COMMENT ON COLUMN "public"."m_version"."craetor" IS '创建者';

-- ----------------------------
-- Records of m_version
-- ----------------------------

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
