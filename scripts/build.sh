#!/usr/bin/env bash
set -e

SERVER_URL="http://127.0.0.1:8000/docs"
SERVER_HOST="127.0.0.1"
SERVER_PORT="8000"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

VENV_DIR="$ROOT_DIR/.venv"
SERVER_DIR="$ROOT_DIR/embedding_server"
LOG_FILE="$SERVER_DIR/embed_server.log"
APP_DIR="$ROOT_DIR/muninx"

# ---------------------------------------------------------------------------
# 1. Start the embedding server if it is not already running.
# ---------------------------------------------------------------------------

if curl -s "$SERVER_URL" > /dev/null; then
    echo "Embed server already running at $SERVER_URL"
else
    echo "Embed server not running. Starting server..."

    if [ ! -d "$VENV_DIR" ]; then
        echo "Error: .venv not found at $VENV_DIR"
        exit 1
    fi

    if [ ! -d "$SERVER_DIR" ]; then
        echo "Error: embedding_server directory not found at $SERVER_DIR"
        exit 1
    fi

    if [ ! -f "$SERVER_DIR/embed.py" ]; then
        echo "Error: embed.py not found at $SERVER_DIR/embed.py"
        exit 1
    fi

    source "$VENV_DIR/bin/activate"

    if ! python -m uvicorn --version > /dev/null 2>&1; then
        echo "Error: uvicorn is not installed in .venv"
        echo "Install it with:"
        echo "  source $VENV_DIR/bin/activate"
        echo "  pip install fastapi uvicorn"
        exit 1
    fi

    nohup python -m uvicorn embed:app \
        --app-dir "$SERVER_DIR" \
        --host "$SERVER_HOST" \
        --port "$SERVER_PORT" \
        > "$LOG_FILE" 2>&1 &

    SERVER_PID=$!
    echo "Started embed server with PID $SERVER_PID"
    echo "Logs: $LOG_FILE"

    sleep 2

    if curl -s "$SERVER_URL" > /dev/null; then
        echo "Embed server started successfully at $SERVER_URL"
    else
        echo "Error: embed server failed to start."
        echo "Check logs:"
        echo "  cat $LOG_FILE"
        exit 1
    fi
fi

# ---------------------------------------------------------------------------
# 2. Build and run the app.
# ---------------------------------------------------------------------------

echo "Building app..."
cd "$APP_DIR"
go build -o bin/muninx .

echo "Starting app..."
exec "$APP_DIR/bin/muninx"
