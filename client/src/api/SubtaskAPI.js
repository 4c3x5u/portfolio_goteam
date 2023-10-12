import axios from 'axios';

const SubtaskAPI = {
  patch: (id, data) => axios.patch(
    `${process.env.REACT_APP_SERVER_URL}/subtask?id=${id}`,
    data,
    { withCredentials: true },
  ),
};

export default SubtaskAPI;
