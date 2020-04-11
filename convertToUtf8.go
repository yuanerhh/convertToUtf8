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

	// 控制台输出所有log信息
	log.Verbose = true

	files := os.Args[1:]
	for _, filePattern := range files {
		fileList, _ := getFileList(filePattern)
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

			fmt.Println(charCode)
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

func getFileList(filename string) ([]string, error) {
	var res = make([]string, 0, 10)
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Error("ioutil.ReadDir failed: %s", err)
		return nil, err
	}

	for _, file :=  range files {
		//fmt.Printf("isdir: %v, filename: %v\n", file.IsDir(), file.Name())
		if file.IsDir() {
			continue
		}
		if match, _ := filepath.Match(filename, file.Name()); match{
			res = append(res, file.Name())
		}
	}
	return res, nil
}
