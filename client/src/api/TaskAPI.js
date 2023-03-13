import axios from 'axios';
import { getAuthHeaders } from '../misc/util';

const TaskAPI = {
  post: (task) => axios.post(
    `${process.env.REACT_APP_BACKEND_URL}/task`,
    task,
    getAuthHeaders(),
  ),

  patch: (taskId, data) => axios.patch(
    `${process.env.REACT_APP_BACKEND_URL}/task?id=${taskId}`,
    data,
    getAuthHeaders(),
  ),

  delete: (taskId) => axios.delete(
    `${process.env.REACT_APP_BACKEND_URL}/task?id=${taskId}`,
    getAuthHeaders(),
  ),
};

export default TaskAPI;
