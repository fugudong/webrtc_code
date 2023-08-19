package connect

// This is a temporary solution to avoid holding a zombie connection forever, by
// setting a 1 day timeout on reading from the WebSocket connection.
const WS_READ_TIMEOUT_SEC = 60 * 60 * 2