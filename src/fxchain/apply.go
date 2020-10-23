package fxchain

import (
	"ximfect/effect"
	"ximfect/tool"

	"github.com/ximfect/ximgy"
)

// Apply parses an FX chain and applies it to the supplied Image.
func Apply(src string, img *ximgy.Image, t *tool.Tool) (*ximgy.Image, error) {
	t.VerboseLn(" - Parsing FX chain...")
	chain := ParseChain(src)
	total := len(chain)
	for i, pair := range chain {
		t.VerboseF(" - [%d/%d] Effect `%s`:\n", i, total, pair.effect)
		t.VerboseLn(" -- Loading effect...")
		eff, err := effect.LoadFromAppdata(pair.effect)
		if err != nil {
			return nil, err
		}
		t.VerboseF(" -- Preparing arguments...")
		a := tool.ArgumentList{}
		for k, v := range pair.params {
			a.NamedArgs[k] = tool.Argument{IsValue: true, Value: v, BoolValue: true}
		}
		t.VerboseLn(" -- Applying effect...")
		err = effect.Apply(eff, img, t, a)
		if err != nil {
			return nil, err
		}
		t.VerboseLn(" -- Updating source...")
		img.SetSource(img.GetOutput())
	}
	return img, nil
}
