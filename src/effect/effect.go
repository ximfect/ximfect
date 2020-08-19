/* effect definitions */

package effect

import (
	"fmt"
	"image/color"
	"os"
	"ximfect/environ"

	"github.com/ximfect/ximgy"

	"github.com/robertkrimen/otto"
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
	vm       *otto.Otto
}

// NewEffect returns an Effect constructed from the given sournce and metadata
func NewEffect(meta *Metadata, src string) *Effect {
	tmp := new(Effect)
	tmp.Metadata = meta
	tmp.SetSource(src)
	return tmp
}

// Load prepares the effect for running
func (e *Effect) Load(vm *otto.Otto) error {
	e.vm = vm
	var err error
	if len(e.Metadata.Preload) > 0 {
		fmt.Println("Preloading...")
		for _, filename := range e.Metadata.Preload {
			file, err := os.Open(
				environ.AppdataPath("effects", e.Metadata.ID, filename))
			if err != nil {
				return fmt.Errorf("error during effect preload: %v", err)
			}
			_, err = e.vm.Run(file)
			if err != nil {
				return fmt.Errorf("error during effect preload: %v", err)
			}
		}
	}
	_, err = e.vm.Run(e.source)
	if err != nil {
		return fmt.Errorf("error while loading effect: %v", err)
	}
	return nil
}

// Run processes the given image on the given VM
func (e *Effect) Run(pixel ximgy.Pixel) color.RGBA {
	def := color.RGBA{0, 0, 0, 0}
	var (
		ret otto.Value
		obj *otto.Object
		tmp int64
		err error
	)
	code := fmt.Sprintf("effect(%d,%d,{r:%d,g:%d,b:%d,a:%d});",
		pixel.X, pixel.Y, pixel.R, pixel.G, pixel.B, pixel.A)
	ret, err = e.vm.Run(code)
	if err != nil {
		return def
	}
	if !ret.IsObject() {
		return def
	}
	obj = ret.Object()
	ret, err = obj.Get("r")
	if err != nil {
		return def
	}
	tmp, err = ret.ToInteger()
	if err != nil {
		return def
	}
	red8 := uint8(tmp)
	ret, err = obj.Get("g")
	if err != nil {
		return def
	}
	tmp, err = ret.ToInteger()
	if err != nil {
		return def
	}
	green8 := uint8(tmp)
	ret, err = obj.Get("b")
	if err != nil {
		return def
	}
	tmp, err = ret.ToInteger()
	if err != nil {
		return def
	}
	blue8 := uint8(tmp)
	ret, err = obj.Get("a")
	if err != nil {
		return def
	}
	tmp, err = ret.ToInteger()
	if err != nil {
		return def
	}
	alpha8 := uint8(tmp)
	return color.RGBA{red8, green8, blue8, alpha8}
}

// SetSource sets the source for the effect
func (e *Effect) SetSource(src string) {
	e.source = src
}
