import axios from 'axios';
import { getAuthHeaders } from '../misc/util';

const UserAPI = {
  /// Only used for adding/removing a user from a board.
  patch: (username, boardId, isActive) => axios.patch(
    `${process.env.SERVER_URL}/user?username=${username}`,
    { board_id: boardId, is_active: isActive },
    getAuthHeaders(),
  ),

  delete: (username) => axios.delete(
    `${process.env.SERVER_URL}/user?username=${username}`,
    getAuthHeaders(),
  ),
};

export default UserAPI;
