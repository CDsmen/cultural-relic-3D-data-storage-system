module artifact_svr

go 1.21.1

replace github.com/apache/thrift => github.com/apache/thrift v0.13.0

require (
	google.golang.org/grpc v1.63.0
	google.golang.org/protobuf v1.33.0
)

require (
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
)

require (
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240227224415-6ceb2ff114de // indirect
	gorm.io/driver/mysql v1.5.6
	gorm.io/gorm v1.25.9
)
