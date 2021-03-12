function effect(pixel)
    if skip then return pixel end
    pixel["r"] = ramp[int((pixel["r"] / 256) * ramp_c) + 1]
    pixel["g"] = ramp[int((pixel["g"] / 256) * ramp_c) + 1]
    pixel["b"] = ramp[int((pixel["b"] / 256) * ramp_c) + 1]
    return pixel
end
