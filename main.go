package main

import (
	"fmt"

	"github.com/mmuflih/golib/filter"
)

func main() {
	f := filter.Where{"name": filter.W{"like": "%abc%"}, "q": filter.W{"deleted_at": nil}}
	fmt.Println(f.GenerateConditionRaw())
}
