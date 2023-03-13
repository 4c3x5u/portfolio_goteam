import axios from 'axios';

const AuthAPI = {
  login: (username, password) => (
    axios.post(`${process.env.SERVER_URL}/login`, {
      username,
      password,
    })
  ),

  register: (username, password) => (
    axios.post(`${process.env.SERVER_URL}/register`, {
      username,
      password,
    })
  ),
};

export default AuthAPI;
