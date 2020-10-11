<h1><img src="img/ximfect.png" alt="ximfect logo" width="32px" height="32px">&nbsp;ximfect</h1>
An effect-based image processing tool.<br />


<i>
Did you know? ximfect has an <a href="https://youtu.be/PGSvlpF07tU">official theme song</a>!

Need help? Join the official <a href="https://discord.gg/AGPZyUE">Discord server</a>!
</i>

# Usage
`ximfect (action) <--namedArgument value --otherArgument other value ...>`

You can see a list of all actions by running `ximfect help` or just `ximfect` by itself.

Dragging & dropping a `.xfp` file on the ximfect executable will unpack & install the effect(s) from the file.

# Docs
Documentation is available [here](https://ximfect.github.io).

# How to install

1. Go to the [Releases](https://github.com/ximfect/ximfect/releases) page.
2. Find the release you wish to install.
3. Download the executable for your OS+architecture combo.
4. Move your executable to a favorable location. (Windows: `C:\ximfect`, Linux: `/usr/bin`)
5. Rename the executable to be just `ximfect`. (`ximfect.exe` on Windows)
6. If you wish to install the pre-packaged effects, look for a tutorial below.
7. (Windows-only) Add `C:\ximfect` to you PATH.
8. **Done!**

# How to install effects

## Release effects

1. Go to the [Releases](https://github.com/ximfect/ximfect/releases) page.
2. Find the release you have installed. (check with `ximfect version`)
3. Download the `effects.zip` file.
4. Using your command prompt/terminal, navigate to the folder you downloaded the file.
5. Unzip the file.
6. For each effect, run `ximfect unpack --file {effect}.xfp`.
6. **Done!**

## User-made effects

1. Download the provided `.xfp` file of your effect.
2. Using your command prompt/terminal, navigate to the folder you downloaded the file.
3. Run `ximfect unpack --file {effect}.xfp`.
4. **Done!**
