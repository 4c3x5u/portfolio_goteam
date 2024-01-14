import axios from 'axios';

const apiUrl = process.env.REACT_APP_TASK_SERVICE_URL + "/task"

const TaskAPI = {
  post: (task) => axios.post(
    apiUrl, task, { withCredentials: true },
  ),

  patch: (task) => axios.patch(
    apiUrl, task, { withCredentials: true },
  ),

  delete: (taskId) => axios.delete(
    apiUrl + "?id=" + taskId, { withCredentials: true },
  ),
};

export default TaskAPI;
