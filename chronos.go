package chronos

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

type MP4 struct {
	Filename string
	File     *os.File
	MVHD     *MVHD
}

type MVHD struct {
	Version          byte
	_                [3]byte // ignore flags
	CreationTime     int32
	ModificationTime int32
	TimeScale        uint32
	Duration         uint32
}

func ReadMP4(filename string) (*MP4, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	mvhd, err := parseMVHD(file)
	if err != nil {
		return nil, err
	}

	return &MP4{
		Filename: filename,
		File:     file,
		MVHD:     mvhd,
	}, nil
}

func parseMVHD(file *os.File) (*MVHD, error) {
	// first, find the mvhd in the header
	pos, err := findBytes([]byte{0x6d, 0x76, 0x68, 0x64}, file)
	if err != nil {
		return nil, err
	}

	fmt.Printf("pos: %v\n", pos)

	// the position is the beginning of the MVHD tag, so we need to set it to the end
	pos += 4

	newMVHD := new(MVHD)

	_, err = file.Seek(pos, 0)
	if err != nil {
		return nil, err
	}

	binary.Read(file, binary.BigEndian, newMVHD)

	return newMVHD, nil
}

func (m *MP4) TimeStep() int {
	return 0
}

func (m *MP4) Duration() int {
	return 0
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
