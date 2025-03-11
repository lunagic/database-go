package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func DD(v any) {
	Dump(v)
	os.Exit(22)
}

func Dump(v any) {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
