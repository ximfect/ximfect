@echo off
echo :: Generating version
py gen-build-ver.py
cd ..\src
echo :: Building
go build
echo :: Deleting old executable
del ..\testing\ximfect.exe
echo :: Moving new executable
move .\ximfect.exe ..\testing > NUL
echo :: DONE
cd ..\testing