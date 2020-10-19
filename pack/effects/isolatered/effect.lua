function effect(pixel)
    local avg = (pixel["r"] + pixel["g"] + pixel["b"]) / 3
    local clr = avg + ((pixel["r"] / 128) * (pixel["r"] - avg))

    if avg > 255 then avg = 255 end
    if clr > 255 then clr = 255 end

    return { r = clr, g = avg, b = avg, a = pixel["a"] }
end