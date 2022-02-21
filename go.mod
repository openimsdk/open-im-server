module Open_IM

go 1.15

require (
	github.com/Shopify/sarama v1.19.0
	github.com/Shopify/toxiproxy v2.1.4+incompatible // indirect
	github.com/alibabacloud-go/darabonba-openapi v0.1.11
	github.com/alibabacloud-go/dysmsapi-20170525/v2 v2.0.8
	github.com/alibabacloud-go/tea v1.1.17
	github.com/antonfisher/nested-logrus-formatter v1.3.0
	github.com/bwmarrin/snowflake v0.3.0
	github.com/eapache/go-resiliency v1.2.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/fatih/structs v1.1.0
	github.com/frankban/quicktest v1.14.0 // indirect
	github.com/garyburd/redigo v1.6.2
	github.com/gin-gonic/gin v1.7.0
	github.com/go-playground/validator/v10 v10.4.1
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.1.0
	github.com/golang/protobuf v1.5.2
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/jinzhu/copier v0.3.4
	github.com/jinzhu/gorm v1.9.16
	github.com/jinzhu/now v1.1.3 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.4 // indirect
	github.com/lib/pq v1.2.0 // indirect
	github.com/mattn/go-sqlite3 v1.14.6 // indirect
	github.com/minio/minio-go/v7 v7.0.22
	github.com/mitchellh/mapstructure v1.4.2
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/olivere/elastic/v7 v7.0.23
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pkg/errors v0.9.1
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/viper v1.9.0
	github.com/stretchr/testify v1.7.0
	github.com/tencentyun/qcloud-cos-sts-sdk v0.0.0-20210325043845-84a0811633ca
	github.com/tmc/grpc-websocket-proxy v0.0.0-20190109142713-0ad062ec5ee5 // indirect
	go.etcd.io/etcd v0.0.0-20200402134248-51bdeb39e698
	go.mongodb.org/mongo-driver v1.8.3
	golang.org/x/image v0.0.0-20210220032944-ac19c3e999fb
	golang.org/x/net v0.0.0-20210917221730-978cfadd31cf
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.29.1
