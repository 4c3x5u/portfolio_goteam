import axios from 'axios';

const apiUrl = process.env.REACT_APP_TASK_SERVICE_URL + "/tasks"

const TasksAPI = {
  get: (boardID) => axios.get(
    apiUrl + "?boardID=" + boardID, { withCredentials: true },
  ),

  patch: (data) => axios.patch(apiUrl, data, { withCredentials: true }),
};

export default TasksAPI;
