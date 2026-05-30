package parser

import (
	"fmt"
	"strings"
)

type Printable interface {
	Print()
}

type BoxHeader struct {
	Size uint32
	Type [4]byte
}

func printBoxHeader(h *BoxHeader, level int) {
	prefeix := strings.Repeat("\t", level)
	fmt.Printf("%sBlock type: %s\n", prefeix, h.BoxType())
	fmt.Printf("%sBlock size: %d\n", prefeix, h.Size)

}

type FullBox struct {
	Version uint8
	Flags   [3]byte
}

type MovieHeaderv0 struct {
	CreationTime     uint32
	ModificationTime uint32
	Timescale        uint32
	Duration         uint32
}

func (m MovieHeaderv0) Print() {
	fmt.Println("Movie header version 0")
	fmt.Printf("Creation time: %d\n", m.CreationTime)
	fmt.Printf("Modification time: %d\n", m.ModificationTime)
	fmt.Printf("Timescale: %d\n", m.Timescale)
	fmt.Printf("Duration: %d\n", m.Duration)
	fmt.Printf("Duration in seconds: %.2f\n", m.GetDurationInSeconds())
}

func (m MovieHeaderv0) GetDurationInSeconds() float64 {
	return float64(m.Duration) / float64(m.Timescale)
}

type MovieHeaderv1 struct {
	CreationTime     uint64
	ModificationTime uint64
	Timescale        uint32
	Duration         uint64
}

func (m MovieHeaderv1) Print() {
	fmt.Println("Movie header version 1")
	fmt.Printf("Creation time: %d\n", m.CreationTime)
	fmt.Printf("Modification time: %d\n", m.ModificationTime)
	fmt.Printf("Timescale: %d\n", m.Timescale)
	fmt.Printf("Duration: %d\n", m.Duration)
	fmt.Printf("Duration in seconds: %.2f\n", m.GetDurationInSeconds())
}

func (m MovieHeaderv1) GetDurationInSeconds() float64 {
	return float64(m.Duration) / float64(m.Timescale)
}

func (b BoxHeader) BoxType() string {
	return string(b.Type[:])
}

type TrackHeaderv0 struct {
	CreationTime     uint32
	ModificationTime uint32
	TrackId          uint32
	Reserved1        [4]byte
	Duration         uint32
	Reserved2        [8]byte
	Layer            int16
	AlternateGroup   int16
	Volume           int16
	Reserved3        [2]byte
	Matrix           [36]byte
	Width            uint32
	Height           uint32
}

func (t TrackHeaderv0) Print() {
	fmt.Printf("Track Id: %d\n", t.TrackId)
	actualWidth := float64(t.Width) / (1 << 16)
	actualHeight := float64(t.Height) / (1 << 16)
	fmt.Printf("Width: %.2f\n", actualWidth)
	fmt.Printf("Height: %.2f\n", actualHeight)
}

type TrackHeaderv1 struct {
	CreationTime     uint64
	ModificationTime uint64
	TrackId          uint32
	Reserved1        [4]byte
	Duration         uint64
	Reserved2        [8]byte
	Layer            int16
	AlternateGroup   int16
	Volume           int16
	Reserved3        [2]byte
	Matrix           [36]byte
	Width            uint32
	Height           uint32
}

func (t TrackHeaderv1) Print() {
	fmt.Printf("Track Id: %d", t.TrackId)
	actualWidth := float64(t.Width) / (1 << 16)
	actualHeight := float64(t.Height) / (1 << 16)
	fmt.Printf("Width: %.2f\n", actualWidth)
	fmt.Printf("Height: %.2f\n", actualHeight)
}

type HandlerReference struct {
	Predef      [4]byte
	HandlerType [4]byte
}

func (h HandlerReference) GetHandlerType() string {
	return string(h.HandlerType[:])
}

func (h HandlerReference) Print() {
	fmt.Printf("The track is of type: %s\n", h.GetHandlerType())
}

type MediaHeaderv0 struct{
	CreationTime uint32
	ModificationTime uint32
	Timescale uint32
	Duration uint32
	Language uint16
	Pre_defined uint16
}

func (m MediaHeaderv0) LanguageCode() string {
    c1 := byte((m.Language>>10)&0x1F) + 0x60
    c2 := byte((m.Language>>5)&0x1F) + 0x60
    c3 := byte(m.Language&0x1F) + 0x60

    return string([]byte{c1, c2, c3})
}

func (m MediaHeaderv0) Print(){
	fmt.Printf("The track duration in second: %.2f\n",float64(m.Duration)/float64(m.Timescale))
	fmt.Printf("The language code is: %s\n",m.LanguageCode())
}

type MediaHeaderv1 struct{
	CreationTime uint64
	ModificationTime uint64
	Timescale uint32
	Duration uint64
	Language uint16
	Pre_defined uint16
}

func (m MediaHeaderv1) LanguageCode() string {
    c1 := byte((m.Language>>10)&0x1F) + 0x60
    c2 := byte((m.Language>>5)&0x1F) + 0x60
    c3 := byte(m.Language&0x1F) + 0x60

    return string([]byte{c1, c2, c3})
}

func (m MediaHeaderv1) Print(){
	fmt.Printf("The track duration in second: %.2f\n",float64(m.Duration)/float64(m.Timescale))
	fmt.Printf("The language code is: %s\n",m.LanguageCode())
}
