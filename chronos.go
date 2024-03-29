package chronos

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"time"
)

type MP4 struct {
	File         *os.File
	Metadata     *Metadata
	MVHDPosition int64
}

// MVHDData reflects the raw byte structure of the MP4 metadata
type MVHDData struct {
	Version          byte
	_                [3]byte
	CreationTime     int32
	ModificationTime int32
	TimeScale        uint32
	Duration         uint32
}

// Metadata is an internal representation of the MP4 metadata
type Metadata struct {
	Version          int           `json:"version"`
	CreationTime     time.Time     `json:"creationTime"`
	ModificationTime time.Time     `json:"modificationTime"`
	TimeScale        uint32        `json:"timeScale"`
	Duration         uint32        `json:"duration"`
	TimeDuration     time.Duration `json:"timeDuration,string"`
}

func ReadMP4(file *os.File) (*MP4, error) {
	pos, err := getMVHDPosition(file)

	mvhd, err := parseMVHD(file, pos)
	if err != nil {
		return nil, err
	}

	return &MP4{
		File:         file,
		Metadata:     mvhdToMetadata(mvhd),
		MVHDPosition: pos,
	}, nil
}

func getMVHDPosition(file *os.File) (int64, error) {
	pos, err := findBytes([]byte{0x6d, 0x76, 0x68, 0x64}, file)
	if err != nil {
		return 0, err
	}

	// the position is the beginning of the MVHD tag, so we need to set it to the end
	pos += 4

	return pos, nil
}

func parseMVHD(file *os.File, pos int64) (*MVHDData, error) {
	newMVHD := new(MVHDData)

	_, err := file.Seek(pos, 0)
	if err != nil {
		return nil, err
	}

	binary.Read(file, binary.BigEndian, newMVHD)

	return newMVHD, nil
}

func (mp4 *MP4) WriteMVHD() error {
	_, err := mp4.File.Seek(mp4.MVHDPosition, 0)
	if err != nil {
		return err
	}

	err = binary.Write(mp4.File, binary.BigEndian, metadataToMVHD(mp4.Metadata))
	if err != nil {
		return err
	}

	return nil
}

func findBytes(key []byte, file *os.File) (int64, error) {
	if len(key) < 1 {
		return 0, errors.New("findBytes: key must contain at least 1 value")
	}

	info, err := file.Stat()
	if err != nil {
		return 0, err
	}

	bs := make([]byte, len(key))
	for i := int64(0); i < info.Size(); i++ {
		_, err := file.ReadAt(bs, i)
		if err != nil {
			return 0, fmt.Errorf("findBytes: %w", err)
		}
		if bytes.Equal(bs, key) {
			return i, nil
		}
	}
	return 0, errors.New("findBytes: match not found")
}

func mvhdToMetadata(mvhd *MVHDData) *Metadata {
	return &Metadata{
		Version:          int(mvhd.Version),
		CreationTime:     time.Unix(int64(mvhd.CreationTime), 0),
		ModificationTime: time.Unix(int64(mvhd.ModificationTime), 0),
		TimeScale:        mvhd.TimeScale,
		Duration:         mvhd.Duration,
		TimeDuration:     time.Duration(mvhd.Duration) * time.Second / time.Duration(mvhd.TimeScale),
	}
}

func metadataToMVHD(metadata *Metadata) *MVHDData {
	return &MVHDData{
		Version:          byte(metadata.Version),
		CreationTime:     int32(metadata.CreationTime.Unix()),
		ModificationTime: int32(metadata.ModificationTime.Unix()),
		TimeScale:        metadata.TimeScale,
		Duration:         metadata.Duration,
	}
}
