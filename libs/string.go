package libs

import (
	"strconv"
)

/**
 * Created by Muhammad Muflih Kholidin
 * https://github.com/mmuflih
 * muflic.24@gmail.com
 **/

func AddZero(num, length int) string {
	numS := strconv.Itoa(num)
	numLength := len(numS)
	zero := ""
	for i := numLength; i < length; i++ {
		zero += "0"
	}
	return zero + numS
}
