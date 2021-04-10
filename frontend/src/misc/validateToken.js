import axios from 'axios';

const validateToken = (setAuthenticated) => (
  axios
    .post(`${process.env.REACT_APP_BACKEND_URL}/verify-token/`, {
      username: sessionStorage.getItem('username'),
      token: sessionStorage.getItem('auth-token'),
    })
    .then(() => setAuthenticated(true))
    .catch(() => setAuthenticated(false))
);

export default validateToken;
