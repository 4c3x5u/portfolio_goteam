import axios from 'axios';
import { getAuthHeaders } from '../misc/util';

const UsersAPI = {
  // not in use – maintained for demonstration purposes
  get: (teamId, boardId) => {
    // team_id is mandatory
    let queryString = `?team_id=${teamId}`;
    if (boardId) { queryString += `&board_id=${boardId}`; }
    return axios.get(
      `${process.env.REACT_APP_BACKEND_URL}/users/${queryString}`,
      getAuthHeaders(),
    );
  },

  /// Only used for adding/removing a user from a board.
  post: (username, boardId, isActive) => axios.post(
    `${process.env.REACT_APP_BACKEND_URL}/users/`,
    {
      username,
      board_id: boardId,
      is_active: isActive,
    },
    getAuthHeaders(),
  ),

  delete: (username) => axios.delete(
    `${process.env.REACT_APP_BACKEND_URL}/users/?username=${username}`,
    getAuthHeaders(),
  ),
};

export default UsersAPI;
