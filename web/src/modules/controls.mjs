import { changePlayerColor, numColors, player_positions, player_num } from "./board.mjs";
import { sendMove } from "./socket.mjs";

export function initControls() {
    for (let colorNum = 0; colorNum < numColors; ++colorNum) {
        let control = document.getElementById(`diamond-control${colorNum}`);
        control.addEventListener("click", (_) => {
            // TODO Do some validation first to make sure the control is pressable
            // @ts-ignore
            changePlayerColor(player_positions[player_num], control.dataset.color);
            sendMove(colorNum)
        });
    }
}

export function deactivateControls() {
    for (let i = 0; i < numColors; ++i) {
        let control = document.getElementById(`diamond-control${i}`);
        control.setAttribute('disabled', 'disabled');
    }
}

export function reactivateControls() {
    for (let i = 0; i < numColors; ++i) {
        let control = document.getElementById(`diamond-control${i}`);
        control.removeAttribute('disabled');
    }
}

/**
 * @param {boolean} my_turn
 */
export function setTurn(my_turn) {
    if (my_turn) {
        reactivateControls()
    } else {
        deactivateControls()
    }
}
