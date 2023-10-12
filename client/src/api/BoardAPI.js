import axios from 'axios';

const BoardAPI = {
  post: (boardData) => axios.post(
    `${process.env.REACT_APP_SERVER_URL}/board`,
    boardData,
    { withCredentials: true },
  ),

  delete: (boardId) => axios.delete(
    `${process.env.REACT_APP_SERVER_URL}/board?id=${boardId}`,
    { withCredentials: true },
  ),

  patch: (boardId, boardData) => axios.patch(
    `${process.env.REACT_APP_SERVER_URL}/board?id=${boardId}`,
    boardData,
    { withCredentials: true },
  ),
};

export default BoardAPI;
