# convertToUtf8
自动检测非UTF8编码的文件并转换为UTF8编码。

## 安装
一、源码安装：
go get github.com/gpmgo/gopm/modules/log
go get github.com/saintfish/chardet
go get golang.org/x/net/html/charset

go build convertToUtf8.go
go install

如果GOBIN路径追加到系统PATH中，则convertToUtf8就可以随意使用。

二、可执行文件安装：
bin目录下有各系统平台的编译好的可执行文件，直接将对应的目录追加到系统PATH中即可直接使用convertToUtf8。

## 使用方式
convertToUtf8 [filename | 文件名通配(如*.cpp)]

