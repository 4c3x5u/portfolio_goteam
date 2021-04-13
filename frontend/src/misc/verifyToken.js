import axios from 'axios';

const verifyToken = () => (
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

export default verifyToken;
