package main

import (
	"os"
	"strings"
	"log"
	"path/filepath"
	"flag"
	"io/ioutil"
	"path"
	"unicode"
	"github.com/foolin/gocsv"
	"encoding/json"
)

var csvpath = flag.String("csv", "", "exmaple: xxx/data/demo.csv or dir: xxx/data")
var outpath = flag.String("out", "", "exmaple: xxx/data/demo.json or dir: xxx/out")

func main() {
	//abs, err := filepath.Abs("./../")
	//log.Printf(filepath.Base(abs))
	flag.Parse()
	if *csvpath == "" {
		flag.Usage()
		return
	}
	fileInfo, err := os.Stat(*csvpath)
	if err != nil {
		log.Panic(err)
		return
	}
	isOutOneFile := false
	if *outpath != "" && (strings.ToLower(filepath.Ext(*outpath)) == ".json" || strings.ToLower(filepath.Ext(*outpath)) == ".cson") {
		isOutOneFile = true
	}
	if *outpath == "" {
		*outpath = *csvpath
	}

	mapAllList := make(map[string][][]string, 0)
	if fileInfo.IsDir() {
		infos, err := ioutil.ReadDir(*csvpath)
		if err != nil {
			log.Panic(err)
			return
		}
		if *outpath == "" {
			*outpath = *csvpath
		}
		for _, info := range infos {
			ext := filepath.Ext(info.Name())
			if ext != ".csv" {
				continue
			}

			name := filename(info.Name());
			list, err := gocsv.ReadLines(path.Join(*csvpath, info.Name()), true)
			if err != nil {
				log.Fatalf("read csv: %v, error: %v", info.Name(), err)
				return
			}
			if isOutOneFile {
				mapAllList[name] = list
			} else {
				jsonfile := strings.Replace(info.Name(), ".csv", ".json", -1)
				outFile := path.Join(*outpath, jsonfile)
				err = writeJsonFile(outFile, list)
				if err != nil {
					log.Fatalf("write file: %v error: %v", outFile, err)
					return
				}
				log.Printf("write file: %v", outFile)
			}

		}

	} else {
		name := upper(filename(fileInfo.Name()));
		list, err := gocsv.ReadLines(path.Join(*csvpath, fileInfo.Name()), true)
		if err != nil {
			log.Fatalf("read csv error: %v", err)
			return
		}
		if isOutOneFile {
			mapAllList[name] = list
		} else {
			jsonfile := strings.Replace(fileInfo.Name(), ".csv", ".json", -1)
			outFile := path.Join(*outpath, jsonfile)
			err = writeJsonFile(outFile, list)
			if err != nil {
				log.Fatalf("write file: %v error: %v", outFile, err)
				return
			}
			log.Printf("write file: %v", outFile)
		}
	}

	if isOutOneFile {
		outFile := *outpath
		err = writeJsonFile(outFile, mapAllList)
		if err != nil {
			log.Fatalf("write file: %v error: %v", outFile, err)
			return
		}
		log.Printf("write file: %v", outFile)
	}

	log.Print("generator done!")
}

func writeJsonFile(outFile string, data interface{}) error {
	//mkdir
	outAbs, _ := filepath.Abs(outFile)
	outDir := filepath.Dir(outAbs)
	err := os.MkdirAll(outDir, 0755)
	if err != nil {
		return err
	}
	jsonContent, err := json.Marshal(data)
	if err != nil {
		return err
	}
	//write file
	err = ioutil.WriteFile(outFile, jsonContent, 0755)
	if err != nil {
		return err
	}
	return nil
}

func upper(str string) string {
	if str == "" {
		return str
	}
	ret := make([]rune, 0)
	isNeedUpper := true        //首字母大写
	for _, c := range str {
		//log.Printf("%#v\n", c)
		if c == '_' {
			isNeedUpper = true        //下划线大写
			continue
		}
		if isNeedUpper {
			ret = append(ret, unicode.ToUpper(c))
		} else {
			ret = append(ret, c)
		}
		isNeedUpper = false
	}
	return string(ret)
}

func filename(filename string) string {
	name := filepath.Base(filename)
	return strings.TrimSuffix(name, filepath.Ext(filename))
}