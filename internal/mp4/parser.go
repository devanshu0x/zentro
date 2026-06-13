package mp4

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func ReadBlock(r io.Reader, data any) error {
	return binary.Read(r, binary.BigEndian, data)
}

func isContainer(boxType string) bool {
	switch boxType {
	case "moov", "trak", "mdia", "minf", "stbl":
		return true
	}
	return false
}

func ParseFile(path string, verbose bool) (*Movie, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	movie := &Movie{}

	ctx := &ParseContext{
		Verbose: verbose,
		Movie:   movie,
	}

	if err := ParseBoxes(
		file,
		info.Size(),
		0,
		ctx,
	); err != nil {
		return nil, err
	}

	return movie, nil
}

func ParseBoxes(file io.ReadSeeker, end int64, level int, ctx *ParseContext) error {

	for {
		curr, err := file.Seek(0, io.SeekCurrent)
		if err != nil {
			return err
		}

		if curr >= end {
			break
		}

		var header BoxHeader

		if err := ReadBlock(file, &header); err != nil {
			return err
		}

		if header.Size < 8 {
			return fmt.Errorf(
				"invalid box size %d for box %s",
				header.Size,
				header.BoxType(),
			)
		}

		if ctx.Verbose{
			PrintLevel(header,level)
		}

		boxEnd := curr + int64(header.Size)

		if isContainer(header.BoxType()) {

			if header.BoxType() == "trak" {

				oldTrack := ctx.CurrentTrack

				track := &Track{}

				ctx.Movie.Tracks = append(
					ctx.Movie.Tracks,
					track,
				)

				ctx.CurrentTrack = track

				err := ParseBoxes(
					file,
					boxEnd,
					level+1,
					ctx,
				)

				ctx.CurrentTrack = oldTrack

				if err != nil {
					return err
				}

			} else {

				if err := ParseBoxes(
					file,
					boxEnd,
					level+1,
					ctx,
				); err != nil {
					return err
				}
			}

			_, err := file.Seek(boxEnd, io.SeekStart)
			if err != nil {
				return err
			}

			continue
		}

		switch header.BoxType() {

		case "mvhd":

			mvhd, err := ParseMvhd(file)
			if err != nil {
				return err
			}
			
			if ctx.Verbose{
				PrintLevel(mvhd,level)
			}

			ctx.Movie.Duration =
				mvhd.DurationSeconds()

		case "tkhd":

			tkhd, err := ParseTkhd(file)
			if err != nil {
				return err
			}

			if ctx.Verbose{
				PrintLevel(tkhd,level)
			}

			if ctx.CurrentTrack != nil {

				ctx.CurrentTrack.ID =
					tkhd.TrackID

				ctx.CurrentTrack.Width =
					tkhd.Width

				ctx.CurrentTrack.Height =
					tkhd.Height

				ctx.CurrentTrack.TrackHeader =
					tkhd
			}

		case "mdhd":

			mdhd, err := ParseMdhd(file)
			if err != nil {
				return err
			}

			if ctx.Verbose{
				PrintLevel(mdhd,level)
			}

			if ctx.CurrentTrack != nil {

				ctx.CurrentTrack.Duration =
					mdhd.DurationSeconds()

				ctx.CurrentTrack.MediaHeader =
					mdhd
			}

		case "hdlr":

			hdlr, err := ParseHdlr(file)
			if err != nil {
				return err
			}

			if ctx.Verbose{
				PrintLevel(hdlr,level)
			}

			if ctx.CurrentTrack != nil {
				ctx.CurrentTrack.Type =
					hdlr.HandlerType
			}

		default:

			_, err := file.Seek(
				int64(header.Size)-8,
				io.SeekCurrent,
			)

			if err != nil {
				return err
			}

			continue
		}

		_, err = file.Seek(
			boxEnd,
			io.SeekStart,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
