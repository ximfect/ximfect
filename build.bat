@echo off
cd src
setlocal

:: Windows
set GOOS=windows
echo = Building for Windows...

    :: 32bit
    set GOARCH=386
    go build -o ximfect-windows32.exe
    echo -- Built for 32 bit.
    :: 64bit
    set GOARCH=amd64
    go build -o ximfect-windows64.exe
    echo -- Built for 64 bit.

:: Move executables to favorable location
move ximfect-* ..\pack > NUL
cd ..

echo --- Done!
echo.
endlocal