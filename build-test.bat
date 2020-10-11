@echo off
echo :: Generating version
py gen-build-ver.py
cd src
echo :: Building
go build
echo :: Deleting old executable
del ..\ximfect.exe
echo :: Moving new executable
move .\ximfect.exe .. > NUL
echo :: DONE
cd ..