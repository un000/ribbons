package main

import (
	"encoding/json"
	"fmt"

	"github.com/un000/ribbons"
)

func main() {
	jsonArray := []byte(`[1000, 1001, 1337]`)

	{
		set, err := ribbons.NewFromJSON(jsonArray)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Items: %+v\n", set.List())
		fmt.Printf("Has 1000: %v\n", set.Has(1000))
		fmt.Printf("Len 1000: %v\n", set.Len())

		set.Add(33)
		fmt.Printf("Items after adding 33: %+v\n", set.List())

		set.Delete(33)
		fmt.Printf("Items after deleting 33: %+v\n", set.List())
	}

	fmt.Println("-------")

	// custom json unmarshalling also works
	{
		set := ribbons.New()
		err := json.Unmarshal(jsonArray, &set)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Items: %+v\n", set.List())
		fmt.Printf("Has 1000: %v\n", set.Has(1000))
	}
}
