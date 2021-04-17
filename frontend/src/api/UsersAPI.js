import axios from 'axios';
import { getAuthHeaders } from '../misc/util';

const UsersAPI = {
  get: (teamId) => axios.get(
    `${process.env.REACT_APP_BACKEND_URL}/users/?team_id=${teamId}`,
    getAuthHeaders(),
  ),
};

export default UsersAPI;
