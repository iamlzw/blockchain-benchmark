package tools

import (
	"encoding/csv"
	"log"
	"os"
)

func ReadTestData() [][]string{
	//准备读取文件
	fileName := "./test.csv"
	fs, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("can not open the file, err is %+v", err)
	}
	defer fs.Close()
	fs1, _ := os.Open(fileName)
	r1 := csv.NewReader(fs1)
	content, err := r1.ReadAll()
	if err != nil {
		log.Fatalf("can not read all, err is %+v", err)
	}
	return content
}
