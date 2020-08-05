package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"strings"
	"ximfect/effect"
	"ximfect/tool"
)

func main() {
	args := tool.GetArgv(os.Args)

	if _, silent := args.NamedArgs["silent"]; !silent {
		fmt.Println("ximfect v" + tool.Version)
		fmt.Println("Learn more at https://github.com/QeaML/ximfect")
		fmt.Println("")
	}

	if len(args.PosArgs) == 0 {
		fmt.Println("Not enough positional arguments!")
		return
	}

	action := args.PosArgs[0].Value

	eff, hasEffect := args.NamedArgs["effect"]
	filename, hasFile := args.NamedArgs["file"]
	outFilename, hasOutFile := args.NamedArgs["out"]

	switch action {
	case "about":
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
			preload := fx.Metadata.Preload
			fmt.Printf("======== About %s ========\n", eff.Value)
			fmt.Printf("Name:           %s\n", name)
			fmt.Printf("Version:        %s\n", version)
			fmt.Printf("Author:         %s\n", author)
			fmt.Printf("Description:    %s\n", desc)
			if len(preload) > 0 {
				fmt.Printf("Preload:         %v\n", strings.Join(preload, ", "))
			}
		} else {
			fmt.Println("Please specify an effect with --effect <id>")
		}
	case "apply":
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
					} else {
						fmt.Println("Please specify an output file with --out <filename>")
					}
				}
			} else {
				fmt.Println("Please specify an input file with --file <filename>")
			}
		} else {
			fmt.Println("Please specify an effect with --effect <id>")
		}
	case "version":
		fmt.Println(tool.Version)
	}
}
