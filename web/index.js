import { makeRandomGrid } from "./modules/board.mjs";
import { initControls } from "./modules/controls.mjs";
import { connect, WS_SERVER_ADDRESS } from "./modules/socket.mjs"; 


function main() {
  document.addEventListener("DOMContentLoaded", makeRandomGrid, false);
  document.addEventListener("DOMContentLoaded", initControls, false);
  document.addEventListener("DOMContentLoaded", () => connect(WS_SERVER_ADDRESS), false);
}

main();
