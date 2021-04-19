import axios from 'axios';
import { getAuthHeaders } from '../misc/util';

const TeamsAPI = {
  get: (teamId) => axios.get(
    `${process.env.REACT_APP_BACKEND_URL}/teams/?team_id=${teamId}`,
    getAuthHeaders(),
  ),
};

export default TeamsAPI;
