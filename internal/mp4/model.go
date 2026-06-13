package mp4

type Movie struct{
	Duration float64
	Tracks []*Track
}

type Track struct{
	ID uint32
	Type string
	Duration float64
	Width float64
	Height float64

	MediaHeader *MediaHeader
	TrackHeader *TrackHeader
}

type BoxHeader struct {
	Size uint32
	Type [4]byte
}

func (b BoxHeader) BoxType() string {
	return string(b.Type[:])
}

type FullBox struct {
	Version uint8
	Flags   [3]byte
}

type MovieHeader struct {
	Version uint8
	CreationTime     uint64
	ModificationTime uint64
	Timescale uint32
	Duration  uint64
}

func (m MovieHeader) DurationSeconds() float64 {
	return float64(m.Duration) / float64(m.Timescale)
}

type TrackHeader struct {
	Version uint8
	CreationTime     uint64
	ModificationTime uint64
	TrackID uint32
	Duration uint64
	Width  float64
	Height float64
}

type HandlerReference struct {
	HandlerType string
}

type MediaHeader struct {
	Version uint8
	CreationTime     uint64
	ModificationTime uint64
	Timescale uint32
	Duration  uint64
	Language string
}

func (m MediaHeader) DurationSeconds() float64 {
	return float64(m.Duration) / float64(m.Timescale)
}

type ParseContext struct {
	Verbose bool
	Movie *Movie
	CurrentTrack *Track
}

type SampleDescription struct {
	Codec string
}

type TimeToSample struct {
}

type SampleSize struct {
}

type SampleToChunk struct {
}

type ChunkOffset struct {
}


