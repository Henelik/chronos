package main

import (
	"fmt"

	"github.com/Henelik/chronos"
)

func main() {
	mp4, err := chronos.ReadMP4("grubhub.mp4")
	if err != nil {
		panic(err)
	}

	fmt.Printf("MP4: %#v\n", mp4)
	fmt.Printf("MVHD: %#v\n", mp4.MVHD)
}
