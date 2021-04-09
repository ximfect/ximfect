package fxchain

import (
	"fmt"
	"ximfect/tool"
	"ximfect/vm"

	"github.com/ximfect/ximgy"
)

// Apply parses an FX chain and applies it to the supplied Image.
func Apply(src string, img *ximgy.Image, ctx *tool.Context) (*ximgy.Image, error) {
	// create a sub-log
	log := ctx.Log.Sub("FXChain")

	// parse the fx chain "script"
	log.Debug("Parsing FX chain...")
	chain := ParseChain(src)
	// amount of effects in the chain
	total := len(chain)

	// apply every effect in the chain
	for i, pair := range chain {
		//logging
		log.Debug(fmt.Sprintf("[%d/%d] Effect `%s`:\n", i, total, pair.effect))

		// load effect
		log.Debug("Loading effect...")
		eff, err := vm.LoadAppdataEffect(pair.effect)
		if err != nil {
			return nil, err
		}

		// put the specified arguments into an ArgumentList
		log.Debug("Preparing arguments...")
		a := tool.ArgumentList{}
		for k, v := range pair.params {
			a.NArgs[k] = tool.Argument{IsValue: true, Value: v, BoolValue: true}
		}

		// apply the effect
		log.Debug("Applying effect...")
		err = eff.Apply(img, ctx)
		if err != nil {
			return nil, err
		}

		// set the source to the output (makes the effects stack on each other)
		log.Debug("Updating source...")
		img.SetSource(img.GetOutput())
	}
	return img, nil
}
