package fxchain

import (
	"ximfect/effect"
	"ximfect/tool"

	"github.com/ximfect/ximgy"
)

// Apply parses an FX chain and applies it to the supplied Image.
func Apply(src string, img *ximgy.Image, t *tool.Tool) (*ximgy.Image, error) {
	chain := ParseChain(src)
	total := len(chain)
	for i, pair := range chain {
		t.VerboseF("== [%d/%d] Applying effect `%s`:\n", i, total, pair.effect)
		eff, err := effect.LoadFromAppdata(pair.effect)
		if err != nil {
			return nil, err
		}
		a := tool.ArgumentList{}
		for k, v := range pair.params {
			a.NamedArgs[k] = tool.Argument{IsValue: true, Value: v, BoolValue: true}
		}
		err = effect.Apply(eff, img, t, a)
		if err != nil {
			return nil, err
		}
		img.SetSource(img.GetOutput())
	}
	t.VerboseLn("Chain finished.")
	return img, nil
}
