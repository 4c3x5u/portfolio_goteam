import axios from 'axios';
import { getAuthHeaders } from '../misc/util';

const SubtasksAPI = {
  patch: (id, data) => axios.patch(
    `${process.env.REACT_APP_BACKEND_URL}/subtasks/?id=${id}`,
    data,
    getAuthHeaders(),
  ),
};

export default SubtasksAPI;
