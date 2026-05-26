package main

import (
	"fmt"
	"os"
	"github.com/devanshu0x/zentro/parser"
)

func main(){
	if len(os.Args)<2{
		fmt.Println("Usage of tool: zentro <path to target mp4 file>")
		return 
	}

	filePath:=os.Args[1]

	info,err:=os.Stat(filePath)
	if err!=nil{
		fmt.Printf("NO such file exists at given path: %s\n",filePath)
		return 
	}
	if info.IsDir(){
		fmt.Printf("Its a directory path not a file: %s",filePath)
		return 
	}

	parser.ReadFile(filePath)

}