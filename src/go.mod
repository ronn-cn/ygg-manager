module manager

go 1.16

replace (
	libs/convert => ../../go_library_v1/convert
	libs/logger => ../../go_library_v1/logger
	libs/network => ../../go_library_v1/network
	libs/ouid => ../../go_library_v1/ouid
)

require (
	github.com/gin-gonic/gin v1.7.7
	github.com/pelletier/go-toml v1.9.4
	gorm.io/driver/postgres v1.3.1
	gorm.io/gorm v1.23.3
	libs/convert v0.0.0-00010101000000-000000000000
	libs/logger v0.0.0-00010101000000-000000000000
	libs/network v0.0.0-00010101000000-000000000000
	libs/ouid v0.0.0-00010101000000-000000000000
)
