package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/axgle/mahonia"
	"io"
	"log"
	"os"
	"strings"
)

const (
	LINE_RDWR_SIZE = 4096 * 1024
	projectApi     = "project-api"
	wikiApi        = "wiki-api"
	projectWeb     = "project-web"
	wikiWeb        = "wiki-web"
)

var directories = []string{projectApi, wikiApi, projectWeb, wikiWeb}

type translateLine struct {
	directory string
	tag       string
	value     string
}

func readFile(filePath string, region string) []translateLine {
	res := make([]translateLine, 0)
	for _, dir := range directories {
		fileName := filePath + "/" + dir + "/" + region + ".json"
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatal("Open file Failed", err)
			return res
		}
		defer file.Close()

		br := bufio.NewReaderSize(file, LINE_RDWR_SIZE)
		lineNum := 0
		for {
			byteLine, _, c := br.ReadLine()
			lineNum++
			if c == io.EOF {
				break
			}
			strLine := string(byteLine)
			if len(strLine) <= 0 {
				continue
			}
			//fmt.Println(lineNum, len(strLine), strLine)
			if strLine[0] == '{' || strLine[0] == '}' {
				continue
			}
			if strLine[len(strLine)-1] == ',' {
				strLine = strLine[:len(strLine)-1]
			}
			splitIdx := strings.Index(strLine, ": ")
			if splitIdx == -1 {
				fmt.Printf("%s split error", strLine)
				continue
			}
			tag := strings.TrimSpace(strLine[:splitIdx])
			val := strings.TrimSpace(strLine[splitIdx+2 : len(strLine)-1])

			// "tag"
			tag = tag[1 : len(tag)-1]
			// "val",
			val = val[1:len(val)]

			res = append(res, translateLine{directory: dir, tag: tag, value: val})
		}
	}
	return res
}

func saveAsNewFile(fileName string, data map[string]string) {
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	for key := range data {
		_, err := file.WriteString(key + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
}

func saveAsNewFileJson(fileName string, data []translateLine) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal("create file failed", err)
		return
	}
	defer file.Close()
	bw := bufio.NewWriterSize(file, LINE_RDWR_SIZE)
	bw.Write([]byte("{\n"))
	lastIndex := len(data) - 1
	for i, lineData := range data {
		suffix := "\",\n"
		if i == lastIndex {
			suffix = "\"\n"
		}
		wl := "  \"" + lineData.tag + "\": \"" + lineData.value + suffix
		enc := mahonia.NewDecoder("UTF-8")
		resData := enc.ConvertString(string(wl))
		bw.WriteString(resData)
		//bw.WriteString(wl)
	}
	bw.WriteString("}")
	bw.Flush()
	//fmt.Printf("generated file:%s line:%d\n ", fileName, len(data))
}

var filePath string
var region string

func init() {
	flag.StringVar(&filePath, "filePath", "/Users/chenweiqiang/go/src/github.com/bangwork/my_tools/citiao", "文件路径")
	flag.StringVar(&region, "region", "en", "语言类型")
	flag.Parse()
}

func main() {
	data := readFile(filePath, region)

	resTagEqualValueNoEqualData := make([]translateLine, 0)
	resTagNoEqualValueMap := make(map[string]string)
	resTagValueMap := make(map[string]string)
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data); j++ {
			if data[i].tag == data[j].tag && data[i].value != data[j].value {
				_, ok := resTagNoEqualValueMap[data[i].directory+"_"+data[i].tag]
				if !ok {
					resTagNoEqualValueMap[data[i].directory+"_"+data[i].tag] = data[i].value
					resTagValueMap[data[i].tag] = data[i].value
					resTagEqualValueNoEqualData = append(resTagEqualValueNoEqualData, translateLine{tag: data[i].directory + "_" + data[i].tag, value: data[i].value})
				}
			}
		}
	}

	resTagEqualValueData := make([]translateLine, 0)
	resTagEqualValueMap := make(map[string]string)
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data); j++ {
			if data[i].tag != data[j].tag || (data[i].tag == data[j].tag && data[i].value == data[j].value) {
				_, ok := resTagEqualValueMap[data[i].tag]
				_, ok1 := resTagValueMap[data[i].tag]
				if !ok && !ok1 {
					resTagEqualValueMap[data[i].tag] = data[i].value
					resTagEqualValueData = append(resTagEqualValueData, translateLine{tag: data[i].tag, value: data[i].value})
				}
			}
		}
	}

	// 输出到文件-tag相同value相同（只输出一次）或者 tag不同 的有效key值
	saveAsNewFileJson("/Users/chenweiqiang/go/src/github.com/bangwork/my_tools/output/efficient.json", resTagEqualValueData)
	// 输出到文件-tag相同但是value不相同
	saveAsNewFileJson("/Users/chenweiqiang/go/src/github.com/bangwork/my_tools/output/conflict.json", resTagEqualValueNoEqualData)
}
