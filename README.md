# gocsv

go csv helper

Features
---------

* Read for struct/map/parser
* Support encoding

Usage
---------

```go

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


```

csv file:

    Goods Id,Goods Name,价格
    id,name,price
    int,string,float
    1,Apple,5999.99
    2,小米,3.89



output:

    ------------- read map  -------------
    []map[string]interface {}{map[string]interface {}{"id":1, "name":"Apple", "price":5999.99}, map[string]interface {}{"id":2, "name":"小米", "price":3.89}}

    ------------- read object  -------------
    []main.Goods{main.Goods{ID:1, Name:"Apple", Cost:5999.99}, main.Goods{ID:2, Name:"小米", Cost:3.89}}

    ------------- read parser  -------------
    -line 1
    gocsv.Field{Name:"id", Value:"1", Kind:"int"}
    gocsv.Field{Name:"name", Value:"Apple", Kind:"string"}
    gocsv.Field{Name:"price", Value:"5999.99", Kind:"float"}
    -line 2
    gocsv.Field{Name:"id", Value:"2", Kind:"int"}
    gocsv.Field{Name:"name", Value:"小米", Kind:"string"}
    gocsv.Field{Name:"price", Value:"3.89", Kind:"float"}
