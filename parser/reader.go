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

type MovieHeaderv0 struct{
	Flags [3]byte
	CreationTime uint32
	ModificationTime uint32
	Timescale uint32
	Duration uint32
}

func (m MovieHeaderv0) Print(){
	fmt.Println("Movie header version 0")
	fmt.Printf("Creation time: %d\n",m.CreationTime)
	fmt.Printf("Modification time: %d\n",m.ModificationTime)
	fmt.Printf("Timescale: %d\n",m.Timescale)
	fmt.Printf("Duration: %d\n",m.Duration)
	fmt.Printf("Duration in seconds: %.2f\n",m.GetDurationInSeconds())
}

func (m MovieHeaderv0) GetDurationInSeconds() float64 {
	return float64(m.Duration)/float64(m.Timescale)
}

type MovieHeaderv1 struct{
	Flags [3]byte
	CreationTime uint64
	ModificationTime uint64
	Timescale uint32
	Duration uint64
}

func (m MovieHeaderv1) Print(){
	fmt.Println("Movie header version 1")
	fmt.Printf("Creation time: %d\n",m.CreationTime)
	fmt.Printf("Modification time: %d\n",m.ModificationTime)
	fmt.Printf("Timescale: %d\n",m.Timescale)
	fmt.Printf("Duration: %d\n",m.Duration)
	fmt.Printf("Duration in seconds: %.2f\n",m.GetDurationInSeconds())
}

func (m MovieHeaderv1) GetDurationInSeconds() float64 {
	return float64(m.Duration)/float64(m.Timescale)
}

func (b BoxHeader) BoxType() string {
	return string(b.Type[:])
}

func ReadBlock(file io.Reader, b interface{}) error {
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

func ReadMovieHeaderData(file io.ReadSeeker, end int64) error {
	var version uint8
	err:=binary.Read(file,binary.BigEndian,&version)
	if err!=nil{
		return fmt.Errorf("Failed to read version in mvhd: %w",err)
	}
	switch version{
	case 0:
		mvhd:= &MovieHeaderv0{}
		if err:=ReadBlock(file,mvhd);err!=nil{
			return fmt.Errorf("Failed to read mvhdv0 struct: %w",err)
		}
		mvhd.Print()

	case 1:
		mvhd:= &MovieHeaderv1{}
		if err:=ReadBlock(file,mvhd);err!=nil{
			return fmt.Errorf("Failed to read mvhdv1 struct: %w",err)
		}
		mvhd.Print()
	}

	 _,err=file.Seek(end,io.SeekStart)
	return err
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
		if readErr := ReadBlock(file, header); readErr != nil {
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

			switch boxType:= header.BoxType(); boxType{
			case "mvhd":
				ReadMovieHeaderData(file,curr+int64(header.Size))
			default:
				_, err := file.Seek(int64(header.Size)-8, io.SeekCurrent)
				if err != nil {
					return err
				}
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
