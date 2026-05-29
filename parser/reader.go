package parser

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

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
	if headerType == "moov" || headerType == "trak" || headerType == "mdia" {
		return true
	}
	return false
}

func ReadHeaderDataVersioned(file io.ReadSeeker, end int64, v0 Printable, v1 Printable) error {
	fullBox := FullBox{}
	err := binary.Read(file, binary.BigEndian, &fullBox)
	if err != nil {
		return fmt.Errorf("Failed to read fullbox: %w", err)
	}
	switch fullBox.Version {
	case 0:
		if err := ReadBlock(file, v0); err != nil {
			return fmt.Errorf("Failed to read %T struct: %w",v0 ,err)
		}
		v0.Print()

	case 1:
		if err := ReadBlock(file, v1); err != nil {
			return fmt.Errorf("Failed to read %T struct: %w",v1 ,err)
		}
		v1.Print()
	}

	_, err = file.Seek(end, io.SeekStart)
	return err
}

func ReadHeaderData(file io.ReadSeeker, end int64, h Printable) error {
	err:=binary.Read(file,binary.BigEndian,h)
	if err!=nil{
		return err
	}
	h.Print()

	_, err = file.Seek(end, io.SeekStart)
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

			switch boxType := header.BoxType(); boxType {
			case "mvhd":
				ReadHeaderDataVersioned(file, curr+int64(header.Size),&MovieHeaderv0{},&MovieHeaderv1{})
			case "tkhd":
				ReadHeaderDataVersioned(file, curr+int64(header.Size),&TrackHeaderv0{},&TrackHeaderv1{})
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
