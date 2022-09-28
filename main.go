package main

import (
	"fmt"

	"github.com/mmuflih/golib/filter"
)

func main() {
	// f := filter.Where{"name": filter.W{"like": "%abc%"}, "q": filter.W{"deleted_at": nil}}
	f := filter.Where{"deleted_at": nil}
	// f["created_at"] = filter.W{"raw": "between (a and b)"}
	// f := filter.Where{}
	f["ilike"] = filter.W{"name": "%a%", "label": "%a%"}
	// for i := 0; i < 10; i++ {
	fmt.Println(f.GenerateConditionRaw())
	// }
}
