# Sequence

[中文版](https://github.com/dbunion/com/blob/master/uid/README_cn.md) 
[English](https://github.com/dbunion/com/blob/master/uid/README.md) 

This module provides the function of generating sequence, and provides the following ways. Users can choose according to their needs
* mysql
* redis
* snowflake

## mysql
The generation rule based on MySQL is to use the sequence table in the database to control the auto increment ID, which can ensure the global non auto increment; the generation of sequence strictly depends on MySQL, and the client will get one at a time when using it.
The cache range is retrieved after the cache is used, so it can only guarantee that the global self increment cannot guarantee continuity. 
MySQL mode can provide int32 and Int64 mode sequences, and users can meet their own needs. Specify a different table name. If the provided user has the permission to create a table and modify the parameter autocreatetable = true, the program will automatically create the required sequence table structure

The sequence table structure is as follows. Please initialize it in the database before use
```sql
CREATE TABLE `sequence` (
  `id` int(11) NOT NULL,
  `next_id` bigint(20) DEFAULT NULL,
  `cache` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='uid sequence';
```

How to use
```go
package main
import (
    "fmt"
    "github.com/dbunion/com/uid"
    "time"
)

func main(){
    s, err := uid.NewUID(uid.TypeMySQL, uid.Config{
        Server:          "127.0.0.1",
        Port:            3306,
        User:            "test",
        Password:        "123456",
        InitValue:       time.Now().Unix(),
        Step:            10,
        DBName:          "test",
        TableName:       "int32_seq",
        AutoCreateTable: true,
    })

    if err != nil {
        fmt.Printf("start err:%v\n", err)
        return
    }

   if s.HasInt32() {
        fmt.Printf("Int32:%v\n", s.NextUID32())
   }   
}

```
## redis
The principle of generating sequence in redis is to fix a key and use the incby method to do auto increment. In principle,  it can only ensure that the auto increment is discontinuous and the usage method is similar to MySQL. The same can be supported in redis method. The sequence generation of int32 and Int64 needs to use different keys to create different uid objects for processing.

```go
package main
import (
    "fmt"
    "github.com/dbunion/com/uid"
    "time"
)

func main(){
    s, err := uid.NewUID(uid.TypeRedis, uid.Config{
    	Key: "uid", 
        Server: "127.0.0.1", 
        Password: "password",
        Port: 6379,
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

## SnowFlake
The principle of snowflake algorithm is to use different node ID in different nodes to rely on the globally unique ID generated by snowflake algorithm. Snowflake algorithm can only generate Int64 bit sequence; it cannot provide int32 sequence, Please pay attention to it when using it. The way of using it is similar.

```go
package main
import (
    "fmt"
    "github.com/dbunion/com/uid"
    "time"
)

func main(){
    s, err := uid.NewUID(uid.TypeSnowFlake, uid.Config{
        NodeID: 1,
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

# Thanks
This module is completed on the shoulders of the predecessors, integrating the achievements of the predecessors. If there is any infringement or improper use, please inform me in time. Thank you.
* github.com/go-redis/redis
* github.com/go-sql-driver/mysql
* github.com/bwmarrin/snowflake
