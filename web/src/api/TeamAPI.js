import axios from 'axios';

const TeamAPI = {
  get: () => axios.get(
    process.env.REACT_APP_TEAM_SERVICE_URL + "/team", { withCredentials: true },
  ),
};

export default TeamAPI;
