bits = arg("bits")
if bits == nil then
    info("bit depth unspecified. defaulting to 1. (specify using --bits [num])")
    bits = 1
end
bits = int(bits)
skip = false
if bits >= 8 then
    warn("bit depth is 8 or more. result will be the same as original image.")
    bits = 8
    skip = true
end
if not skip then
    ramp = {}
    ramp_step = int(256 / (2 ^ bits))
    for x=0,256,ramp_step do
        table.insert(ramp, x)
    end
    ramp_c = #ramp
    ramp[ramp_c] = ramp[ramp_c] - 1
end