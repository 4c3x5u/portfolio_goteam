import axios from 'axios';

export const verifyToken = () => (
  axios.post(`${process.env.REACT_APP_BACKEND_URL}/verify-token/`, {
    username: sessionStorage.getItem('username'),
    token: sessionStorage.getItem('auth-token'),
  }).then((res) => ({
    username: res.data.username,
    token: sessionStorage.getItem('auth-token'),
    teamId: res.data.teamId,
    isAdmin: res.data.isAdmin,
    isAuthenticated: true,
  }))
);

const authHeaders = {
  headers: {
    'auth-user': sessionStorage.getItem('username'),
    'auth-token': sessionStorage.getItem('auth-token'),
  },
};

export const getBoards = (teamId) => axios.get(
  `${process.env.REACT_APP_BACKEND_URL}/boards/?team_id=${teamId}`,
  authHeaders,
);

export const getColumns = (boardId) => axios.get(
  `${process.env.REACT_APP_BACKEND_URL}/columns/?board_id=${boardId}`,
  authHeaders,
);

export const getTasks = (columnId) => axios.get(
  `${process.env.REACT_APP_BACKEND_URL}/tasks/?column_id=${columnId}`,
  authHeaders,
);

export const getSubtasks = (taskId) => axios.get(
  `${process.env.REACT_APP_BACKEND_URL}/subtasks/?task_id=${taskId}`,
  authHeaders,
);
