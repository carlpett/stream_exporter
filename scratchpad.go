package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func main() {
	v := make(map[string]interface{})
	err := json.Unmarshal([]byte("{\"hello\":{\"there\": {\"cruel\": \"world\"}}}"), &v)
	if err != nil {
		panic(err)
	}
	pathSegments := strings.Split("hello.there.cruel", ".")

	extractor := func(input map[string]interface{}) []interface{} {
		res := make([]interface{}, 0)
		current := input
		fmt.Printf("%v\n", current)
		for _, segment := range pathSegments {
			fmt.Printf("Looking for %s in %v\n", segment, current)
			current, ok := current[segment]
			fmt.Printf("Found %v\n", current)
			if !ok {
				fmt.Printf("Could not find %s in %v\n", segment, current)
				return res
			}
			fmt.Printf("Found %v\n", current)
		}

		res = append(res, current)

		return res
	}

	result := extractor(v)
	fmt.Println(result)
}
