@echo off
cd pack

:: Add icon to windows executable
echo = Adding icon...

    for %%f in (ximfect-*.exe) do (
        ResourceHacker -open %%f -save %%f -action addskip -res ..\img\ximfect.ico -mask ICONGROUP,MAINICON, -log NUL
        echo -- Added icon to %%f.
    )
    
echo --- Done!
echo.

:: Package the effects into a .zip file using 7z
echo = Packaging effects...

    if exist effects.zip del effects.zip
    cd effects
    7z a ..\effects.zip * -bso0
    cd ..
    
echo --- Done!
echo.

:: Move all release files into the `release` directory
echo = Moving into release directory

    move effects.zip release > NUL
    echo Moved effects.
    move ximfect-* release > NUL
    echo Moved executables.
    
echo --- Done!
echo.

cd ..