import axios from 'axios';

const verifyToken = (setCurrentUser) => (
  axios.post(`${process.env.REACT_APP_BACKEND_URL}/verify-token/`, {
    username: sessionStorage.getItem('username'),
    token: sessionStorage.getItem('auth-token'),
  }).then((res) => setCurrentUser({
    username: res.data.username,
    teamId: res.data.teamId,
    isAdmin: res.data.isAdmin,
    isAuthenticated: true,
  })).catch(() => setCurrentUser({
    username: '',
    teamId: null,
    isAdmin: false,
    isAuthenticated: false,
  }))
);

export default verifyToken;
