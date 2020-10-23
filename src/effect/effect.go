/* effect definitions */

package effect

import (
	"errors"
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"image/color"
	"ximfect/environ"

	"github.com/ximfect/ximgy"
)

// Metadata contains additional information about an Effect
type Metadata struct {
	Name    string
	Version string
	ID      string
	Author  string
	Desc    string
	Preload []string
}

// Effect represents an effect that can be applied to an Image
type Effect struct {
	Metadata *Metadata
	source   string
	vm       *lua.LState
}

// NewEffect returns an Effect constructed from the given source and metadata
func NewEffect(meta *Metadata, src string) *Effect {
	tmp := new(Effect)
	tmp.Metadata = meta
	tmp.SetSource(src)
	return tmp
}

// Load prepares the effect for running
func (e *Effect) Load(vm *lua.LState) error {
	e.vm = vm
	var err error
	if len(e.Metadata.Preload) > 0 {
		fmt.Println(" - Preloading...")
		for _, filename := range e.Metadata.Preload {
			err := e.vm.DoFile(environ.AppdataPath("effects", e.Metadata.ID, filename))
			if err != nil {
				return fmt.Errorf("error during effect preload: %v", err)
			}
		}
	}
	err = e.vm.DoString(e.source)
	if err != nil {
		return fmt.Errorf("error while loading effect: %v", err)
	}
	return nil
}

// Run processes the given image on the given VM
func (e *Effect) Run(pixel ximgy.Pixel) (color.RGBA, error) {
	def := color.RGBA{}
	inTable := e.vm.CreateTable(6, 1)
	inTable.RawSet(lua.LString("r"), lua.LNumber(pixel.R))
	inTable.RawSet(lua.LString("g"), lua.LNumber(pixel.G))
	inTable.RawSet(lua.LString("b"), lua.LNumber(pixel.B))
	inTable.RawSet(lua.LString("a"), lua.LNumber(pixel.A))
	inTable.RawSet(lua.LString("x"), lua.LNumber(pixel.X))
	inTable.RawSet(lua.LString("y"), lua.LNumber(pixel.Y))
	err := e.vm.CallByParam(lua.P{
		Fn:      e.vm.GetGlobal("effect"),
		NRet:    1,
		Protect: true,
	}, inTable)
	if err != nil {
		return def, err
	}
	ret := e.vm.Get(-1)
	e.vm.Pop(1)
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

}

// SetSource sets the source for the effect
func (e *Effect) SetSource(src string) {
	e.source = src
}
