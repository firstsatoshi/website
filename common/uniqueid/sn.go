package uniqueid

import (
	"fmt"
	"math/rand"
	"time"
)

// 生成sn单号
type SnPrefix string

const (
	KC_RAND_KIND_NUM   = 0 // 纯数字
	KC_RAND_KIND_LOWER = 1 // 小写字母
	KC_RAND_KIND_UPPER = 2 // 大写字母
	KC_RAND_KIND_ALL   = 3 // 数字、大小写字母
)

// 随机字符串
func krand(size int, kind int) string {
	ikind, kinds, result := kind, [][]int{{10, 48}, {26, 97}, {26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return string(result)
}

// 生成单号
func GenSn(snPrefix string) string {
	time.Sleep(time.Nanosecond)
	return fmt.Sprintf("%s%s%s", snPrefix, time.Now().UTC().In(time.UTC).Format("20060102150405"), krand(8, KC_RAND_KIND_NUM))
}
