package main

import "fmt"

type Shard struct {
	name  string
	start int
	end   int
	limit int
}

func selectShard(productId int, shards []Shard) Shard {
	// Используем бинарный поиск
	l := 0
	r := len(shards) - 1
	for r-l >= 0 {
		m := (l + r) / 2
		if (productId >= shards[m].start) && (productId <= shards[m].end) {
			return shards[m]
		} else if productId <= shards[m].start {
			r = m - 1
		} else {
			l = m + 1
		}
	}

	return Shard{}
}

func shardsInfo() []Shard {
	return []Shard{
		{"1", 1, 20_000, 20_000},
		{"2", 20_001, 40_000, 20_000},
		{"3", 40_001, 60_000, 20_000},
		{"4", 60_001, 80_000, 20_000},
		{"5", 80_001, 90_000, 10_000},
		{"6", 90_001, 100_000, 10_000},
	}
}

func main() {
	shards := shardsInfo()
	productId := 55_555
	shard := selectShard(productId, shards)
	fmt.Println(shard.name) //3
}
