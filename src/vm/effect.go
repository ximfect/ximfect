/* effect definitions */

package vm

import (
	"errors"
	lua "github.com/yuin/gopher-lua"
	"image/color"
	"ximfect/tool"

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

func (e *Effect) run(img *ximgy.Image, vm *lua.LState, matrix [][][]lua.LValue) error {
	inTable := vm.CreateTable(6, 1)
	for x := 0; x < img.Size.X; x++ {
		for y := 0; y < img.Size.Y; y++ {
			pixel := img.At(x, y)
			inTable.RawSetString("r", lua.LNumber(pixel.R))
			inTable.RawSetString("g", lua.LNumber(pixel.G))
			inTable.RawSetString("b", lua.LNumber(pixel.B))
			inTable.RawSetString("a", lua.LNumber(pixel.A))
			inTable.RawSetString("x", lua.LNumber(x))
			inTable.RawSetString("y", lua.LNumber(y))
			err := vm.CallByParam(lua.P{
				Fn:      vm.GetGlobal("effect"),
				NRet:    1,
				Protect: true,
			}, inTable)
			if err != nil {
				return err
			}
			ret := vm.Get(-1)
			vm.Pop(1)
			tbl, ok := ret.(*lua.LTable)
			if !ok {
				return errors.New("could not convert return value to table")
			}
			matrix[x][y][0] = tbl.RawGetString("r")
			matrix[x][y][1] = tbl.RawGetString("g")
			matrix[x][y][2] = tbl.RawGetString("b")
			matrix[x][y][3] = tbl.RawGetString("a")
		}
	}
	return nil
}

// Apply applies this effect to the given Image
func (e *Effect) Apply(img *ximgy.Image, ctx *tool.Context) error {
	log := ctx.Log.Sub("Effect")

	log.Debug("Creating VM state...")
	vm, err := e.vm(img, ctx)
	if err != nil {
		return err
	}

	log.Debug("Running effect...")
	matrix := [][][]lua.LValue{}
	for x := 0; x < img.Size.X; x++ {
		matrix = append(matrix, [][]lua.LValue{})
		for y := 0; y < img.Size.Y; y++ {
			matrix[x] = append(matrix[x], []lua.LValue{nil, nil, nil, nil})
		}
	}
	err = e.run(img, vm, matrix)
	if err != nil {
		return err
	}

	log.Debug("Applying changes...")
	for x := 0; x < img.Size.X; x++ {
		for y := 0; y < img.Size.Y; y++ {
			rRaw := matrix[x][y][0]
			gRaw := matrix[x][y][1]
			bRaw := matrix[x][y][2]
			aRaw := matrix[x][y][3]
			if rRaw.Type() != lua.LTNumber {
				return errors.New("red value is not a number")
			}
			if gRaw.Type() != lua.LTNumber {
				return errors.New("green value is not a number")
			}
			if bRaw.Type() != lua.LTNumber {
				return errors.New("blue value is not a number")
			}
			if aRaw.Type() != lua.LTNumber {
				return errors.New("alpha value is not a number")
			}
			r := uint8(rRaw.(lua.LNumber))
			g := uint8(gRaw.(lua.LNumber))
			b := uint8(bRaw.(lua.LNumber))
			a := uint8(aRaw.(lua.LNumber))
			img.Set(x, y, color.RGBA{r,g,b,a})
		}
	}

	log.Debug("Done!")
	return nil
}