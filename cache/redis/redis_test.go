// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package redis

import (
	"testing"
	"time"

	"github.com/dbunion/com/cache"
	"github.com/gomodule/redigo/redis"
)

const (
	cacheKey  = "cacheKey"
	cacheKey1 = "cacheKey1"
)

func TestRedisCache(t *testing.T) {
	bm, err := cache.NewCache(cache.TypeRedisCache, cache.Config{Server: "127.0.0.1", Port: 6379})
	if err != nil {
		t.Fatalf("create new cache error, err:%v", err)
	}

	timeoutDuration := 10 * time.Second
	if err = bm.Put(cacheKey, 1, timeoutDuration); err != nil {
		t.Fatal("put key error", err)
	}
	if !bm.IsExist(cacheKey) {
		t.Fatal("check key exist error")
	}

	time.Sleep(11 * time.Second)

	if bm.IsExist(cacheKey) {
		t.Fatal("check key exist error")
	}
	if err = bm.Put(cacheKey, 1, timeoutDuration); err != nil {
		t.Fatal("put key error", err)
	}

	if v, _ := redis.Int(bm.Get(cacheKey), err); v != 1 {
		t.Fatal("get key error")
	}

	if err = bm.Incr(cacheKey); err != nil {
		t.Fatalf("incr error, err:%v", err)
	}

	if v, _ := redis.Int(bm.Get(cacheKey), err); v != 2 {
		t.Fatal("get key error")
	}

	if err = bm.Decr(cacheKey); err != nil {
		t.Fatal("Decr Error", err)
	}

	if v, _ := redis.Int(bm.Get(cacheKey), err); v != 1 {
		t.Fatalf("get key error, expect:%v actual:%v", 1, v)
	}
	if err = bm.Delete(cacheKey); err != nil {
		t.Fatal("del Error", err)
	}

	if bm.IsExist(cacheKey) {
		t.Fatal("delete err")
	}

	//test string
	if err = bm.Put(cacheKey, "author", timeoutDuration); err != nil {
		t.Fatal("put key error", err)
	}
	if !bm.IsExist(cacheKey) {
		t.Fatal("check key exist error")
	}

	if v, _ := redis.String(bm.Get(cacheKey), err); v != "author" {
		t.Fatal("get key error")
	}

	// test GetMulti
	if err = bm.Put(cacheKey1, "author1", timeoutDuration); err != nil {
		t.Fatal("put key error", err)
	}
	if !bm.IsExist(cacheKey1) {
		t.Fatal("check key exist error")
	}

	vv := bm.GetMulti([]string{cacheKey, cacheKey1})
	if len(vv) != 2 {
		t.Fatal("get multi error")
	}
	if v, _ := redis.String(vv[0], nil); v != "author" {
		t.Fatal("get multi error")
	}
	if v, _ := redis.String(vv[1], nil); v != "author1" {
		t.Fatal("get multi error")
	}

	// test clear all
	if err = bm.ClearAll(); err != nil {
		t.Fatal("clear all err")
	}
}
