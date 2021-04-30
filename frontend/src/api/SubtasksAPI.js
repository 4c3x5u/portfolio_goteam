import axios from 'axios';
import { getAuthHeaders } from '../misc/util';

const SubtasksAPI = {
  // not in use – maintained for demonstration purposes
  get: (taskId) => axios.get(
    `${process.env.REACT_APP_BACKEND_URL}/subtasks/?task_id=${taskId}`,
    getAuthHeaders(),
  ),

  patch: (id, data) => axios.patch(
    `${process.env.REACT_APP_BACKEND_URL}/subtasks/?id=${id}`,
    data,
    getAuthHeaders(),
  ),
};

export default SubtasksAPI;
