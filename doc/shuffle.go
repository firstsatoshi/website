package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {

	boxIds := []int{1, 2, 3, 4, 5, 6}
	rand.Seed(int64(time.Now().UnixNano()))
	rand.Shuffle(len(boxIds), func(i, j int) { boxIds[i], boxIds[j] = boxIds[j], boxIds[i] })
	fmt.Printf("%v", boxIds)
	lockBoxIds := boxIds[:1]
	fmt.Printf("%v", lockBoxIds)

}
