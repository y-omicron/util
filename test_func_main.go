package main

import (
	"fmt"
	"github.com/y-omicron/util/Util"
	"strings"
)

func main() {
	aaa := Util.HttpXFileVerify(false, "172.24.189.26:8888", "http://127.0.0.1:8080", 50)
	for _, key := range aaa {
		fmt.Printf("%s\n", strings.Join(key, ", "))
	}
}
