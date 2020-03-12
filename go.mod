module HommyUpdate

go 1.14

replace GrpcCommon => ../GrpcCommon

require (
	GrpcCommon v0.0.0-00010101000000-000000000000
	github.com/garyburd/redigo v1.6.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/jinzhu/gorm v1.9.12
	golang.org/x/net v0.0.0-20190404232315-eb5bcb51f2a3
	google.golang.org/grpc v1.28.0
	gopkg.in/yaml.v2 v2.2.8
)
