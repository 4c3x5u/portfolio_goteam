import axios from 'axios';
import { getAuthHeaders } from '../misc/util';

const UsersAPI = {
  get: (teamId, boardId) => axios.get(
    `${
      process.env.REACT_APP_BACKEND_URL
    }/users/?team_id=${
      teamId
    }&board_id=${
      boardId
    }`,
    getAuthHeaders(),
  ),

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
