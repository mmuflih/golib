package main

import (
	"fmt"

	"github.com/mmuflih/golib/filter"
)

func main() {
	// f := filter.Where{"name": filter.W{"like": "%abc%"}, "q": filter.W{"deleted_at": nil}}
	f := filter.Where{}
	f["created_at"] = filter.W{"raw": "between (a and b)"}
	f["like"] = filter.W{"name": "%a%", "label": "%a%"}
	fmt.Println(f.GenerateConditionRaw())
}
