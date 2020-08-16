cd src

# windows
env GOOS=windows echo = Building for Windows...

	env GOARCH=386 go build -o ximfect-windows32.exe
	echo -- Built for 32 bit.
	
	env GOARCH=amd64 go build -o ximfect-windows64.exe
	echo -- Built for 64 bit.
	
# linux
env GOOS=linux echo = Building for Linux...

	env GOARCH=386 go build -o ximfect-linux32
	echo -- Built for 32 bit.
	
	env GOARCH=amd64 go build -o ximfect-linux64
	echo -- Built for 64 bit.
	
# move files to favorable location
sudo mv ximfect-* ../pack
cd ..

echo --- Done!
