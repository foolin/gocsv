package main
import (
	"github.com/foolin/gocsv"
	"fmt"
)

//Goods goods struct for csv
type Goods struct {
	ID   int `csv:"id"`	//id => ID
	Name string	// name => Name (default, first letter lowercase)
	Cost float64 `csv:"price"`	// rename price => cost
}

func main() {

	var err error

	//======================= read map[string]interface{} ===================//
	fmt.Println("\n------------- read map  -------------")
	//datautf8.csv utf8 file
	data, err := gocsv.ReadMap("datautf8.csv", true)
	if err != nil {
		panic(fmt.Sprintf("read error: %v", err))
		return
	}
	fmt.Printf("%#v\n", data)


	//======================= read object ===================//
	fmt.Println("\n------------- read object  -------------")
	var out []Goods
	//data.csv ANSI(excel default)
	err = gocsv.ReadObject("data.csv", false, &out)
	if err != nil {
		fmt.Printf("read error: %v", err)
		return
	}
	fmt.Printf("%#v\n", out)

	//======================= read parser ===================//
	fmt.Println("\n------------- read parser  -------------")
	line := 1
	err = gocsv.Read("data.csv", false, func(fields []gocsv.Field) error {
		fmt.Printf("-line %v\n", line)
		for _, f := range fields {
			fmt.Printf("%#v\n", f)
		}
		line = line + 1
		return nil
	})
	if err != nil {
		fmt.Printf("read error: %v", err)
		return
	}

}
