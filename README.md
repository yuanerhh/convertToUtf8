# convertToUtf8
自动检测非UTF8编码的文件并转换为UTF8编码。  
  
## 安装

一、源码安装  

    go get github.com/gpmgo/gopm 
    go get github.com/saintfish/chardet  
    
    mkdir -p $GOPATH/src/golang.org/x
    cd $GOPATH/src/golang.org/x
    git clone https://github.com/golang/net.git
    git clone https://github.com/golang/text.git
  
    go build convertToUtf8.go  
    go install 
   
  
如果GOBIN路径已追加到系统PATH中，则convertToUtf8可以随意使用。  
  
二、可执行文件安装  
  
bin目录下有各系统平台的编译好的可执行文件，直接将对应的目录追加到系统PATH中即可直接使用convertToUtf8。  
  
## 使用方式

    convertToUtf8 [filename | 文件名通配(如*.cpp) | path/filename | path/*.cpp]  
    
## 注意事项

1. 因为convertToUtf8查询的是当前路径中的所有文件，所以参数只能传入当前目录的文件或者子目录文件。
2. 如果有文件需要进行UTF8转码，convertToUtf8会以filename.ori保留原文件的备份，然后将原文件转换为UTF8编码。
