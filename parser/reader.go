package parser

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

type BoxHeader struct{
	Size uint32
	Type [4]byte
}

func (b BoxHeader) BoxType() string{
	return string(b.Type[:])
}

func ReadBoxHeader(file io.Reader,b *BoxHeader) error {
	err:=binary.Read(file,binary.BigEndian,b)
	return err
}

func printBoxHeader(h *BoxHeader){
	fmt.Printf("Header size: %d\n",h.Size)
	fmt.Printf("Header type: %s\n",h.BoxType())

}

func ReadFile(filePath string){
	file,err:=os.Open(filePath)
	if err!=nil{
		log.Fatalf("Failed to open file: %s",err)
	}

	defer file.Close()

	header:=&BoxHeader{}
	for{
		if err:=ReadBoxHeader(file,header);err==nil{
			printBoxHeader(header)
			if _,err:= file.Seek(int64(header.Size)-8,io.SeekCurrent); err!=nil{
				log.Fatalf("Failed to move read offset: %v",err)
			}
		}else if err==io.EOF{
			break;
		}else{
			log.Fatalf("Failed to read Box header: %v",err)
		}
	}

	fmt.Println("Read the entire file!")

}




