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


	//======================= read list struct ===================//
	fmt.Println("\n------------- read list struct -------------")
	var list1 []Goods
	//data.csv ANSI(excel default)
	err = gocsv.ReadList("data.csv", true, &list1)
	if err != nil {
		fmt.Printf("read error: %v", err)
		return
	}
	fmt.Printf("%#v\n", list1)

	//======================= read list ptr ===================//
	fmt.Println("\n------------- read list ptr  -------------")
	var list2 []*Goods
	//data.csv ANSI(excel default)
	err = gocsv.ReadList("data.csv", true, &list2)
	if err != nil {
		fmt.Printf("read error: %v", err)
		return
	}
	fmt.Printf("%#v\n", list2)


	//======================= read map struct ===================//
	fmt.Println("\n------------- read map struct -------------")
	var map1 map[int]Goods
	//data.csv ANSI(excel default)
	err = gocsv.ReadMap("data.csv", true, "id", &map1)
	if err != nil {
		fmt.Printf("read error: %v", err)
		return
	}
	fmt.Printf("%#v\n", map1)

	//======================= read map ptr ===================//
	fmt.Println("\n------------- read map ptr  -------------")
	var map2 map[int]*Goods
	//data.csv ANSI(excel default)
	err = gocsv.ReadMap("data.csv", true, "id", &map2)
	if err != nil {
		fmt.Printf("read error: %v", err)
		return
	}
	fmt.Printf("%#v\n", map2)

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
