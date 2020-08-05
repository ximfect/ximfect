<h1><img src="img/ximfect.png" alt="ximfect logo" width="32px" height="32px">&nbsp;ximfect</h1>
<i>An effect-based image processing tool.</i>


*Did you know? ximfect has an [official theme song](https://youtu.be/PGSvlpF07tU)!*

# Need help?
[Join the official Discord server!](https://discord.gg/AGPZyUE)

# Usage
`ximfect --apply --effect <effect> --file <source filename, only supports png> --out <output filename, only supports png>`

Effects are placed in APPDATA/ximfect/effects in separate subdirectories.

# How to install (Windows)
Go to the releases page, and download the latest version. Inside the downloaded zip file, you should find a `ximfect.exe` file and a folder called `effects`. Extract `ximfect.exe` to a location with a favorable name (e.g. `C:\ximfect`) and add that location to your PATH. Next, find your way to your `%APPDATA%` folder. In there, create a folder named `ximfect` and extract the `effects` folder into the folder you just created.

# How to effects
Effects are recognized by their id, which is the name of the folder containing their files.

The brains of the effect is the `effect.js` file, which is structured like this:
```js
function effect(x, y, pixel) { // the effect's function, called on every pixel.
    /*
        `x` and `y` are coordinates of the currently processed picture. always 
        above 0 and below ImageSize()

        `pixel` is an object which contains the `r`,`g`,`b`,`a` valuse of the 
        pixel. This same structure is returned by ImageAt(x, y) and must be 
        returned by this function.
    */
   return {r: pixel.r, g: pixel.g, b: pixel.b, a: pixel.a};
}
```

And what describes the effect is the `effect.meta` file, which contains various metadata about the effect, structured like this:
```
# display(!) name of the effect
Name = nothing
# semantic version
Version = 1.0.1
# the author's name and email
Author = qeaml <qeaml@pm.me>
# a short description
Desc = Does literally nothing.
```
