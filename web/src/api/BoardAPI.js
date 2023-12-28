import axios from 'axios';

const apiUrl = process.env.REACT_APP_TEAM_SERVICE_URL + "/board"

const BoardAPI = {
  post: (boardData) => axios.post(apiUrl, boardData, { withCredentials: true }),

  delete: (boardId) => axios.delete(
    apiUrl + "?id=" + boardId, { withCredentials: true },
  ),

  patch: (boardId, boardData) => axios.patch(
    apiUrl + "?id=" + boardId, boardData, { withCredentials: true },
  ),
};

export default BoardAPI;
