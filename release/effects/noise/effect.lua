function effect(pixel)
	noiseamt = (random() * 128) - 64
	finalr = pixel["r"] + noiseamt
	finalg = pixel["g"] + noiseamt
	finalb = pixel["b"] + noiseamt
	if finalr > 255 then
		finalr = 255
	end
	if finalr < 0 then
		finalr = 0
	end
	if finalg > 255 then
		finalg = 255
	end
	if finalg < 0 then
		finalg = 0
	end
	if finalb > 255 then
		finalb = 255
	end
	if finalb < 0 then
		finalb = 0
	end
	return {r=finalr, g=finalg, b=finalb, a=255}
end