package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dbunion/com/uid"
	// import mysql driver
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"sync/atomic"
	"time"
)

var (
	errAlreadyCache = errors.New("already alloc cache")
)

/**
CREATE TABLE `int32_seq` (
  `id` int(11) NOT NULL,
  `next_id` bigint(20) DEFAULT NULL,
  `cache` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='uid sequence';

CREATE TABLE `int64_seq` (
  `id` int(11) NOT NULL,
  `next_id` bigint(20) DEFAULT NULL,
  `cache` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='uid sequence';

**/

var sequenceTemplate = `
CREATE TABLE IF NOT EXISTS %s.%s ( 
  id int(11) NOT NULL,
  next_id bigint(20) DEFAULT NULL,
  cache bigint(20) DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='uid sequence';
`

// MyUID - mysql uid gen
type MyUID struct {
	config   *uid.Config
	done     chan struct{}
	current  int64
	remain   int64
	mutex    sync.Mutex
	db       *sql.DB
	reqInt32 bool
	reqInt64 bool
}

// NewMyUID - create new uid with default collection name.
func NewMyUID() uid.UID {
	return &MyUID{
		done: make(chan struct{}),
	}
}

// HasInt32 - has int32 uid
func (r *MyUID) HasInt32() bool {
	return true
}

func (r *MyUID) next() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if atomic.LoadInt64(&r.remain) > 0 {
		return errAlreadyCache
	}

	txn, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		_ = txn.Rollback()
	}()

	query := fmt.Sprintf("select next_id, cache from %s where id = 0 for update", r.config.TableName)
	rows, err := txn.Query(query)
	if err != nil {
		return err
	}

	var nextID int64 = r.config.InitValue
	var cache int64 = r.config.Step

	if !rows.Next() {
		if err := rows.Close(); err != nil {
			return err
		}

		sql := fmt.Sprintf("insert into %s (id, next_id, cache) values(0, %d, %d)", r.config.TableName, nextID, cache)
		result, err := txn.Exec(sql)
		if err != nil {
			return err
		}

		rowAffected, err := result.RowsAffected()
		if err != nil || rowAffected != 1 {
			return err
		}
	} else {
		if err := rows.Scan(&nextID, &cache); err != nil {
			return rows.Close()
		}
		if err := rows.Close(); err != nil {
			return err
		}
	}

	newLast := nextID + cache

	// update new id
	query = fmt.Sprintf("update %s set next_id = %d where id = 0", r.config.TableName, newLast)
	result, err := txn.Exec(query)
	if err != nil {
		return err
	}

	rowAffected, err := result.RowsAffected()
	if err != nil || rowAffected != 1 {
		return err
	}

	atomic.StoreInt64(&r.current, nextID-1)
	atomic.StoreInt64(&r.remain, cache)
	return txn.Commit()
}

// NextUID32 - next int32 uid
func (r *MyUID) NextUID32() int32 {
	if r.reqInt64 {
		return -1
	}

	r.reqInt32 = true

	if atomic.LoadInt64(&r.remain) > 0 && atomic.AddInt64(&r.remain, -1) >= 0 {
		return int32(atomic.AddInt64(&r.current, 1))
	}

	// request new cache
	if err := r.next(); err != nil && err != errAlreadyCache {
		return -1
	}

	if atomic.LoadInt64(&r.remain) > 0 && atomic.AddInt64(&r.remain, -1) >= 0 {
		return int32(atomic.AddInt64(&r.current, 1))
	}

	return -1
}

// NextUID64 - next int64 uid
func (r *MyUID) NextUID64() int64 {
	if r.reqInt32 {
		return -1
	}

	r.reqInt64 = true

	if atomic.LoadInt64(&r.remain) > 0 && atomic.AddInt64(&r.remain, -1) >= 0 {
		return atomic.AddInt64(&r.current, 1)
	}

	// request new cache
	if err := r.next(); err != nil && err != errAlreadyCache {
		return -1
	}

	if atomic.LoadInt64(&r.remain) > 0 && atomic.AddInt64(&r.remain, -1) >= 0 {
		return atomic.AddInt64(&r.current, 1)
	}

	return -1
}

// Close - close connection
func (r *MyUID) Close() error {
	// release resource
	r.done <- struct{}{}
	return r.db.Close()
}

// StartAndGC start uid adapter.
// config is like {"node":"11"}
// so no gc operation.
func (r *MyUID) StartAndGC(config uid.Config) error {
	r.config = &config

	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?&timeout=30s", r.config.User, r.config.Password, r.config.Server, r.config.Port, r.config.DBName))

	if err != nil {
		return err
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour * 8)

	if _, err := db.Exec("select now()"); err != nil {
		return err
	}

	// check table
	if config.AutoCreateTable {
		sql := fmt.Sprintf(sequenceTemplate, config.DBName, config.TableName)
		if _, err := db.Exec(sql); err != nil {
			return fmt.Errorf("prepare init table err:%v sql:%v", err, sql)
		}
	}

	r.db = db
	go func() {
		ticker := time.NewTicker(time.Second * 30)
		for {
			select {
			case <-ticker.C:
				if err := db.Ping(); err != nil {
					fmt.Printf("ping error:%v", err)
				}
			case <-r.done:
				fmt.Printf("context done\n")
				return
			}
		}
	}()

	return nil
}

func init() {
	uid.Register(uid.TypeMySQL, NewMyUID)
}
