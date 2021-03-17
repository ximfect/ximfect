/* effect definitions */

package vm

import (
	"errors"
	"image/color"
	"ximfect/tool"

	lua "github.com/yuin/gopher-lua"

	"github.com/ximfect/ximgy"
)

// EffectMetadata contains additional information about an Effect
type EffectMetadata struct {
	Name    string
	Version string
	ID      string
	Author  string
	Desc    string
	Preload []string
}

// Effect represents an effect that can be applied to an Image
type Effect struct {
	Metadata *EffectMetadata
	Dir      string
}

// NewEffect returns an Effect constructed from the given dir and metadata
func NewEffect(meta *EffectMetadata, dir string) *Effect {
	return &Effect{meta, dir}
}

func (e *Effect) run(img *ximgy.Image, vm *lua.LState, mat [][]lua.LValue) error {
	// create an input table (reused for every iteration)
	inTable := vm.CreateTable(6, 1)
	// iterate through the image
	for x := 0; x < img.Size.X; x++ {
		for y := 0; y < img.Size.Y; y++ {
			// get the pixel
			pixel := img.At(x, y)

			// set the values in the input table
			inTable.RawSetString("r", lua.LNumber(pixel.R))
			inTable.RawSetString("g", lua.LNumber(pixel.G))
			inTable.RawSetString("b", lua.LNumber(pixel.B))
			inTable.RawSetString("a", lua.LNumber(pixel.A))
			inTable.RawSetString("x", lua.LNumber(x))
			inTable.RawSetString("y", lua.LNumber(y))

			// call the effect() function
			err := vm.CallByParam(lua.P{
				Fn:      vm.GetGlobal("effect"),
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
func (e *Effect) Apply(img *ximgy.Image, ctx *tool.Context) error {
	log := ctx.Log.Sub("Effect")

	// create a vm state
	log.Debug("Creating VM state...")
	vm, err := e.vm(img, ctx)
	if err != nil {
		return err
	}

	// run the effect
	log.Debug("Running effect...")
	// create a matrix to store the results
	matrix := [][]lua.LValue{}
	// fill it with nil
	for x := 0; x < img.Size.X; x++ {
		matrix = append(matrix, []lua.LValue{})
		for y := 0; y < img.Size.Y; y++ {
			matrix[x] = append(matrix[x], nil)
		}
	}
	// run the effect on the matrix
	err = e.run(img, vm, matrix)
	if err != nil {
		return err
	}

	// apply changes
	log.Debug("Applying changes...")
	for x := 0; x < img.Size.X; x++ {
		for y := 0; y < img.Size.Y; y++ {
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
				return errors.New("return value is not a table")
			}

			// get raw r, g, b, a objects and cast them to numbers
			if rRaw, ok = tbl.RawGetString("r").(lua.LNumber); !ok {
				return errors.New("red value is not a number")
			}
			if gRaw, ok = tbl.RawGetString("g").(lua.LNumber); !ok {
				return errors.New("green value is not a number")
			}
			if bRaw, ok = tbl.RawGetString("b").(lua.LNumber); !ok {
				return errors.New("blue value is not a number")
			}
			if aRaw, ok = tbl.RawGetString("a").(lua.LNumber); !ok {
				return errors.New("alpha value is not a number")
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
	return nil
}
