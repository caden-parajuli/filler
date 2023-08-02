when defined("c"):
  # Karax build html
  import karax / [karaxdsl, vdom]

  template kxi(): int = 0
  template addEventHandler(n: VNode; k: EventKind; action: string; kxi: int) =
    n.setAttr($k, action)

  let pageHtml = buildHtml(html):
    head:
      script(src="index.js")
    body:
      p:
        "This is my little test"
      
  echo pageHtml


when defined("js"):
  # This gets compiled to js
  