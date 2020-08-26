function effect(x, y, pixel) {
	var avg = (pixel.r + pixel.g + pixel.b) / 3;
	var red = avg + ((pixel.r / 128) * (pixel.r - avg));

	if (red > 255.0) red = 255;
	if (avg > 255.0) avg = 255;

	return {r: red, g: avg, b: avg, a: pixel.a};
}
