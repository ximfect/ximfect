echo : Generate src/tool/const.go
python3 ./generate_const.py
cd ../src

echo : Build
go build
export ret="$?"
if [[ $ret != 0 ]]; then echo "Exit code: $ret" && exit; fi;

echo : Delete old executable
rm ../testing/ximfect

echo : Move new executable
mv ./ximfect ../testing
cd ../testing
