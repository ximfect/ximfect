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
	source   string
}

// NewEffect returns an Effect constructed from the given source and metadata
func NewEffect(meta *EffectMetadata, src string) *Effect {
	tmp := new(Effect)
	tmp.Metadata = meta
	tmp.SetSource(src)
	return tmp
}

// SetSource sets the source for the effect
func (e *Effect) SetSource(src string) {
	e.source = src
}

func (e *Effect) run(img *ximgy.Image, vm *lua.LState) func(pixel ximgy.Pixel) (color.RGBA, error) {
	return (func(pixel ximgy.Pixel) (color.RGBA, error) {
		def := color.RGBA{}
		inTable := vm.CreateTable(6, 1)
		inTable.RawSet(lua.LString("r"), lua.LNumber(pixel.R))
		inTable.RawSet(lua.LString("g"), lua.LNumber(pixel.G))
		inTable.RawSet(lua.LString("b"), lua.LNumber(pixel.B))
		inTable.RawSet(lua.LString("a"), lua.LNumber(pixel.A))
		inTable.RawSet(lua.LString("x"), lua.LNumber(pixel.X))
		inTable.RawSet(lua.LString("y"), lua.LNumber(pixel.Y))
		err := vm.CallByParam(lua.P{
			Fn:      vm.GetGlobal("effect"),
			NRet:    1,
			Protect: true,
		}, inTable)
		if err != nil {
			return def, err
		}
		ret := vm.Get(-1)
		vm.Pop(1)
		tbl, ok := ret.(*lua.LTable)
		if !ok {
			return def, errors.New("could not convert return value to table")
		}

		redRaw := tbl.RawGet(lua.LString("r"))
		if redRaw.Type() == lua.LTNil {
			return def, errors.New("returned red value is nil")
		}
		redNum, ok := redRaw.(lua.LNumber)
		if !ok {
			return def, errors.New("returned red value is not a number")
		}
		red := uint8(redNum)

		greenRaw := tbl.RawGet(lua.LString("g"))
		if greenRaw.Type() == lua.LTNil {
			return def, errors.New("returned green value is nil")
		}
		greenNum, ok := greenRaw.(lua.LNumber)
		if !ok {
			return def, errors.New("returned green value is not a number")
		}
		green := uint8(greenNum)

		blueRaw := tbl.RawGet(lua.LString("b"))
		if blueRaw.Type() == lua.LTNil {
			return def, errors.New("returned blue value is nil")
		}
		blueNum, ok := blueRaw.(lua.LNumber)
		if !ok {
			return def, errors.New("returned blue value is not a number")
		}
		blue := uint8(blueNum)

		alphaRaw := tbl.RawGet(lua.LString("a"))
		if alphaRaw.Type() == lua.LTNil {
			return def, errors.New("returned alpha value is nil")
		}
		alphaNum, ok := alphaRaw.(lua.LNumber)
		if !ok {
			return def, errors.New("returned alpha value is not a number")
		}
		alpha := uint8(alphaNum)

		return color.RGBA{R: red, G: green, B: blue, A: alpha}, nil
	})
}

// Apply applies this effect to the given Image
func (e *Effect) Apply(img *ximgy.Image, ctx *tool.Context) error {
	log := ctx.Log.Sub("Effect")

	log.Debug("Creating VM state...")
	vm, err := e.vm(img, ctx)
	if err != nil {
		return err
	}

	log.Debug("Iterating image...")
	err = img.Iterate(e.run(img, vm))
	if err != nil {
		return err
	}

	log.Debug("Done!")
	return nil
}