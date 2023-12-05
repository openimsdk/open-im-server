module github.com/openimsdk/open-im-server/v3/tools/up35

go 1.19

require (
	github.com/go-sql-driver/mysql v1.7.1
	github.com/openimsdk/open-im-server/v3 v3.5.0
	github.com/openimsdk/open-im-server/v3/tools/data-conversion v0.0.0-00010101000000-000000000000
	go.mongodb.org/mongo-driver v1.12.1
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/driver/mysql v1.5.1
	gorm.io/gorm v1.25.4
)

require (
	github.com/OpenIMSDK/protocol v0.0.31 // indirect
	github.com/OpenIMSDK/tools v0.0.18 // indirect
	github.com/bwmarrin/snowflake v0.3.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/jinzhu/copier v0.4.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/image v0.13.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sync v0.4.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231012201019-e917dd12ba7a // indirect
	google.golang.org/grpc v1.59.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)

replace (
	github.com/openimsdk/open-im-server/v3 => ./../../../open-im-server
	github.com/openimsdk/open-im-server/v3/tools/data-conversion => ./../data-conversion
)
