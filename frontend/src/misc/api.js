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

const getAuthHeaders = () => ({
  headers: {
    'auth-user': sessionStorage.getItem('username'),
    'auth-token': sessionStorage.getItem('auth-token'),
  },
});

export const getBoards = (teamId) => axios.get(
  `${process.env.REACT_APP_BACKEND_URL}/boards/?team_id=${teamId}`,
  getAuthHeaders(),
);

export const getColumns = (boardId) => axios.get(
  `${process.env.REACT_APP_BACKEND_URL}/columns/?board_id=${boardId}`,
  getAuthHeaders(),
);

export const getTasks = (columnId) => axios.get(
  `${process.env.REACT_APP_BACKEND_URL}/tasks/?column_id=${columnId}`,
  getAuthHeaders(),
);

export const postTask = (task) => axios.post(
  `${process.env.REACT_APP_BACKEND_URL}/tasks/`,
  task,
  getAuthHeaders(),
);

export const patchTask = (data) => axios.patch(
  `${process.env.REACT_APP_BACKEND_URL}/tasks/`,
  data,
  getAuthHeaders(),
);

export const getSubtasks = (taskId) => axios.get(
  `${process.env.REACT_APP_BACKEND_URL}/subtasks/?task_id=${taskId}`,
  getAuthHeaders(),
);
