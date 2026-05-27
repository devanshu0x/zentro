package parser

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type BoxHeader struct {
	Size uint32
	Type [4]byte
}

func (b BoxHeader) BoxType() string {
	return string(b.Type[:])
}

func ReadBoxHeader(file io.Reader, b *BoxHeader) error {
	err := binary.Read(file, binary.BigEndian, b)
	return err
}

func printBoxHeader(h *BoxHeader, level int) {
	prefeix := strings.Repeat("\t", level)
	fmt.Printf("%sBlock type: %s\n", prefeix, h.BoxType())
	fmt.Printf("%sBlock size: %d\n", prefeix, h.Size)

}

func isContainer(headerType string) bool {
	if headerType == "moov" || headerType == "trak" {
		return true
	}
	return false
}

func ParseBoxes(file io.ReadSeeker, end int64, level int) error {
	header := &BoxHeader{}

	for {
		curr, err := file.Seek(0, io.SeekCurrent)
		if err != nil {
			return err
		}
		if curr >= end {
			break
		}
		if readErr := ReadBoxHeader(file, header); readErr != nil {
			return readErr
		}
		if header.Size < 8 {
			return fmt.Errorf("invalid box size: %d", header.Size)
		}

		printBoxHeader(header, level)
		if isContainer(header.BoxType()) {
			err := ParseBoxes(file, curr+int64(header.Size), level+1)
			if err != nil {
				return err
			}
		} else {
			_, err := file.Seek(int64(header.Size)-8, io.SeekCurrent)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ReadFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}

	defer file.Close()
	// constainerLevel:=0
	// for{
	// 	if err:=ReadBoxHeader(file,header);err==nil{
	// 		printBoxHeader(header,constainerLevel)
	// 		if isContainer(header.BoxType()){
	// 			constainerLevel++
	// 			continue
	// 		}
	// 		constainerLevel--
	// 		if _,err:= file.Seek(int64(header.Size)-8,io.SeekCurrent); err!=nil{
	// 			log.Fatalf("Failed to move read offset: %v",err)
	// 		}
	// 	}else if err==io.EOF{
	// 		break;
	// 	}else{
	// 		log.Fatalf("Failed to read Box header: %v",err)
	// 	}
	// }

	info, err := file.Stat()
	if err != nil {
		log.Fatalf("Failed to load file info: %v", err)
	}
	fileEnd := info.Size()
	if err := ParseBoxes(file, fileEnd, 0); err != nil {
		log.Fatalf("Failed to read the file: %v", err)
	}

	fmt.Println("Read the entire file!")

}
