when defined(c):
  # Karax build html
  import karax / [karaxdsl, vdom]

  template kxi(): int = 0
  template addEventHandler(n: VNode; k: EventKind; action: string; kxi: int) =
    n.setAttr($k, action)

  let pageHtml = buildHtml(html):
    head:
      link(rel = "stylesheet", href = "filler.css")
      script(src = "filler.js")
    body:
      tdiv(id = "joinGameDiv"): 
        form(id = "joinForm", onsubmit = "return joinGame();"):
          text "Name: "
          input(type = "text", id = "nameInput", name = "nameInput", value = "")
          input(type = "submit", id = "submitNameButton", name = "submitNameButton", value = "Join")

  echo pageHtml


when defined(js):
  import std/[dom]
  import jswebsockets

  var name: cstring = "Name"

  proc hide(element: Element) {.inline.} =
    element.style.display = "none"
  proc show(element: Element) {.inline.} =
    element.style.display = "block"
  
  
  proc joinGame(): bool {.exportc.} =
    name = getElementById("nameInput").value
    getElementById("joinForm").FormElement.reset()
    getElementById("joinGameDiv").hide()
    echo "Name: " & $name
    echo "Opening socket at ws://" & $window.location.host & "/ws"
    var
      socket = newWebSocket(cstring("ws://" & $window.location.host & "/ws"))

    socket.onOpen = proc (e: Event) =
      echo("sent: Client opened")
      socket.send("Client opened")
      
    socket.onMessage = proc (e: MessageEvent) =
      echo("received: ",e.data)
      socket.close(StatusCode(1000),"received msg")
      
    socket.onClose = proc (e: CloseEvent) =
      echo("closing: ",e.reason)
    return false