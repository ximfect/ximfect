function effect(pixel)
    local x = pixel["x"] + offset_x
    local y = pixel["y"] + offset_y

    while x >= img_size["x"] + jump_x do
        x = x - img_size["x"]
    end
    while y >= img_size["x"] + jump_x do
        y = y - img_size["x"]
    end

    return at(x, y)
end