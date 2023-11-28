
SCRIPT_DIR="$(dirname "$(realpath "$0")")"
PROJ_DIR="$(realpath "$SCRIPT_DIR/..")"
REMOTE_JSON_URL=https://debricked.com/api/1.0/open/files/supported-formats
LOCAL_JSON_DIR=$PROJ_DIR/internal/file/embedded
LOCAL_JSON_FILE=$PROJ_DIR/internal/file/embedded/supported_formats.json

echo "Supported-formats is downloaded from remote for offline backup"
mkdir -p $LOCAL_JSON_DIR && wget -O $LOCAL_JSON_FILE $REMOTE_JSON_URL