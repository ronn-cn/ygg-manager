package main

// 配置文件结构体
type Configure struct {
	Ouid     string `toml:"Ouid"`
	Pin      string `toml:"Pin"`
	PriKey   string `toml:"PriKey"`
	PubKey   string `toml:"PubKey"`
	Port     string `toml:"Port"`
	DbHost   string `toml:"DbHost"`
	DbPort   string `toml:"DbPort"`
	DbUser   string `toml:"DbUser"`
	DbPasswd string `toml:"DbPasswd"`
	DbName   string `toml:"DbName"`
}

// gorm.Model 的定义
type Model struct {
	ID        int64 `json:"id" gorm:"column:id;autoIncrement;default:NULL"`
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime;NOT NULL"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

// 数据结构表
type Account struct {
	Model
	OUID           string      `json:"ouid" gorm:"column:ouid;primaryKey;NOT NULL"`                  // 统一标识
	Name           string      `json:"name" gorm:"column:name;"`                                     // 账号名称
	Account        string      `json:"account" gorm:"column:account;NOT NULL"`                       // 账号
	Passwd         string      `json:"passwd" gorm:"column:passwd;NOT NULL"`                         // 密码
	PermissionType *Permission `json:"type" gorm:"ForeignKey:TypeID;AssociationForeignKey:ID"`       // 账户权限类型
	TypeID         int         `json:"type_id" gorm:"column:type;NOT NULL"`                          // 账户权限类型ID
	Contact        string      `json:"contact" gorm:"column:contact"`                                // 联系方式
	Detail         string      `json:"detail" gorm:"column:detail"`                                  // 详细
	Status         int         `json:"status" gorm:"column:status;NOT NULL"`                         // 账户状态，启用1 禁用0
	Company        *Company    `json:"company" gorm:"ForeignKey:CompanyID;AssociationForeignKey:ID"` // 账户单位
	CompanyID      int         `json:"company_id" gorm:"column:company;default:NULL"`                // 账户的单位id
	IP             string      `json:"-" gorm:"-"`
}

func (Account) TableName() string {
	return "m_accounts"
}

// 返回账号接口类
type AccountApi struct {
	Model
	OUID           string      `json:"ouid" gorm:"column:ouid;primaryKey;NOT NULL"`                  // 统一标识
	Name           string      `json:"name" gorm:"column:name;"`                                     // 账号名称
	Account        string      `json:"account" gorm:"column:account;NOT NULL"`                       // 账号
	PermissionType *Permission `json:"type" gorm:"ForeignKey:TypeID;AssociationForeignKey:ID"`       // 账户权限类型
	TypeID         *int        `json:"type_id" gorm:"column:type;NOT NULL"`                          // 账户权限类型ID
	Contact        string      `json:"contact" gorm:"column:contact"`                                // 联系方式
	Detail         string      `json:"detail" gorm:"column:detail"`                                  // 详细
	Status         int         `json:"status" gorm:"column:status;NOT NULL"`                         // 账户状态，启用1 禁用0
	Company        *Company    `json:"company" gorm:"ForeignKey:CompanyID;AssociationForeignKey:ID"` // 账户单位
	CompanyID      *int        `json:"company_id" gorm:"column:company"`                             // 账户的单位id
	IP             string      `json:"-" gorm:"-"`
}

func (AccountApi) TableName() string {
	return "m_accounts"
}

type Permission struct {
	ID        int    `json:"id" gorm:"column:id;primaryKey;NOT NULL"` // 自增ID
	Name      string `json:"name" gorm:"column:name;NOT NULL"`        // 权限类名称
	AllowMenu string `json:"allow_menu" gorm:"column:allow_menu"`     // 允许使用的菜单列表，逗号,分割
}

func (Permission) TableName() string {
	return "m_permission"
}

type Device struct {
	Model
	OUID         string       `json:"ouid" gorm:"column:ouid;primaryKey;NOT NULL"`                      // 统一标识
	Name         string       `json:"name" gorm:"column:name"`                                          // 设备名称
	PIN          string       `json:"pin" gorm:"column:pin"`                                            // 设备PIN码
	System       *System      `json:"system" gorm:"ForeignKey:SystemOUID;AssociationForeignKey:OUID"`   // 设备系统
	SystemOUID   *string      `json:"system_ouid" gorm:"column:system"`                                 // 设备系统OUID
	DeviceStatus DeviceStatus `json:"status" gorm:"ForeignKey:StatusID;AssociationForeignKey:ID"`       // 设备系统
	StatusID     *int         `json:"status_id" gorm:"column:status"`                                   // 设备状态
	License      *License     `json:"license" gorm:"ForeignKey:LicenseCode;AssociationForeignKey:Code"` // 注册信息
	LicenseCode  *string      `json:"license_code" gorm:"column:license"`                               // 注册授权码
	ProductJson  string       `json:"product_json" gorm:"column:product_json"`                          // 产品信息
	InstallJson  string       `json:"install_json" gorm:"column:install_json"`                          // 安装信息
	SoldTime     *int64       `json:"sold_time" gorm:"column:sold_time"`                                // 出售时间
	InstalledAt  *int64       `json:"installed_at" gorm:"column:installed_at"`                          // 安装时间
	LastTime     *int64       `json:"last_time" gorm:"column:last_time"`                                // 最后登录时间
	OwnerCompany *Company     `json:"owner" gorm:"ForeignKey:OwnerID;AssociationForeignKey:ID"`         // 所有者单位
	OwnerID      *int64       `json:"owner_id" gorm:"column:owner"`                                     // 所有者ID
	Remark       string       `json:"remark" gorm:"column:remark"`                                      // 备注
	Manufacturer *int64       `json:"manufacturer" gorm:"column:manufacturer"`                          // 生产者
	Distributor  *int64       `json:"distributor" gorm:"column:distributor"`                            // 销售者
	Customer     *int64       `json:"customer" gorm:"column:customer"`                                  // 使用客户
}

func (Device) TableName() string {
	return "m_devices"
}

type System struct {
	Model
	OUID           string   `json:"ouid" gorm:"column:ouid;primaryKey;NOT NULL"`                    // 统一标识
	Name           string   `json:"name" gorm:"column:name"`                                        // 系统名称
	List           string   `json:"list" gorm:"column:list"`                                        // 关系
	Main           string   `json:"main" gorm:"column:main"`                                        // 主程序
	CreatorAccount *Account `json:"creator" gorm:"ForeignKey:CreatorID;AssociationForeignKey:OUID"` // 创建者
	CreatorID      *string  `json:"creator_id" gorm:"column:creator"`                               // 创建者ID
	Status         int      `json:"status" gorm:"column:status"`                                    // 状态
	Remark         string   `json:"remark" gorm:"column:remark"`                                    // 备注
}

func (System) TableName() string {
	return "m_system"
}

type License struct {
	Model
	Code           string   `json:"code" gorm:"column:code;primaryKey;NOT NULL"`                    // 授权码
	Count          int      `json:"count" gorm:"column:count"`                                      // 授权注册数量
	UseCount       int      `json:"use_count" gorm:"column:use_count"`                              // 已使用的数量
	Permit         string   `json:"permit" gorm:"column:permit"`                                    // 允许注册的系统
	ExpiresAt      int64    `json:"expires_at" gorm:"column:expires_at"`                            // 过期时间
	VendorJson     string   `json:"vendor_json" gorm:"column:vendor_json"`                          // 厂商信息
	CreatorAccount *Account `json:"creator" gorm:"ForeignKey:CreatorID;AssociationForeignKey:OUID"` // 创建者
	CreatorID      *string  `json:"creator_id" gorm:"column:creator"`                               // 创建者ID
	Remark         string   `json:"remark" gorm:"column:remark"`                                    // 备注
}

func (License) TableName() string {
	return "m_license"
}

type DeviceStatus struct {
	ID   int    `json:"id" gorm:"column:id;primaryKey;NOT NULL"` // 自增id
	Name string `json:"name" gorm:"column:name"`                 // 设备状态名称
}

func (DeviceStatus) TableName() string {
	return "m_device_status"
}

type Company struct {
	ID          int64  `json:"id" gorm:"column:id;primaryKey;NOT NULL"`            // 自增ID
	Name        string `json:"name" gorm:"column:name"`                            // 单位名称
	Type        string `json:"type" gorm:"column:type"`                            // 单位类型
	Address     string `json:"address" gorm:"column:address"`                      // 单位地址
	Telephone   string `json:"telephone" gorm:"column:telephone"`                  // 单位电话
	Email       string `json:"email" gorm:"column:email"`                          // 单位邮箱
	Website     string `json:"website" gorm:"column:website"`                      // 单位网站
	Description string `json:"description" gorm:"column:description"`              // 单位简介
	CreatedAt   int64  `json:"created_at" gorm:"column:created_at;autoCreateTime"` // 创建时间
	UpdatedAt   int64  `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"` // 更新时间
}

func (Company) TableName() string {
	return "m_company"
}

type Application struct {
	Model
	Appid          string   `json:"appid" gorm:"column:appid;primaryKey;NOT NULL"`                  // 应用ID
	Name           string   `json:"name" gorm:"column:name"`                                        // 应用名称
	Type           int      `json:"type" gorm:"column:type"`                                        // 应用类型
	Latest         string   `json:"latest" gorm:"column:latest"`                                    // 应用最新版本
	CreatorAccount *Account `json:"creator" gorm:"ForeignKey:CreatorID;AssociationForeignKey:OUID"` // 创建者
	CreatorID      *string  `json:"creator_id" gorm:"column:creator"`                               // 创建者ID
	Status         int      `json:"status" gorm:"column:status"`                                    // 应用的状态，正式应用，开发应用
	Remark         string   `json:"remark" gorm:"column:remark"`                                    // 备注信息
	Depend         string   `json:"depend" gorm:"column:depend"`                                    // 依赖于(应用)
}

func (Application) TableName() string {
	return "m_application"
}

type Version struct {
	ID             int64   `json:"id" gorm:"column:id;primaryKey;NOT NULL"`                        // 自增ID
	Version        string  `json:"version" gorm:"column:version"`                                  // 版本号
	Appid          string  `json:"appid" gorm:"column:appid"`                                      // 应用id
	Method         int     `json:"method" gorm:"column:method"`                                    // 更新方式
	Description    string  `json:"description" gorm:"column:description"`                          // 版本简介
	CreatedAt      int64   `json:"created_at" gorm:"column:created_at;autoCreateTime"`             // 创建时间
	From           string  `json:"from" gorm:"column:from"`                                        // 来自台主机上传
	CreatorAccount Account `json:"creator" gorm:"ForeignKey:CreatorID;AssociationForeignKey:OUID"` // 创建者
	CreatorID      string  `json:"creator_id" gorm:"column:creator"`                               // 创建者ID
}

func (Version) TableName() string {
	return "m_version"
}

type Setting struct {
	ID    int64  `json:"id" gorm:"column:id;primaryKey;NOT NULL"` // 自增ID
	Key   string `json:"key" gorm:"column:key;NOT NULL"`          // Key
	Value string `json:"value" gorm:"column:value"`               // Value
}

func (Setting) TableName() string {
	return "m_setting"
}

type Record struct {
	ID        int64  `json:"id" gorm:"column:id;primaryKey;NOT NULL"`            // 自增ID
	Type      int    `json:"type" gorm:"column:type"`                            // 日志类型
	Record    string `json:"record" gorm:"column:record"`                        // 日志记录
	CreatedAt int64  `json:"created_at" gorm:"column:created_at;autoCreateTime"` // 创建时间
}

func (Record) TableName() string {
	return "m_record"
}

type Menu struct {
	ID   int64  `json:"id" gorm:"column:id;primaryKey;NOT NULL"` // 自增ID
	Name string `json:"name" gorm:"column:name"`                 // 菜单名称
	Path string `json:"path" gorm:"column:path"`                 // 菜单路径
}

func (Menu) TableName() string {
	return "m_menu"
}
