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
};

export default UsersAPI;
