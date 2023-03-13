import axios from 'axios';
import { getAuthHeaders } from '../misc/util';

const BoardAPI = {
  post: (boardData) => axios.post(
    `${process.env.SERVER_URL}/board`,
    boardData,
    getAuthHeaders(),
  ),

  delete: (boardId) => axios.delete(
    `${process.env.SERVER_URL}/board?id=${boardId}`,
    getAuthHeaders(),
  ),

  patch: (boardId, boardData) => axios.patch(
    `${process.env.SERVER_URL}/board?id=${boardId}`,
    boardData,
    getAuthHeaders(),
  ),
};

export default BoardAPI;
