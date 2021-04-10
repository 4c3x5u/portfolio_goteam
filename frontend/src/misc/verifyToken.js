import axios from 'axios';

const verifyToken = () => (
  axios.post(`${process.env.REACT_APP_BACKEND_URL}/verify-token/`, {
    username: sessionStorage.getItem('username'),
    token: sessionStorage.getItem('auth-token'),
  })
);

export default verifyToken;
