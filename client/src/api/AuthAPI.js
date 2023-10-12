import axios from 'axios';

const AuthAPI = {
  login: (username, password) => (
    axios.post(
      `${process.env.REACT_APP_SERVER_URL}/login`,
      { username, password },
      { withCredentials: true },
    )
  ),

  register: (username, password, inviteCode) => (
    axios.post(
      `${process.env.REACT_APP_SERVER_URL}/register?invite_code=${inviteCode}`,
      { username, password },
      { withCredentials: true },
    )
  ),
};

export default AuthAPI;
