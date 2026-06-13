package mp4

import (
	"encoding/binary"
	"io"
	"fmt"
)

func Fixed1616ToFloat(a uint32) float64{
	return float64(a)/(1<<16)
}

func DecodeLanguage(lang uint16) string {
	c1 := byte((lang>>10)&0x1F) + 0x60
	c2 := byte((lang>>5)&0x1F) + 0x60
	c3 := byte(lang&0x1F) + 0x60

	return string([]byte{c1, c2, c3})
}

func ReadFullBox(r io.Reader) (*FullBox,error) {
	var fb FullBox

	if err:=binary.Read(r,binary.BigEndian, &fb);err!=nil{
		return nil,err
	}

	return &fb,nil
}

func ParseMvhd(r io.Reader) (*MovieHeader, error){
	fb,err:=ReadFullBox(r)
	if err!=nil{
		return nil,err
	}
	mvhd:=&MovieHeader{
		Version: fb.Version,
	}

	switch fb.Version {
	case 0:
		var creation uint32
		var modification uint32
		var duration uint32

		if err := binary.Read(r, binary.BigEndian, &creation); err != nil {
			return nil, err
		}

		if err := binary.Read(r, binary.BigEndian, &modification); err != nil {
			return nil, err
		}

		if err := binary.Read(r, binary.BigEndian, &mvhd.Timescale); err != nil {
			return nil, err
		}

		if err := binary.Read(r, binary.BigEndian, &duration); err != nil {
			return nil, err
		}
		mvhd.CreationTime=uint64(creation)
		mvhd.ModificationTime= uint64(modification)
		mvhd.Duration=uint64(duration)
	case 1:

		if err := binary.Read(r, binary.BigEndian, &mvhd.CreationTime); err != nil {
			return nil, err
		}

		if err := binary.Read(r, binary.BigEndian, &mvhd.ModificationTime); err != nil {
			return nil, err
		}

		if err := binary.Read(r, binary.BigEndian, &mvhd.Timescale); err != nil {
			return nil, err
		}

		if err := binary.Read(r, binary.BigEndian, &mvhd.Duration); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported mvhd version %d", fb.Version)
	}
	return mvhd,nil
}

func ParseTkhd(r io.Reader) (*TrackHeader, error) {
	fb, err := ReadFullBox(r)
	if err != nil {
		return nil, err
	}

	tkhd := &TrackHeader{
		Version: fb.Version,
	}

	switch fb.Version {

	case 0:
		var creation uint32
		var modification uint32
		var duration uint32
		var reserved [4]byte

		if err := binary.Read(r, binary.BigEndian, &creation); err != nil {
			return nil, err
		}

		if err := binary.Read(r, binary.BigEndian, &modification); err != nil {
			return nil, err
		}

		if err := binary.Read(r, binary.BigEndian, &tkhd.TrackID); err != nil {
			return nil, err
		}

		if err := binary.Read(r, binary.BigEndian, &reserved); err != nil {
			return nil, err
		}

		if err := binary.Read(r, binary.BigEndian, &duration); err != nil {
			return nil, err
		}

		tkhd.CreationTime = uint64(creation)
		tkhd.ModificationTime = uint64(modification)
		tkhd.Duration = uint64(duration)

	case 1:
		var reserved [4]byte

		if err := binary.Read(r, binary.BigEndian, &tkhd.CreationTime); err != nil {
			return nil, err
		}

		if err := binary.Read(r, binary.BigEndian, &tkhd.ModificationTime); err != nil {
			return nil, err
		}

		if err := binary.Read(r, binary.BigEndian, &tkhd.TrackID); err != nil {
			return nil, err
		}

		if err := binary.Read(r, binary.BigEndian, &reserved); err != nil {
			return nil, err
		}

		if err := binary.Read(r, binary.BigEndian, &tkhd.Duration); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported tkhd version %d", fb.Version)
	}

	var reserved2 [8]byte
	var layer int16
	var alternateGroup int16
	var volume int16
	var reserved3 [2]byte
	var matrix [36]byte

	if err := binary.Read(r, binary.BigEndian, &reserved2); err != nil {
		return nil, err
	}

	if err := binary.Read(r, binary.BigEndian, &layer); err != nil {
		return nil, err
	}

	if err := binary.Read(r, binary.BigEndian, &alternateGroup); err != nil {
		return nil, err
	}

	if err := binary.Read(r, binary.BigEndian, &volume); err != nil {
		return nil, err
	}

	if err := binary.Read(r, binary.BigEndian, &reserved3); err != nil {
		return nil, err
	}

	if err := binary.Read(r, binary.BigEndian, &matrix); err != nil {
		return nil, err
	}

	var width uint32
	var height uint32

	if err := binary.Read(r, binary.BigEndian, &width); err != nil {
		return nil, err
	}

	if err := binary.Read(r, binary.BigEndian, &height); err != nil {
		return nil, err
	}

	tkhd.Width = Fixed1616ToFloat(width)
	tkhd.Height = Fixed1616ToFloat(height)

	return tkhd, nil
}

func ParseMdhd(r io.Reader) (*MediaHeader, error) {
	fb, err := ReadFullBox(r)
	if err != nil {
		return nil, err
	}

	mdhd := &MediaHeader{
		Version: fb.Version,
	}

	var language uint16
	var predefined uint16

	switch fb.Version {

	case 0:
		var creation uint32
		var modification uint32
		var duration uint32

		if err := binary.Read(r, binary.BigEndian, &creation); err != nil {
			return nil, err
		}

		if err := binary.Read(r, binary.BigEndian, &modification); err != nil {
			return nil, err
		}

		if err := binary.Read(r, binary.BigEndian, &mdhd.Timescale); err != nil {
			return nil, err
		}

		if err := binary.Read(r, binary.BigEndian, &duration); err != nil {
			return nil, err
		}

		if err := binary.Read(r, binary.BigEndian, &language); err != nil {
			return nil, err
		}

		if err := binary.Read(r, binary.BigEndian, &predefined); err != nil {
			return nil, err
		}

		mdhd.CreationTime = uint64(creation)
		mdhd.ModificationTime = uint64(modification)
		mdhd.Duration = uint64(duration)

	case 1:

		if err := binary.Read(r, binary.BigEndian, &mdhd.CreationTime); err != nil {
			return nil, err
		}

		if err := binary.Read(r, binary.BigEndian, &mdhd.ModificationTime); err != nil {
			return nil, err
		}

		if err := binary.Read(r, binary.BigEndian, &mdhd.Timescale); err != nil {
			return nil, err
		}

		if err := binary.Read(r, binary.BigEndian, &mdhd.Duration); err != nil {
			return nil, err
		}

		if err := binary.Read(r, binary.BigEndian, &language); err != nil {
			return nil, err
		}

		if err := binary.Read(r, binary.BigEndian, &predefined); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported mdhd version %d", fb.Version)
	}

	mdhd.Language = DecodeLanguage(language)

	return mdhd, nil
}

func ParseHdlr(r io.Reader) (*HandlerReference, error) {
	_, err := ReadFullBox(r)
	if err != nil {
		return nil, err
	}

	var preDefined [4]byte
	var handlerType [4]byte

	if err := binary.Read(r, binary.BigEndian, &preDefined); err != nil {
		return nil, err
	}

	if err := binary.Read(r, binary.BigEndian, &handlerType); err != nil {
		return nil, err
	}

	return &HandlerReference{
		HandlerType: string(handlerType[:]),
	}, nil
}

