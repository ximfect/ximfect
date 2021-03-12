-- GLOBALS --
-- offset map
offsetmap = {}
-- image size (for semantics)
imgsize = size()
-- offset negation mod and step
mod = imgsize["y"] / 20
step = mod / 2

-- offset randomisation
for i=0,imgsize["y"],1 do
    local offset = 0

    -- 1/2 chance to get a random offset
    if random() > 0.5 then
        -- random offset between 0 and a fraction of the image size
        offset = randint(imgsize["y"] / 20)
    end

    -- negate offset
	if i % mod < step then
		offset = offset - (offset + offset)
    end

    -- save to offset map
    offsetmap[i + 1] = offset
end
