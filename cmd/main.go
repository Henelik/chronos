package main

import (
	"fmt"
	"os"

	"github.com/Henelik/chronos"
)

func main() {
	file, err := os.Open("grubhub.mp4")
	if err != nil {
		panic(err)
	}

	mp4, err := chronos.ReadMP4(file)
	if err != nil {
		panic(err)
	}

	fmt.Printf("MP4: %#v\n", mp4)
	fmt.Printf("MVHD: %#v\n", mp4.MVHD)
}
