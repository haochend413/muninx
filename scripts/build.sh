cd ..

# Start embed server if not already running
if ! curl -s http://127.0.0.1:8000/docs > /dev/null; then
    echo "Embed server not running. Starting server..."

    cd embedding_server || exit 1
    nohup uvicorn embed_server:app --host 127.0.0.1 --port 8000 > embed_server.log 2>&1 &
    cd ..

    sleep 2
fi

# Build muninx
if [ -f ./muninx/bin/muninx ]; then
    rm ./muninx/bin/muninx
fi

cd muninx
go build -o ./bin/muninx
./bin/muninx