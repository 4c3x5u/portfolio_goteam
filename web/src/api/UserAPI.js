import axios from 'axios';

var apiUrl = process.env.REACT_APP_USER_SERVICE_URL

const UserAPI = {
  login: (username, password) => (
    axios.post(
      apiUrl + "/login",
      { username, password },
      { withCredentials: true },
    )
  ),

  register: (username, password, inviteToken) => (
    axios.post(
      apiUrl + "/register?inviteToken=" + inviteToken,
      { username, password },
      { withCredentials: true },
    )
  ),
};

export default UserAPI;
