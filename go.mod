module github.com/dbunion/com

go 1.13

require (
	github.com/RichardKnop/machinery v1.7.7
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b
	github.com/bwmarrin/snowflake v0.3.0
	github.com/fsnotify/fsnotify v1.4.7
	github.com/go-redis/redis/v7 v7.2.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/uuid v1.1.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/iancoleman/strcase v0.0.0-20191112232945-16388991a334
	github.com/jinzhu/gorm v1.9.12 // indirect
	github.com/juju/errors v0.0.0-20200330140219-3fe23663418f
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/prometheus/client_golang v1.6.0 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.4.0
	github.com/youtube/vitess v2.1.1+incompatible // indirect
	github.com/zssky/log v1.0.4
	github.com/zssky/tc v0.0.0-20200328060218-603c6a2939da
	gitlab.com/opennota/check v0.0.0-20181224073239-ccaba434e62a // indirect
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a
	google.golang.org/grpc v1.29.1
	honnef.co/go/tools v0.0.1-2020.1.4 // indirect
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v0.18.2
	k8s.io/utils v0.0.0-20200414100711-2df71ebbae66 // indirect
	sigs.k8s.io/yaml v1.2.0
	vitess.io/vitess v0.0.0-00010101000000-000000000000
)

replace vitess.io/vitess => vitess.io/vitess v0.0.0-20181209180904-4f192d1003d1
