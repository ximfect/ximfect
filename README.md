<h1><img src="img/ximfect.png" alt="ximfect logo" width="32px" height="32px">&nbsp;ximfect</h1>
<i>An effect-based image processing tool.</i>


*Did you know? ximfect has an [official theme song](https://youtu.be/PGSvlpF07tU)!*

*Need help? Join the official [Discord server](https://discord.gg/AGPZyUE)*!

# note
the below instructions are for the latest development version.

# Usage
`ximfect (action) <--namedArgument value --otherArgument other value ...>`

You can install effects from `.zip` files using the `unpack` action.

# How to install

1. Go to the [Releases](https://github.com/QeaML/ximfect/releases) page.
2. Find the release you wish to install.
3. Download the executable for your OS+architecture combo.
4. Move your executable to a favorable location. (Windows: `C:\ximfect`, Linux: `/usr/bin`)
5. Rename the executable to be just `ximfect`. (`ximfect.exe` on Windows)
6. If you wish to install the pre-packaged effects, look for a tutorial below.
7. (Windows-only) Add `C:\ximfect` to you PATH.
8. **Done!**

# How to install effects

## Release effects

1. Go to the [Releases](https://github.com/QeaML/ximfect/releases) page.
2. Find the release you have installed. (check with `ximfect version`)
3. Download the `effects.zip` file.
4. Using your command prompt/terminal, navigate to the folder you downloaded the file.
5. Run `ximfect unpack --file effects.zip`.
6. **Done!**

## User-made effects

1. Download the provided `.zip` file of your effect(s).
2. Using your command prompt/terminal, navigate to the folder you downloaded the file.
3. Run `ximfect unpack --file (effect(s) .zip file)`.
4. **Done!**

# How to create effects
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

You can distribute the effect you made with `ximfect pack`. Simply running `ximfect pack --effect (your effect's id)` will make ximfect drop a `.zip` file in the folder you ran it, which contains the effect in a distributable form. Above is a tutorial on installing effects.
