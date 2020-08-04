function effect(x, y, pixel) {
    x += offset;
    y += offset;
    if(x >= size.x) {
        x -= size.x;
    }
    if(y >= size.y) {
        y -= size.y;
    }
    return ImageAt(x, y);
}