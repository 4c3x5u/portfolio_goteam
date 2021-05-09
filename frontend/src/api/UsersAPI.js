import axios from 'axios';
import { getAuthHeaders } from '../misc/util';

const UsersAPI = {
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
