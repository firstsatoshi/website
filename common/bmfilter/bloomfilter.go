package bmfilter

import (
	"fmt"
	"strings"
	"sync"

	"github.com/zeromicro/go-zero/core/bloom"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// BloomFilter is a thread safe bloom filter wrapper
type BloomFilter struct {
	rwLock      *sync.RWMutex // protect bloomFilter
	bloomFilter bloom.Filter
}

func NewUpgwBloomFilter(redis *redis.Redis, coinType string) *BloomFilter {

	cointype := strings.ToLower(coinType)
	key := fmt.Sprintf("depositaddressbloomfilter:%v", cointype)
	return &BloomFilter{
		bloomFilter: *bloom.New(redis, key, 20*100000),
		rwLock:      new(sync.RWMutex),
	}

}
func (b *BloomFilter) Add(data []byte) error {
	// get write lock
	b.rwLock.Lock()
	defer b.rwLock.Unlock()

	return b.bloomFilter.Add(data)
}

func (b *BloomFilter) Exists(data []byte) (bool, error) {
	// get read lock
	b.rwLock.RLock()
	defer b.rwLock.RUnlock()

	return b.bloomFilter.Exists(data)
}
