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

	if _, ver := args["version"]; ver {
		fmt.Println(tool.Version)
		return
	}

	fmt.Println("ximfect v" + tool.Version)
	fmt.Println("Learn more at https://github.com/QeaML/ximfect")
	fmt.Println("")

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
			name := fx.Metadata.Name
			version := fx.Metadata.Version
			author := fx.Metadata.Author
			desc := fx.Metadata.Desc
			fmt.Printf("======== About %s ========\n", eff.Value)
			fmt.Printf("Name:           %s\n", name)
			fmt.Printf("Version:        %s\n", version)
			fmt.Printf("Author:         %s\n", author)
			fmt.Printf("Description:    %s\n", desc)
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
				if err != nil {
					fmt.Println(err)
					return
				}
				imgR, _, err := image.Decode(file)
				if err != nil {
					fmt.Println(err)
					return
				}
				if img, ok := imgR.(*image.RGBA); ok {
					err = effect.Apply(fx, img)
					if err != nil {
						fmt.Println(err)
					} else if hasOutFile {
						outFile, err := os.Create(outFilename.Value)
						if err != nil {
							fmt.Println(err)
							return
						}
						err = png.Encode(outFile, img)
						if err != nil {
							fmt.Println(err)
							return
						}
					}
				}
			}
		}
	}
}
