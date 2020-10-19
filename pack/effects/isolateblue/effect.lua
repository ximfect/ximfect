function effect(pixel)
    local avg = (pixel["r"] + pixel["g"] + pixel["b"]) / 3
    local clr = avg + ((pixel["b"] / 128) * (pixel["b"] - avg))

    if avg > 255 then avg = 255 end
    if clr > 255 then clr = 255 end

    return { r = avg, g = clr, b = avg, a = pixel["a"] }
end