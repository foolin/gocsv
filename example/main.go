package main
import (
	"github.com/foolin/gocsv"
	"fmt"
)

//Goods goods struct for csv
type Goods struct {
	ID   int `csv:"id"`	//id => ID
	Name string	// name => Name (default, first letter lowercase)
	Price float64 `csv:"cost"`	// rename price => cost
}

func main() {

	var err error

	//======================= read map[string]interface{} ===================//
	fmt.Println("\n------------- read  -------------")
	//datautf8.csv utf8 file
	data, err := gocsv.Read("datautf8.csv", false)
	if err != nil {
		panic(fmt.Sprintf("read error: %v", err))
		return
	}
	fmt.Printf("%#v\n", data)


	//======================= read list ===================//
	fmt.Println("\n------------- read list  -------------")
	var list []Goods
	//data.csv ANSI(excel default)
	err = gocsv.ReadList("data.csv", true, &list)
	if err != nil {
		fmt.Printf("read error: %v", err)
		return
	}
	fmt.Printf("%#v\n", list)


	//======================= read map ===================//
	fmt.Println("\n------------- read map  -------------")
	var vmap map[int]Goods
	//data.csv ANSI(excel default)
	err = gocsv.ReadMap("data.csv", true, "id", &vmap)
	if err != nil {
		fmt.Printf("read error: %v", err)
		return
	}
	fmt.Printf("%#v\n", vmap)

	//======================= read parser ===================//
	fmt.Println("\n------------- read parser  -------------")
	line := 1
	err = gocsv.ReadRaw("data.csv", true, func(fields []gocsv.Field) error {
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
