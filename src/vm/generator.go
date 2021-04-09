/* generator definitions */

package vm

import (
	"errors"
	"image"
	"image/color"
	"ximfect/tool"

	lua "github.com/yuin/gopher-lua"

	"github.com/ximfect/ximgy"
)

// GeneratorMetadata contains additional information about a Generator
type GeneratorMetadata struct {
	Name    string
	Version string
	ID      string
	Author  string
	Desc    string
	Preload []string
}

// Generator represents an Generator that can be create a new image
type Generator struct {
	Metadata *GeneratorMetadata
	Dir      string
}

// NewGenerator returns a Generator constructed from the given dir and metadata
func NewGenerator(meta *GeneratorMetadata, dir string) *Generator {
	return &Generator{meta, dir}
}

func (g *Generator) run(size image.Point, vm *lua.LState, mat [][]lua.LValue) error {
	// create an input table (reused for every iteration)
	inTable := vm.CreateTable(6, 1)
	// iterate through the image
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			inTable.RawSetString("x", lua.LNumber(x))
			inTable.RawSetString("y", lua.LNumber(y))
			// call the effect() function
			err := vm.CallByParam(lua.P{
				Fn:      vm.GetGlobal("generate"),
				NRet:    1,
				Protect: true,
			}, inTable)
			if err != nil {
				return err
			}

			// get the return value
			ret := vm.Get(-1)
			// remove it from the stack
			vm.Pop(1)

			// put the return value into matrix
			mat[x][y] = ret
		}
	}
	return nil
}

// Apply applies this effect to the given Image
func (g *Generator) Apply(sizeX, sizeY int, ctx *tool.Context) (*ximgy.Image, error) {
	log := ctx.Log.Sub("Generator")
	size := image.Rect(0, 0, sizeX, sizeY).Size()

	// create a vm state
	log.Debug("Creating VM state...")
	vm, err := g.vm(size, ctx)
	if err != nil {
		return nil, err
	}

	// run the effect
	log.Debug("Running effect...")
	// create a matrix to store the results
	matrix := [][]lua.LValue{}
	// fill it with nil
	for x := 0; x < size.X; x++ {
		matrix = append(matrix, []lua.LValue{})
		for y := 0; y < size.Y; y++ {
			matrix[x] = append(matrix[x], nil)
		}
	}
	// run the effect on the matrix
	err = g.run(size, vm, matrix)
	if err != nil {
		return nil, err
	}

	img := ximgy.MakeEmpty(image.Rect(0, 0, sizeX, sizeY))

	// apply changes
	log.Debug("Applying changes...")
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			// get the return value for this pixel
			raw := matrix[x][y]

			// context variables
			var (
				tbl  *lua.LTable
				ok   bool
				rRaw lua.LNumber
				gRaw lua.LNumber
				bRaw lua.LNumber
				aRaw lua.LNumber
			)

			// cast return value to table
			if tbl, ok = raw.(*lua.LTable); !ok {
				return nil, errors.New("return value is not a table")
			}

			// get raw r, g, b, a objects and cast them to numbers
			if rRaw, ok = tbl.RawGetString("r").(lua.LNumber); !ok {
				return nil, errors.New("red value is not a number")
			}
			if gRaw, ok = tbl.RawGetString("g").(lua.LNumber); !ok {
				return nil, errors.New("green value is not a number")
			}
			if bRaw, ok = tbl.RawGetString("b").(lua.LNumber); !ok {
				return nil, errors.New("blue value is not a number")
			}
			if aRaw, ok = tbl.RawGetString("a").(lua.LNumber); !ok {
				return nil, errors.New("alpha value is not a number")
			}

			// cast from lua number to uint8
			r := uint8(rRaw)
			g := uint8(gRaw)
			b := uint8(bRaw)
			a := uint8(aRaw)

			// set pixel in image
			img.Set(x, y, color.RGBA{r, g, b, a})
		}
	}

	log.Debug("Done!")
	return img, nil
}
