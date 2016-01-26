package gocsv

import (
	"fmt"
	"encoding/csv"
	"errors"
	"strconv"
	"os"
	"golang.org/x/text/transform"
	"golang.org/x/text/encoding/simplifiedchinese"
	"reflect"
	"strings"
)

//Field field info
type Field struct{
	Name string
	Value string
	Kind string
}

//Read csv
func Read(file string, isUtf8 bool,  parser func([]Field) error) error  {
	//open file
	fi, err := os.Open(file)
	defer fi.Close()
	//get reader
	var reader *csv.Reader
	if isUtf8{
		reader = csv.NewReader(fi)
	} else{
		//transform gbk to utf8
		r := transform.NewReader(fi, simplifiedchinese.GBK.NewDecoder())
		reader = csv.NewReader(r)
	}

	lines, err := reader.ReadAll()
	if err != nil {
		return err
	}
	lineNum := len(lines)
	if (lineNum < 3) {
		return errors.New(fmt.Sprintf("Csv %v is invalid"))
	}
	names, kinds := lines[1], lines[2]
	fieldNum := len(names)
	//从第三行开始
	for i := 3; i < lineNum; i++ {
		line := lines[i]
		itemFields := make([]Field, fieldNum, fieldNum)
		for j := 0; j < fieldNum; j++ {
			itemField := Field{
				Name: names[j],
				Value: line[j],
				Kind: kinds[j],

			}
			itemFields[j] = itemField
		}
		perr := parser(itemFields)
		//如果返回解析错误，则跳过，直接返回
		if perr != nil{
			return perr
		}
	}
	return  nil
}

//ReadMap read for map
func ReadMap(file string, isUtf8 bool) ([]map[string]interface{}, error) {
	list := make([]map[string]interface{}, 0);
	err := Read(file, isUtf8, func(fields []Field) error {
		item := make(map[string]interface{})
		for _, f := range fields{
			if len(f.Name) <= 0 {
				continue
			}
			var itemValue interface{}
			var innerr error
			switch f.Kind {
			case "int":
				itemValue, innerr = strconv.ParseInt(f.Value, 10, 64)
				if innerr != nil{
					itemValue = 0
				}
			case "float":
				itemValue, innerr = strconv.ParseFloat(f.Value, 64)
				if innerr != nil{
					itemValue = 0
				}
			default:
				itemValue = f.Value
			}
			item[f.Name] = itemValue
		}
		list = append(list, item)
		return nil
	})
	return list, err
}

//ReadObject read for struct
func ReadObject(file string, isUtf8 bool, out interface{}) error {

	if out == nil {
		return errors.New("Cannot remake from <nil>")
	}

	outv := reflect.ValueOf(out)

	outt := outv.Type()
	outk := outt.Kind()

	if outk != reflect.Ptr {
		return errors.New("Cannot reflect into non-pointer")
	}
	slicev := outv.Elem()
	slicet := slicev.Type()
	slicek := slicev.Kind()

	if slicek != reflect.Slice {
		return errors.New("Pointer must point to a slice")
	}

	elmt := slicet.Elem()

	//map field => value
	idxs := make(map[string]int)
	for i := 0; i < elmt.NumField(); i++ {
		name := elmt.Field(i).Tag.Get("csv")
		if len(name) <= 0 {
			name = elmt.Field(i).Name
		}
		idxs[format(name)] = i
	}

	err := Read(file, isUtf8, func(fields []Field) error {
		elmv := reflect.Indirect(reflect.New(elmt))
		for _, f := range fields{
			if len(f.Name) <= 0 {
				continue
			}
			idx, ok := idxs[format(f.Name)]
			if !ok{
				continue
			}
			switch f.Kind {
			case "int":
				itemValue, innerr := strconv.ParseInt(f.Value, 10, 64)
				if innerr != nil{
					itemValue = 0
				}
				elmv.Field(idx).SetInt(itemValue)
			case "float":
				itemValue, innerr := strconv.ParseFloat(f.Value, 64)
				if innerr != nil{
					itemValue = 0
				}
				elmv.Field(idx).SetFloat(itemValue)
			default:
				itemValue := f.Value
				elmv.Field(idx).SetString(itemValue)
			}
		}
		slicev.Set(reflect.Append(slicev, elmv))
		return nil
	})

	return err
}

//format format name
func format(name string) string {
	return fmt.Sprintf("%v%v", strings.ToLower(name[0:1]), name[1:])
}