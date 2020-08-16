cd src
	
# linux
env GOOS=linux echo = Building for Linux...

	env GOARCH=386 go build -o ximfect-linux32
	echo -- Built for 32 bit.
	
	env GOARCH=amd64 go build -o ximfect-linux64
	echo -- Built for 64 bit.
	
# move files to favorable location
mv ximfect-* ../pack
cd ..

echo --- Done!
