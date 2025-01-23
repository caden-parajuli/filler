/** 
 * @typedef {IdMessageRaw | GameParamsReqRaw | GameParamsRespRaw | MoveMessageRaw | ClientMoveMessageRaw | JoinGameReqRaw | JoinGameRespRaw | OtherClientJoinRaw} Message
 */

/** 
 * @typedef IdMessageRaw
 * @property {'id_message'} message_type
 * @property {IdMessage} message
 *
 * @typedef IdMessage
 * @property {string} id
 */

/** 
 * @typedef GameParamsReqRaw
 * @property {'game_params_req'} message_type
 * @property {GameParamsReq} message
 *
 * @typedef GameParamsReq
 * @property {string} id
 * @property {boolean} is_diamonds
 * @property {number} num_rows
 * @property {number} num_cols
 * @property {number} num_colors
 */

/** 
 * @typedef GameParamsRespRaw
 * @property {'game_params_resp'} message_type
 * @property {GameParamsResp} message
 *
 * @typedef GameParamsResp
 * @property {Board} board
 * @property {number} game_id
 */

/** 
 * @typedef MoveMessageRaw
 * @property {'move_message'} message_type
 * @property {MoveMessage} message
 *
 * @typedef MoveMessage
 * @property {Board} board
 * @property {boolean} my_turn
 */

/** 
 * @typedef ClientMoveMessageRaw
 * @property {'client_move_message'} message_type
 * @property {ClientMoveMessage} message
 *
 * @typedef ClientMoveMessage
 * @property {string} id
 * @property {number} color
 */

/** 
 * @typedef JoinGameReqRaw
 * @property {'join_game_req'} message_type
 * @property {JoinGameReq} message
 *
 * @typedef JoinGameReq
 * @property {string} id
 * @property {number} game_id
 */

/** 
 * @typedef JoinGameRespRaw
 * @property {'join_game_resp'} message_type
 * @property {JoinGameResp} message
 *
 * @typedef JoinGameResp
 * @property {boolean} success
 * @property {number} player_num
 * @property {Board} board
 * @property {boolean} my_turn
 */

/**
 * @typedef OtherClientJoinRaw
 * @property {'other_client_join'} message_type
 * @property {OtherClientJoin} message
 *
 * @typedef OtherClientJoin
 * @property {boolean} my_turn
 */
