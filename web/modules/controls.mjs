import { changePlayerColor, numColors } from "./board.mjs";

export function initControls() {
  for (let i = 0; i < numColors; ++i) {
    let control = document.querySelector(`#diamond-control${i}`);
    control.addEventListener("click", (_) => {
      // TODO Do some validation first to make sure the control is pressable
      // @ts-ignore
      changePlayerColor(player_positions[player_num], control.dataset.color)
      // TODO message server
    });
  }
}
