when defined(c):
  # Karax build html
  import karax / [karaxdsl, vdom]

  template kxi(): int = 0
  template addEventHandler(n: VNode; k: EventKind; action: string; kxi: int) =
    n.setAttr($k, action)

  let pageHtml = buildHtml(html):
    head:
      script(src = "filler.js")
    body:
      tdiv(id = "joinGame"): 
        form(id = "joinForm"):
          text "Name: "
          input(type = "text", id = "nameInput", name = "nameInput", value = "")
          input(type = "button", id = "submitNameButton", name = "submitNameButton", value = "Join", onClick = "joinGame()")

  echo pageHtml


when defined(js):
  import std/[dom]
  
  proc joinGame() {.exportc.} =
    var name = getElementById("nameInput").value
    echo "Name: " & $name