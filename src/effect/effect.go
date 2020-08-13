/* effect definitions */

package effect

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"ximfect/environ"

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
}

// NewEffect returns an Effect constructed from the given sournce and metadata
func NewEffect(meta *Metadata, src string) *Effect {
	tmp := new(Effect)
	tmp.Metadata = meta
	tmp.SetSource(src)
	return tmp
}

// Run processes the given image on the given VM
func (e *Effect) Run(vm *otto.Otto, img *image.RGBA) error {
	size := img.Bounds().Size()
	var (
		code    string
		red32   uint32
		green32 uint32
		blue32  uint32
		alpha32 uint32
		red8    uint8
		green8  uint8
		blue8   uint8
		alpha8  uint8
		ret     otto.Value
		obj     *otto.Object
		tmp     int64
		err     error
	)
	if len(e.Metadata.Preload) > 0 {
		fmt.Println("Preloading...")
		for _, filename := range e.Metadata.Preload {
			file, err := os.Open(
				environ.AppdataPath("effects", e.Metadata.ID, filename))
			if err != nil {
				return fmt.Errorf("error during effect preload: %v", err)
			}
			_, err = vm.Run(file)
			if err != nil {
				return fmt.Errorf("error during effect preload: %v", err)
			}
		}
	}
	_, err = vm.Run(e.source)
	if err != nil {
		return fmt.Errorf("error while loading effect: %v", err)
	}
	fmt.Println("Applying effect...")
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			red32, green32, blue32, alpha32 = img.At(x, y).RGBA()
			red8 = uint8(red32 >> 8)
			green8 = uint8(green32 >> 8)
			blue8 = uint8(blue32 >> 8)
			alpha8 = uint8(alpha32 >> 8)
			code = fmt.Sprintf("effect(%d,%d,{r:%d,g:%d,b:%d,a:%d});",
				x, y, red8, green8, blue8, alpha8)
			ret, err = vm.Run(code)
			if err != nil {
				return fmt.Errorf("error while processing image: %v", err)
			}
			if !ret.IsObject() {
				return fmt.Errorf("error while processing image: function return value isn't Object")
			}
			obj = ret.Object()
			ret, err = obj.Get("r")
			if err != nil {
				return fmt.Errorf("error while processing image: %v", err)
			}
			tmp, err = ret.ToInteger()
			if err != nil {
				return fmt.Errorf("error while processing image: %v", err)
			}
			red8 = uint8(tmp)
			ret, err = obj.Get("g")
			if err != nil {
				return fmt.Errorf("error while processing image: %v", err)
			}
			tmp, err = ret.ToInteger()
			if err != nil {
				return fmt.Errorf("error while processing image: %v", err)
			}
			green8 = uint8(tmp)
			ret, err = obj.Get("b")
			if err != nil {
				return fmt.Errorf("error while processing image: %v", err)
			}
			tmp, err = ret.ToInteger()
			if err != nil {
				return fmt.Errorf("error while processing image: %v", err)
			}
			blue8 = uint8(tmp)
			ret, err = obj.Get("a")
			if err != nil {
				return fmt.Errorf("error while processing image: %v", err)
			}
			tmp, err = ret.ToInteger()
			if err != nil {
				return fmt.Errorf("error while processing image: %v", err)
			}
			alpha8 = uint8(tmp)
			img.SetRGBA(x, y, color.RGBA{red8, green8, blue8, alpha8})
		}
	}
	fmt.Println("Finished!")
	return nil
}

// SetSource sets the source for the effect
func (e *Effect) SetSource(src string) {
	e.source = src
}
