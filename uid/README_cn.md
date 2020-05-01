# 全局自增序列

本模块提供全局自增id生成功能，提供以下几种方式，用户可以根据需求自行选择
* mysql
* redis
* 雪花算法

## mysql
基于mysql方式生成规则是利用数据库中的sequence表控制自增id，可以保证全局非自增；序列的生成严格依赖mysql，客户端使用的时候每次会获取一个
cache范围，cache使用完毕之后再去重新获取，所以只能保证全局自增无法保证连续。mysql方式可以提供int32和int64方式的序列，用户根据自己的需求
指定不同的表名称即可

sequence表结构如下，使用前请先在数据库中初始化
```sql
CREATE TABLE `sequence` (
  `id` int(11) NOT NULL,
  `next_id` bigint(20) DEFAULT NULL,
  `cache` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='uid sequence';
```

使用方式
```go
func main(){
    s, err := uid.NewUID("MySqlUID", uid.Config{
        Server: "127.0.0.1",
        Port:   3306, 
        User: "test_user_rw",
        Password:  "test_password",
        InitValue: 1,
        Step:      10,
        DBName:    "test",
        TableName: "int32_sequence",
    })
    if err != nil {
       panic(err)
    }
  
   if s.HasInt32() {
        fmt.Printf("Int32:%v\n", s.NextUID32())
   }   
}

```
## redis
redis生成序列的原理是固定一个key，使用IncBy方式做自增；原则上也是只能保证自增非连续，使用方式类似于mysql；同样redis方式可以支持
int32和int64的序列生成，需要分别使用不同的key创建不通的uid对象来处理。
```go
func main(){
 s, err := uid.NewUID("RedisUID", uid.Config{
    	Key: "uid", 
        Server: "127.0.0.1", 
        Password: "password",
        Port: 6379
 })
    
 if err != nil {
   panic(err)
 }

 // gen uid
 for i := 0; i < 10; i++ {
   if s.HasInt32() {
     fmt.Printf("Int32:%d\n", s.NextUID32())
   }
 }
}
```

## 雪花算法
雪花算法原理是在不通节点使用不同的node_id依托雪花算法生成的全局唯一id，雪花算法只能生成int64位的序列；无法提供int32的序列，
使用的时候需要留意，使用方式类似
```go
func main(){
 s, err := uid.NewUID("SnowflakeUID", uid.Config{
    NodeID: 1
 })

 if err != nil {
    panic(err)
 }

 // gen uid
 for i := 0; i < 10; i++ {
    fmt.Printf("Int64:%v\n", s.NextUID64())
 }
}
```

# 致谢
本模块是站在前辈的肩膀上完成的，把前辈的劳动成果整合实现的，特此表示感谢，如果有侵权或者使用不当的地方请及时告知，谢谢
* github.com/go-redis/redis
* github.com/go-sql-driver/mysql
* github.com/bwmarrin/snowflake
