<h1><img src="img/ximfect.png" alt="ximfect logo" width="32px" height="32px">&nbsp;ximfect</h1>
<i>An effect-based image processing tool.</i>


*Did you know? ximfect has an [official theme song](https://youtu.be/PGSvlpF07tU)!*

# Need help?
[Join the official Discord server!](https://discord.gg/AGPZyUE)

# Usage
`ximfect (action) <--namedArgument value --otherArgument other value ...>`

You can install effects from `.zip` files using the `unpack` action.

# How to install (Windows)

## Release (stable)
1. Go to the [Releases](https://github.com/QeaML/ximfect/releases) page.
2. Download the latest release's `ximfect-v(...).zip` file.
3. Open the file, and extract the `ximfect.exe` file to a favorable location (e.g. `C:\ximfect`)
4. (If you haven't already) Add that location to your PATH.
5. Open your APPDATA folder.
6. (If it doesn't exist) Create a `ximfect` directory.
7. (If there isn't one) Extract the `effects` directory from the zip file into the directory you just created.
8. **Done!**

## Development (no guarantee)
1. Clone this repository to your computer.
2. Inside your command prompt(can't be PowerShell), navigate to the repository's directory. (called ximfect)
3. Run `pack`.
4. A release-ready zip file should've been created. Perform the standard release installation instructions, starting with step 3.
5. **Done!**

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

The easiest way to install an effect, is to store the **directory containing the effect's files** in a zip file, then dragging and dropping the zip file onto the ximfect executable. You can also distribute the zip file, others can install it in the same way.
