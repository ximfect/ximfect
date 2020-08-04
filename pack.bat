@echo off

:: Build ximfect
cls
    echo Building...
    cd src
    if exist ximfect.exe del ximfect.exe
    go build
    move ximfect.exe .. > NUL
    cd ..
echo Done!
echo.

:: Add icon to executable
echo Adding icon...
    ResourceHacker -open ximfect.exe -save ximfect_new.exe -action addskip -res img/ximfect.ico -mask ICONGROUP,MAINICON, -log NUL
    del ximfect.exe
    ren ximfect_new.exe ximfect.exe
echo Done!
echo.

:: Package the executable with effects into a distributable zip file
:: (using 7zip)
echo Packaging...
    move ximfect.exe pack > NUL
    cd pack
    if exist ximfect.zip del ximfect.zip
    7z a ximfect.zip ximfect.exe -bso0
    7z a ximfect.zip effects\ -bso0
    :: Usually the output would be more like `ximfect.zip` but we are in fancy town
    for /f "tokens=*" %%g in ('ximfect --version') do (set zipname=ximfect-v%%g.zip)
    ren ximfect.zip %zipname%
    move %zipname% .. > NUL
    cd ..
echo Done!
echo.