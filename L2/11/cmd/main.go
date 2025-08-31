package main

import (
	"fmt"

	"github.com/S1riyS/wildberries-techschool/L2/11/pkg/anagrams"
)

func main() {
	data := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
	fmt.Println(anagrams.Find(data))
}
