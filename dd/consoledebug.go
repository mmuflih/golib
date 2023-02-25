package dd

import (
	"encoding/json"
	"fmt"
)

/**
 * Created by Muhammad Muflih Kholidin
 * https://github.com/mmuflih
 * muflic.24@gmail.com
 **/

func Log(data ...interface{}) {
	o, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(o))
}
