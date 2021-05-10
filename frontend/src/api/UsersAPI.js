import axios from 'axios';
import { getAuthHeaders } from '../misc/util';

const UsersAPI = {
  /// Only used for adding/removing a user from a board.
  patch: (username, boardId, isActive) => axios.patch(
    `${process.env.REACT_APP_BACKEND_URL}/users/?username=${username}`,
    { board_id: boardId, is_active: isActive },
    getAuthHeaders(),
  ),

  delete: (username) => axios.delete(
    `${process.env.REACT_APP_BACKEND_URL}/users/?username=${username}`,
    getAuthHeaders(),
  ),
};

export default UsersAPI;
