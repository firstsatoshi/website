package bmfilter_test

import (
	"testing"

	"github.com/firstsatoshi/website/common/bmfilter"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

func TestBloomFilter(t *testing.T) {

	// cfg := conf.
	// CacheRedis:
	//   - Host:
	// Pass: Tisd28478fhkhKSDFsdfk

	cfg := redis.RedisConf{
		Type: "node",
		Host: "172.29.0.1:36379",
		Pass: "Tisd28478fhkhKSDFsdfk",
	}

	redis, err := redis.NewRedis(cfg)
	if err != nil {
		panic(err)
	}
	bmf := bmfilter.NewUpgwBloomFilter(redis, "TEST")

	bmf.Add([]byte("aaaaaaaaaaaaaaaaaaaaaa"))
	bmf.Add([]byte("aaaaaaaacccccccccccccc"))
	bmf.Add([]byte("aaaacccccccccccccccccca"))
	bmf.Add([]byte("zzzzzzzzzzzzzzzzzzzzzz"))

	ok, err := bmf.Exists([]byte("hello"))
	if err != nil {
		t.FailNow()
	}

	if ok {
		t.FailNow()
	}

	ok, err = bmf.Exists([]byte("zzzzzzzzzzzzzzzzzzzzzz"))
	if err != nil {
		t.FailNow()
	}

	if !ok {
		t.FailNow()
	}

}
