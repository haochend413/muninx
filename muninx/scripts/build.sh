cd ..
if [ -f ./bin/muninx ]; then
    rm ./bin/muninx
fi
go build -o ./bin/muninx
./bin/muninx