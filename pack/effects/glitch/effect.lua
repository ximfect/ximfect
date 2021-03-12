function effect(pixel)
	outx = pixel["x"] + offsetmap[pixel.y + 1]
	outy = pixel["y"]

	if outx > imgsize["x"] or outx < 0 then
		return {r=0,g=0,b=0,a=255}
	else
		return at(outx, outy)
	end
end
