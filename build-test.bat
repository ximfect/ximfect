@echo off
cd src
go build
del ..\ximfect.exe
move .\ximfect.exe .. > NUL
cd ..