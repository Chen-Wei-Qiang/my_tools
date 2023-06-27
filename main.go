package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/axgle/mahonia"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	LINE_RDWR_SIZE = 4096 * 1024
)

type translateLine struct {
	directory string
	tag       string
	value     string
}

func ListDirectories(directory string) ([]string, error) {
	var directories []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != directory {
			directories = append(directories, info.Name())
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return directories, nil
}

func readFile(filePath string, directories []string, region string) []translateLine {
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
			splitIdx := strings.Index(strLine, ":")
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
var outPath string

func init() {
	flag.StringVar(&filePath, "filePath", "./citiao", "文件路径")
	flag.StringVar(&region, "region", "en", "语言类型")
	flag.StringVar(&outPath, "outPath", "./", "输出路径")
	flag.Parse()
}

func main() {
	directories, err := ListDirectories(filePath)
	if err != nil {
		log.Fatal(err)
	}
	data := readFile(filePath, directories, region)
	fmt.Println("词条总数", len(data))
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

	resTagEqualValueEqualDataOne := make([]translateLine, 0)
	resTagEqualValueMapOnes := make(map[string]int)
	for i := 0; i < len(data); i++ {
		resTagEqualValueMapOnes[data[i].tag]++
	}

	var totle int
	for tag, count := range resTagEqualValueMapOnes {
		if count != 1 {
			_, ok1 := resTagValueMap[tag]
			if !ok1 {
				resTagEqualValueEqualDataOne = append(resTagEqualValueEqualDataOne, translateLine{tag: tag, value: strconv.Itoa(count)})
				totle += count
			}
		}
	}

	// 输出到文件-有效词条文件
	saveAsNewFileJson(outPath+"/efficient.json", resTagEqualValueData)
	fmt.Println("有效词条数据条数", len(resTagEqualValueData))
	// 输出到文件-冲突词条文件
	saveAsNewFileJson(outPath+"/conflict.json", resTagEqualValueNoEqualData)
	fmt.Println("冲突词条数据条数", len(resTagEqualValueNoEqualData))
	// 输出到文件-key相同且value也相同
	saveAsNewFileJson(outPath+"/keyEqualvalue.json", resTagEqualValueEqualDataOne)
	fmt.Println("key且value词条数据相同条数", totle-len(resTagEqualValueEqualDataOne))

	if len(resTagEqualValueData)+len(resTagEqualValueNoEqualData)+totle-len(resTagEqualValueEqualDataOne) == len(data) {
		fmt.Println("条数相同，没有问题～")
	} else {
		fmt.Println("条数不同，请检查～")
	}
}
