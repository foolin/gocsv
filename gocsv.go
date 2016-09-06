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
type Field struct {
	Name  string
	Value string
	Kind  string
}


//Read read for map array
func Read(file string, isGbk bool) (list []map[string]interface{}, err error) {
	//catch panic
	defer func() {
		if rerr := recover(); rerr != nil {
			err = errors.New(fmt.Sprintf("read csv file: %v, error: %v", file, rerr))
		}
	}()

	list = make([]map[string]interface{}, 0);
	err = ReadRaw(file, isGbk, func(fields []Field) error {
		item := make(map[string]interface{})
		for _, f := range fields {
			if len(f.Name) <= 0 {
				continue
			}
			var itemValue interface{}
			var innerr error
			switch f.Kind {
			case "int":
				itemValue, innerr = strconv.ParseInt(f.Value, 10, 64)
				if innerr != nil {
					itemValue = 0
				}
			case "float":
				itemValue, innerr = strconv.ParseFloat(f.Value, 64)
				if innerr != nil {
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

//ReadList read for []struct
func ReadList(file string, isGbk bool, out interface{}) (err error) {
	//catch panic
	defer func() {
		if rerr := recover(); rerr != nil {
			err = errors.New(fmt.Sprintf("read csv file: %v, error: %v", file, rerr))
		}
	}()

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
	elmIsPtr := false
	//element is ptr
	if elmt.Kind() == reflect.Ptr{
		elmt = elmt.Elem()
		elmIsPtr = true
	}

	//map field => value
	idxs := make(map[string]int)

	for i := 0; i < elmt.NumField(); i++ {
		name := elmt.Field(i).Tag.Get("csv")
		if len(name) <= 0 {
			name = elmt.Field(i).Name
		}
		idxs[format(name)] = i
	}

	err = ReadRaw(file, isGbk, func(fields []Field) error {
		elmv := reflect.Indirect(reflect.New(elmt))
		for _, f := range fields {
			if len(f.Name) <= 0 {
				continue
			}
			idx, ok := idxs[format(f.Name)]
			if !ok {
				continue
			}
			fValue := elmv.Field(idx)
			setValue(&fValue, f)
		}
		if elmIsPtr{
			slicev.Set(reflect.Append(slicev, elmv.Addr()))
		}else{
			slicev.Set(reflect.Append(slicev, elmv))
		}
		return nil
	})

	return err
}


//ReadList read for map[interface{}]struct
func ReadMap(file string, isGbk bool, keyField string, out interface{}) (err error) {
	//catch panic
	defer func() {
		if rerr := recover(); rerr != nil {
			err = errors.New(fmt.Sprintf("read csv file: %v, error: %v", file, rerr))
		}
	}()

	if out == nil {
		return errors.New("Cannot remake from <nil>")
	}

	outv := reflect.ValueOf(out)

	outt := outv.Type()
	outk := outt.Kind()

	if outk != reflect.Ptr {
		return errors.New("Cannot reflect into non-pointer")
	}
	mapv := outv.Elem()
	mapt := mapv.Type()
	mapk := mapv.Kind()

	if mapk != reflect.Map {
		return errors.New("Pointer must point to a slice")
	}

	//make map
	if mapv.IsNil() {
		mapv.Set(reflect.MakeMap(mapt))
	}

	elmt := mapt.Elem()
	elmIsPtr := false
	//element is ptr
	if elmt.Kind() == reflect.Ptr{
		elmt = elmt.Elem()
		elmIsPtr = true
	}

	//map field => value
	idxs := make(map[string]int)
	for i := 0; i < elmt.NumField(); i++ {
		name := elmt.Field(i).Tag.Get("csv")
		if len(name) <= 0 {
			name = elmt.Field(i).Name
		}
		idxs[format(name)] = i
	}

	err = ReadRaw(file, isGbk, func(fields []Field) error {
		elmv := reflect.Indirect(reflect.New(elmt))
		keyi := 0
		isMatchKey := false
		for _, f := range fields {
			if len(f.Name) <= 0 {
				continue
			}
			idx, ok := idxs[format(f.Name)]
			if !ok {
				continue
			}
			if f.Name == keyField {
				keyi = idx
				isMatchKey = true
			}
			fValue := elmv.Field(idx)
			setValue(&fValue, f)
		}
		if !isMatchKey{
			return errors.New(fmt.Sprintf("Primary key not found, \"%v\" not has field name \"%v\"",file, keyField))
		}
		if elmIsPtr{
			mapv.SetMapIndex(elmv.Field(keyi), elmv.Addr())
		}else{
			mapv.SetMapIndex(elmv.Field(keyi), elmv)
		}
		return nil
	})

	return err
}


//Read read csv for handle
func ReadLines(file string, isGbk bool) (lines [][]string, err error) {
	//catch panic
	defer func() {
		if rerr := recover(); rerr != nil {
			err = errors.New(fmt.Sprintf("read csv file: %v, error: %v", file, rerr))
		}
	}()

	//open file
	fi, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fi.Close()
	//get reader
	var reader *csv.Reader
	if !isGbk {
		reader = csv.NewReader(fi)
	} else {
		//transform gbk to utf8
		r := transform.NewReader(fi, simplifiedchinese.GBK.NewDecoder())
		reader = csv.NewReader(r)
	}
	lines, err = reader.ReadAll()
	return
}

func setValue(elmv *reflect.Value, f Field)  {
	switch f.Kind {
	case "int", "int64", "long":
		itemValue, innerr := strconv.ParseInt(f.Value, 10, 64)
		if innerr != nil {
			itemValue = 0
		}
		elmv.SetInt(itemValue)
	case "float", "float64", "double":
		itemValue, innerr := strconv.ParseFloat(f.Value, 64)
		if innerr != nil {
			itemValue = 0
		}
		elmv.SetFloat(itemValue)
	case "bool":
		itemValue, innerr := strconv.ParseBool(f.Value)
		if innerr != nil {
			itemValue = false
		}
		elmv.SetBool(itemValue)
	default:
		itemValue := f.Value
		elmv.SetString(itemValue)
	}
}

//Read read csv for handle
func ReadRaw(file string, isGbk bool, handle func([]Field) error) (err error) {
	//catch panic
	defer func() {
		if rerr := recover(); rerr != nil {
			err = errors.New(fmt.Sprintf("read csv file: %v, error: %v", file, rerr))
		}
	}()
	if file == ""{
		return errors.New("read csv file parameter is empty.")
	}
	lines, err := ReadLines(file, isGbk)
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
				Name: trim(names[j]),
				Value: trim(line[j]),
				Kind: trim(kinds[j]),

			}
			itemFields[j] = itemField
		}
		perr := handle(itemFields)
		//如果返回解析错误，则跳过，直接返回
		if perr != nil {
			return perr
		}
	}
	return nil
}

//format format name
func format(name string) string {
	return fmt.Sprintf("%v%v", strings.ToLower(name[0:1]), name[1:])
}

func trim(s string) string  {
	return strings.TrimSpace(s)
}