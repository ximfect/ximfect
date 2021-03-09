package fxchain

import (
	"fmt"
	"ximfect/vm"
	"ximfect/tool"

	"github.com/ximfect/ximgy"
)

// Apply parses an FX chain and applies it to the supplied Image.
func Apply(src string, img *ximgy.Image, ctx *tool.Context) (*ximgy.Image, error) {
	log := ctx.Log.Sub("FXChain")
	log.Debug("Parsing FX chain...")
	chain := ParseChain(src)
	total := len(chain)
	for i, pair := range chain {
		log.Debug(fmt.Sprintf("[%d/%d] Effect `%s`:\n", i, total, pair.effect))
		log.Debug("Loading effect...")
		eff, err := vm.LoadAppdataEffect(pair.effect)
		if err != nil {
			return nil, err
		}
		log.Debug("Preparing arguments...")
		a := tool.ArgumentList{}
		for k, v := range pair.params {
			a.NamedArgs[k] = tool.Argument{IsValue: true, Value: v, BoolValue: true}
		}
		log.Debug("Applying effect...")
		err = eff.Apply(img, ctx)
		if err != nil {
			return nil, err
		}
		log.Debug("Updating source...")
		img.SetSource(img.GetOutput())
	}
	return img, nil
}
