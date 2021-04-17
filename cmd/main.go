package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Henelik/chronos"

	"github.com/urfave/cli/v2"
)

func main() {
	/* file, err := os.OpenFile("grubhub.mp4", os.O_RDWR, os.ModeExclusive)
	if err != nil {
		panic(err)
	}

	mp4, err := chronos.ReadMP4(file)
	if err != nil {
		panic(err)
	}

	mp4.MVHD.Duration = 0x7ffffff // maximum

	err = mp4.WriteMVHD()
	if err != nil {
		panic(err)
	}

	fmt.Printf("MP4: %#v\n", mp4)
	fmt.Printf("MVHD: %#v\n", mp4.MVHD)*/

	app := &cli.App{
		Name:  "chronos",
		Usage: "manipulate mp4 metadata",
		Commands: []*cli.Command{
			{
				Name: "noop",
			},
			{
				Name:     "add",
				Category: "template",
			},
			{
				Name:     "remove",
				Category: "template",
			},
		},
		Action: func(c *cli.Context) error { // Print all info
			filename := c.Args().Get(0)

			fmt.Printf("Loading file %s\n", filename)

			file, err := os.OpenFile(filename, os.O_RDWR, os.ModeExclusive)
			if err != nil {
				return err
			}

			mp4, err := chronos.ReadMP4(file)
			if err != nil {
				return err
			}

			fmt.Printf("%s MVHD:\n", filename)
			fmt.Printf("Version: %v\n", mp4.Metadata.Version)
			fmt.Printf("Creation time: %v\n", mp4.Metadata.CreationTime)
			fmt.Printf("Modification time: %v\n", mp4.Metadata.ModificationTime)
			fmt.Printf("Time scale: %v units per second\n", mp4.Metadata.TimeScale)
			fmt.Printf("Duration: %v\n", mp4.Metadata.Duration)
			fmt.Printf("Video length: %v\n", mp4.Metadata.TimeDuration)

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
