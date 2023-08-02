import mummy, mummy/routers

proc startHandler(request: Request) =
  let websocket = request.upgradeToWebSocket()
  websocket.send("test")

proc socketHandler(websocket: WebSocket, event: WebSocketEvent, message: Message) =
  case event:
  of OpenEvent:
    discard
  of MessageEvent:
    echo message.kind, ": ", message.data
  of ErrorEvent:
    discard
  of CloseEvent:
    discard

proc startServer(port: int = 8080, public = false) =
  let address = if public: "0.0.0.0" else: "localhost"
  var router: Router
  router.get("/ws", startHandler)
  
  let server = newServer(router, socketHandler)
  echo "Serving on " & address & " " & $port

  server.serve(port = Port(port), address = address)