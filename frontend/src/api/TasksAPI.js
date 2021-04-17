import axios from 'axios';
import { getAuthHeaders } from '../misc/util';

const TasksAPI = {
  get: (columnId) => axios.get(
    `${process.env.REACT_APP_BACKEND_URL}/tasks/?column_id=${columnId}`,
    getAuthHeaders(),
  ),
  post: (task) => axios.post(
    `${process.env.REACT_APP_BACKEND_URL}/tasks/`,
    task,
    getAuthHeaders(),
  ),
  patch: (taskId, data) => axios.patch(
    `${process.env.REACT_APP_BACKEND_URL}/tasks/?id=${taskId}`,
    data,
    getAuthHeaders(),
  ),
  delete: (taskId) => axios.delete(
    `${process.env.REACT_APP_BACKEND_URL}/tasks/?id=${taskId}`,
    getAuthHeaders(),
  ),
};

export default TasksAPI;
