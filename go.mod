module github.com/dbunion/com

go 1.13

require (
	github.com/RichardKnop/machinery v1.7.7
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b
	github.com/bwmarrin/snowflake v0.3.0
	github.com/fastly/go-utils v0.0.0-20180712184237-d95a45783239 // indirect
	github.com/fsnotify/fsnotify v1.4.9
	github.com/go-redis/redis/v7 v7.4.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/uuid v1.2.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/iancoleman/strcase v0.0.0-20191112232945-16388991a334
	github.com/jehiah/go-strftime v0.0.0-20171201141054-1d33003b3869 // indirect
	github.com/juju/errors v0.0.0-20220203013757-bd733f3c86b9
	github.com/juju/testing v0.0.0-20220203020004-a0ff61f03494 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.3.0+incompatible
	github.com/lestrrat-go/strftime v1.0.3 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/prometheus/client_golang v1.6.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.7.0
	github.com/tebeka/strftime v0.1.5 // indirect
	github.com/zssky/log v1.0.4
	github.com/zssky/tc v0.0.0-20200328060218-603c6a2939da
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	google.golang.org/grpc v1.35.0
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v0.18.2
	k8s.io/utils v0.0.0-20200414100711-2df71ebbae66 // indirect
	sigs.k8s.io/yaml v1.2.0
	vitess.io/vitess v0.0.0-20200524212726-2bbe82266007
)

replace github.com/RichardKnop/machinery => github.com/dbunion/machinery v0.0.0-20220514145235-db0f13beb54b
