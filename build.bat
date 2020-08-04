@echo off

echo Building...
cd src
if exist ximfect.exe del ximfect.exe
go build
move ximfect.exe .. > NUL
cd ..
echo Done!