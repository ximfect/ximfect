<h1><img src="img/ximfect.png" alt="ximfect logo" width="32px" height="32px">&nbsp;ximfect</h1>
An effect-based image processing tool.<br />


<i>
Need help? Join the official <a href="https://discord.gg/AGPZyUE">Discord server</a>!
</i>

# Usage
`ximfect (action) <--namedArgument value --otherArgument 123 (...)>`

Run `ximfect help` or just `ximfect` by itself to see a list of actions.

Drag & drop a `.xpk` file on the ximfect executable to unpack & install the effect/lib from the file. 

*(effects come in `.fx.xpk` files, libs come in `.lib.xpk` instead)*

# The how-tos below are exclusively for the lastest bleeding-edge version of ximfect. The documentation for the latest release can be found [here](https://ximfect.github.io).

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
6. For each effect, run `ximfect unpack-effect --file {effect}.fx.xpk`.
6. **Done!**

## User-made effects

1. Download the provided `.fx.xpk` file of your effect.
2. Using your command prompt/terminal, navigate to the folder you downloaded the file.
3. Run `ximfect unpack-effect --file {effect}.fx.xpk`.
4. **Done!**

# How to install libs

1. Download the provided `.lib.xpk` file of your effect.
2. Using your command prompt/terminal, navigate to the folder you downloaded the file.
3. Run `ximfect unpack-lib --file {effect}.lib.xpk`.
4. **Done!**
