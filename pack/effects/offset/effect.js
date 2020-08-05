function effect(x, y, pixel) {
    x += offset;
    y += offset;
    if(x >= size.x + jump.x) {
        x -= size.x;
    }
    if(y >= size.y + jump.y) {
        y -= size.y;
    }
    return ImageAt(x, y);
}