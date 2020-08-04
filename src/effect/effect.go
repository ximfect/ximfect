package effect

import (
	"fmt"
	"image"
	"image/color"

	"github.com/robertkrimen/otto"
)

// Metadata contains additional information about an Effect
type Metadata struct {
	Name    string
	Version string
	ID      string
	Author  string
	Desc    string
}

// Effect represents an effect that can be applied to an Image
type Effect struct {
	Metadata *Metadata
	source   string
}

// NewEffect returns an Effect constructed from the given sournce and metadata
func NewEffect(meta *Metadata, src string) *Effect {
	tmp := new(Effect)
	tmp.Metadata = meta
	tmp.SetSource(src)
	return tmp
}

// Run runs the effect on a single pixel
func (e *Effect) Run(x, y int, vm *otto.Otto, img *image.RGBA) (*QueueEntry, error) {
	var (
		code  string
		red   uint32
		green uint32
		blue  uint32
		alpha uint32
		ret   otto.Value
		obj   *otto.Object
		tmp   int64
		err   error
	)
	_, err = vm.Run(e.source)
	if err != nil {
		return nil, fmt.Errorf("error while loading effect: %v", err)
	}

	red, green, blue, alpha = img.At(x, y).RGBA()
	code = fmt.Sprintf("effect(%d,%d,{r:%d,g:%d,b:%d,a:%d});",
		x, y, red, green, blue, alpha)
	ret, err = vm.Run(code)
	if err != nil {
		return nil, fmt.Errorf("error while processing image: %v", err)
	}
	if !ret.IsObject() {
		return nil, fmt.Errorf("error while processing image: function return value isn't Object")
	}
	obj = ret.Object()
	ret, err = obj.Get("r")
	if err != nil {
		return nil, fmt.Errorf("error while processing image: %v", err)
	}
	tmp, err = ret.ToInteger()
	if err != nil {
		return nil, fmt.Errorf("error while processing image: %v", err)
	}
	red = uint32(tmp)
	ret, err = obj.Get("g")
	if err != nil {
		return nil, fmt.Errorf("error while processing image: %v", err)
	}
	tmp, err = ret.ToInteger()
	if err != nil {
		return nil, fmt.Errorf("error while processing image: %v", err)
	}
	green = uint32(tmp)
	ret, err = obj.Get("b")
	if err != nil {
		return nil, fmt.Errorf("error while processing image: %v", err)
	}
	tmp, err = ret.ToInteger()
	if err != nil {
		return nil, fmt.Errorf("error while processing image: %v", err)
	}
	blue = uint32(tmp)
	ret, err = obj.Get("a")
	if err != nil {
		return nil, fmt.Errorf("error while processing image: %v", err)
	}
	tmp, err = ret.ToInteger()
	if err != nil {
		return nil, fmt.Errorf("error while processing image: %v", err)
	}
	alpha = uint32(tmp)
	entry := QueueEntry{
		x, y, &color.RGBA{
			uint8(red), uint8(green), uint8(blue), uint8(alpha)}}
	return &entry, nil
}

// SetSource sets the source for the effect
func (e *Effect) SetSource(src string) {
	e.source = src
}
