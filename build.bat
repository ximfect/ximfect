@echo off
cd src

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

echo.

:: Linux
set GOOS=linux
echo = Building for Linux...

    :: 32bit
    set GOARCH=386
    go build -o ximfect-linux32
    echo -- Built for 32 bit.
    :: 64bit
    set GOARCH=amd64
    go build -o ximfect-linux64
    echo -- Built for 64 bit.
    
echo.

:: Move executables to favorable location
move ximfect-* ..\pack > NUL
cd ..

echo --- Done!
echo.