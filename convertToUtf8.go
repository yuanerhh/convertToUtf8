package main

import (
	"bytes"
	"fmt"
	"github.com/gpmgo/gopm/modules/log"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s + [filename]\n", os.Args[0])
		return
	}

	allFiles, _ := getAllFiles()
	// 控制台输出log信息控制
	// true: 打印log  false: 不打印log
	log.Verbose = false

	files := os.Args[1:]
	for _, filePattern := range files {
		fileList, _ := getFileList(filePattern, allFiles)
		for _, file := range fileList {
			fText, err := ioutil.ReadFile(file)
			if err != nil {
				log.Error("ioutil.ReadFile %s failed: %s", file, err)
				continue
			}

			charCode, err := detectCode(fText)
			if err != nil {
				log.Error("detectCode failed: %s", err)
				continue
			}

			if charCode == "GB-18030" {
				oriFile, err := os.OpenFile(file+".ori", os.O_RDWR | os.O_CREATE, 0666)
				if err != nil {
					log.Error("OpenFile %s failed: %s", file+".ori", err)
					continue
				}

				newFile, err := os.OpenFile(file, os.O_RDWR, 0666)
				if err != nil {
					log.Error("OpenFile %s failed: %s", file, err)
					oriFile.Close()
					continue
				}

				_, err = oriFile.Write(fText)
				if err != nil {
					log.Error("oriFile.Write failed: %s", err)
					oriFile.Close()
					newFile.Close()
					continue
				}

				// github.com/saintfish/chardet 只检测 GB-18030
				// golang.org/x/net/html/charset 只能用gbk
				newContent, err := convertToUtf8(fText, "gbk")
				if err != nil {
					log.Error("convertToUtf8 failed: %s", err)
					oriFile.Close()
					newFile.Close()
					continue
				}

				_, err = newFile.Write(newContent)
				if err != nil {
					log.Error("newFile.Write failed: %s", err)
					oriFile.Close()
					newFile.Close()
					continue
				}

				oriFile.Close()
				newFile.Close()
				fmt.Printf("%s convert from %s to UTF-8 success!\n", file, charCode)
			}
		}
	}

}

func convertToUtf8(src []byte, encode string) ([]byte, error) {
	byteReader := bytes.NewReader(src)
	reader, err := charset.NewReaderLabel(encode, byteReader)
	if err != nil {
		log.Error("charset.NewReaderLabel failed : %s", err)
		return nil, err
	}

	dst, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Error("ioutil.ReadAll failed : %s", err)
		return nil, err
	}
	return dst, nil
}

func detectCode(src []byte) (string, error) {
	detector := chardet.NewTextDetector()
	var result *chardet.Result
	result, err := detector.DetectBest(src)
	if err != nil {
		log.Error("detector.DetectBest failed: %s", err)
		return "", err
	}

	log.Info("charset: %s, language: %s, confidence: %d",
		result.Charset, result.Language, result.Confidence)

	return result.Charset, nil
}

func getFileList(filename string, fileList []string) ([]string, error) {
	var res = make([]string, 0, 10)
	for _, file :=  range fileList {
		if match, _ := filepath.Match(filename, file); match{
			res = append(res, file)
		}
	}
	return res, nil
}

func getAllFiles() ([]string, error) {
	var allFiles = make([]string, 0, 100)
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			allFiles = append(allFiles, path)
		}
		if err != nil {
			log.Error("Walk err: %s", err)
			return err
		}
		return nil
	})

	return allFiles, nil
}
