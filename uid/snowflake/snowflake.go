package snowflake

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/dbunion/com/uid"
)

// Snowflake - Snowflake uuid gen
type Snowflake struct {
	node *snowflake.Node
}

// NewSnowflake create new uid with default collection name.
func NewSnowflake() uid.UID {
	return &Snowflake{}
}

// HasInt32 - has int32 uid
func (s *Snowflake) HasInt32() bool {
	return false
}

// NextUID32 - next int32 uid
func (s *Snowflake) NextUID32() int32 {
	return 0
}

// NextUID64 - next int64 uid
func (s *Snowflake) NextUID64() int64 {
	return s.node.Generate().Int64()
}

// Close - close connection
func (s *Snowflake) Close() error {
	// do nothing
	return nil
}

// StartAndGC start uid adapter.
// config is like {"node":"11"}
// so no gc operation.
func (s *Snowflake) StartAndGC(config uid.Config) (err error) {
	s.node, err = snowflake.NewNode(config.NodeID)
	if err != nil {
		return fmt.Errorf("create new snowflake err:%v", err)
	}

	return nil
}

func init() {
	uid.Register("SnowflakeUID", NewSnowflake)
}
