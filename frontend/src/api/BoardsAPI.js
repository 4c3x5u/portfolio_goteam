import axios from 'axios';
import { getAuthHeaders } from '../misc/util';

const BoardsAPI = {
  get: (boardId, teamId) => {
    let queryString = '';
    if (boardId || teamId) { queryString += '?'; }
    if (boardId) { queryString += `id=${boardId}`; }
    if (boardId && teamId) { queryString += '&'; }
    if (teamId) { queryString += `team_id=${teamId}`; }
    return axios.get(
      `${process.env.REACT_APP_BACKEND_URL}/boards/${queryString}`,
      getAuthHeaders(),
    );
  },
  post: (boardData) => axios.post(
    `${process.env.REACT_APP_BACKEND_URL}/boards/`,
    boardData,
    getAuthHeaders(),
  ),
  delete: (boardId) => axios.delete(
    `${process.env.REACT_APP_BACKEND_URL}/boards/?id=${boardId}`,
    getAuthHeaders(),
  ),
  patch: (boardId, boardData) => axios.patch(
    `${process.env.REACT_APP_BACKEND_URL}/boards/?id=${boardId}`,
    boardData,
    getAuthHeaders(),
  ),
};

export default BoardsAPI;
