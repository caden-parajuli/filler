import mummy, mummy/routers, std/nativesockets

proc startHandler(request: Request) {.raises: [MummyError], gcsafe.} =
  request.headers["Connection"] = "Upgrade"
  request.headers["Upgrade"] = "websocket"
  let websocket = upgradeToWebSocket(request)

  # Message socket
  websocket.send("test")

proc socketHandler(websocket: WebSocket, event: WebSocketEvent, message: Message) =
  case event:
  of OpenEvent:
    echo "Socket opened"
    websocket.send("message received by server")
  of MessageEvent:
    echo message.kind, ": ", message.data
  of ErrorEvent:
    echo "ERROR"
  of CloseEvent:
    echo "Socket closed"

proc startServer(port: int = 8080, public = false) =
  let address = if public: "0.0.0.0" else: "localhost"
  var router: Router
  router.get("/ws", startHandler)
  
  let server = newServer(router, socketHandler)
  echo "Serving on " & address & " " & $port

  server.serve(port = Port(port), address = address)

startServer(port = 8080, public = false)