function effect(pixel)
    local avg = (pixel["r"] + pixel["g"] + pixel["b"]) / 3
    local clr = avg + ((pixel["g"] / 128) * (pixel["g"] - avg))

    if avg > 255 then avg = 255 end
    if clr > 255 then clr = 255 end

    return { r = avg, g = avg, b = clr, a = pixel["a"] }
end