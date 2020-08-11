package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
	"ximfect/effect"
	"ximfect/environ"
	"ximfect/tool"
)

func exitWithError(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func exitWithText(text string) {
	fmt.Println(text)
	os.Exit(1)
}

func main() {
	args := tool.GetArgv(os.Args)

	if _, silent := args.NamedArgs["silent"]; !silent {
		fmt.Println("ximfect v" + tool.Version)
		fmt.Println("Learn more at https://github.com/QeaML/ximfect")
		fmt.Println("")
	}

	if len(args.PosArgs) == 0 {
		exitWithText("Not enough positional arguments!")
	}

	action := args.PosArgs[0].Value

	if strings.HasSuffix(action, ".zip") {
		err := environ.Unzip(action, environ.AppdataPath("effects"))
		if err != nil {
			exitWithError(err)
		}
	}

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
			exitWithText("Please specify an effect with --effect <id>")
		}
	case "apply":
		if hasEffect {
			fx, err := effect.LoadFromAppdata(eff.Value)
			if err != nil {
				exitWithError(err)
			}
			if hasFile {
				file, err := os.Open(filename.Value)
				if err != nil {
					exitWithError(err)
				}
				imgR, _, err := image.Decode(file)
				if err != nil {
					exitWithError(err)
				}
				if img, ok := imgR.(*image.RGBA); ok {
					err = effect.Apply(fx, img)
					if err != nil {
						exitWithError(err)
					}
					if hasOutFile {
						outFile, err := os.Create(outFilename.Value)
						if err != nil {
							exitWithError(err)
						}
						err = png.Encode(outFile, img)
						if err != nil {
							exitWithError(err)
						}
					} else {
						exitWithText("Please specify an output file with --out <filename>")
					}
				}
			} else {
				exitWithText("Please specify an input file with --file <filename>")
			}
		} else {
			exitWithText("Please specify an effect with --effect <id>")
		}
	case "version":
		fmt.Println(tool.Version)
	case "pack":
		if hasEffect {
			fmt.Println("Packing effect...")
			_, err := effect.LoadFromAppdata(eff.Value)
			if err != nil {
				exitWithError(err)
				return
			}
			err = environ.ZipIt(environ.AppdataPath("effects", eff.Value), eff.Value+".zip")
			if err != nil {
				exitWithError(err)
				return
			}
		} else {
			exitWithText("Please specify an effect with --effect <id>")
		}
	case "unpack":
		if hasFile {
			fmt.Println("Unpacking and installing effect...")
			err := environ.Unzip(filename.Value, environ.AppdataPath("effects"))
			if err != nil {
				exitWithError(err)
			}
		} else {
			exitWithText("Please specify an input file with --file <filename>")
		}
	case "test":
		if hasOutFile {
			fmt.Println("Generating test image...")
			img := image.NewRGBA(image.Rect(0, 0, 510, 510))
			var (
				x int
				y int
			)
			for y = 0; y < 510; y++ {
				for x = 0; x < 510; x++ {
					img.SetRGBA(x, y, color.RGBA{uint8(x/2 + 1), uint8(y/2 + 1), 0, 255})
				}
			}
			outFile, err := os.Create(outFilename.Value)
			if err != nil {
				exitWithError(err)
				return
			}
			err = png.Encode(outFile, img)
			if err != nil {
				exitWithError(err)
				return
			}
		} else {
			exitWithText("Please specify an output file with --out <filename>")
		}
	default:
		exitWithText("Unknown action")
	}
	os.Exit(0)
}
