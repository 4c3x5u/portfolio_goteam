import axios from 'axios';

const TaskAPI = {
  post: (task) => axios.post(
    `${process.env.REACT_APP_SERVER_URL}/task`,
    task,
    { withCredentials: true },
  ),

  patch: (taskId, data) => axios.patch(
    `${process.env.REACT_APP_SERVER_URL}/task?id=${taskId}`,
    data,
    { withCredentials: true },
  ),

  delete: (taskId) => axios.delete(
    `${process.env.REACT_APP_SERVER_URL}/task?id=${taskId}`,
    { withCredentials: true },
  ),
};

export default TaskAPI;
