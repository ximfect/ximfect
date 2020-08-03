package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"ximfect/effect"
	"ximfect/tool"
)

func main() {
	args := tool.GetArgv()

	eff, hasEffect := args["effect"]
	filename, hasFile := args["file"]
	outFilename, hasOutFile := args["out"]

	if _, about := args["about"]; about {
		if hasEffect {
			fx, err := effect.LoadFromAppdata(eff.Value)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(fx.Metadata)
		}
	}
	if _, apply := args["apply"]; apply {
		if hasEffect {
			fx, err := effect.LoadFromAppdata(eff.Value)
			if err != nil {
				fmt.Println(err)
				return
			}
			if hasFile {
				file, err := os.Open(filename.Value)
				imgR, _, err := image.Decode(file)
				if img, ok := imgR.(*image.RGBA); ok {
					err = effect.Apply(fx, img)
					if err != nil {
						fmt.Println(err)
					} else if hasOutFile {
						outFile, _ := os.Create(outFilename.Value)
						png.Encode(outFile, img)
					}
				}
			}
		}
	}
}
