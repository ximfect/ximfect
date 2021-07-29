function effect(pixel)
	neighbourhood = {}
	for i=-1,1,1 do
		for j=-1,1,1 do
			newx = pixel["x"] + i
			newy = pixel["y"] + j
			table.insert(neighbourhood, at(newx, newy))
		end
	end
	totalr = 0
	totalg = 0
	totalb = 0
	totala = 0
	amt = 0
	for i, pixel in ipairs(neighbourhood) do
		totalr = totalr + pixel["r"]
		totalg = totalg + pixel["g"]
		totalb = totalb + pixel["b"]
		totala = totala + pixel["a"]
		amt = amt + 1
	end
	return {r=totalr / amt, g=totalg / amt, b=totalb / amt, a=totala / amt}
end