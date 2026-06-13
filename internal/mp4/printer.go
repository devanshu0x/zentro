package mp4

import (
	"fmt"
	"strings"
)

func (m Movie) String() string {
	var sb strings.Builder

	sb.WriteString(
		fmt.Sprintf(
			"Movie Duration: %.2f sec\n",
			m.Duration,
		),
	)

	for _, track := range m.Tracks {
		sb.WriteString(track.String())
		sb.WriteString("\n")
	}

	return sb.String()
}

func (t Track) String() string {
	return fmt.Sprintf(
		"Track ID=%d Type=%s Duration=%.2f Width=%.0f Height=%.0f",
		t.ID,
		t.Type,
		t.Duration,
		t.Width,
		t.Height,
	)
}

func (m MovieHeader) String() string {
	return fmt.Sprintf(
		"MovieHeader{Version=%d Timescale=%d Duration=%d DurationSeconds=%.2f}",
		m.Version,
		m.Timescale,
		m.Duration,
		m.DurationSeconds(),
	)
}

func (t TrackHeader) String() string {
	return fmt.Sprintf(
		"TrackHeader{TrackID=%d Duration=%d Width=%.0f Height=%.0f}",
		t.TrackID,
		t.Duration,
		t.Width,
		t.Height,
	)
}

func (h HandlerReference) String() string {
	return fmt.Sprintf(
		"HandlerReference{Type=%s}",
		h.HandlerType,
	)
}

func (m MediaHeader) String() string {
	return fmt.Sprintf(
		"MediaHeader{Language=%s Timescale=%d Duration=%d DurationSeconds=%.2f}",
		m.Language,
		m.Timescale,
		m.Duration,
		m.DurationSeconds(),
	)
}

func (b BoxHeader) String() string{
	return fmt.Sprintf(
		"Box: {Type=%s, size=%d}",
		b.BoxType(),
		b.Size,
	)
}

type Printable interface {
	String() string
}

func PrintLevel(p Printable,level int){
	fmt.Printf("%s%v\n",strings.Repeat(" ",level*4),p)
}