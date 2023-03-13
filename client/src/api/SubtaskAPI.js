import axios from 'axios';
import { getAuthHeaders } from '../misc/util';

const SubtaskAPI = {
  patch: (id, data) => axios.patch(
    `${process.env.SERVER_URL}/subtask?id=${id}`,
    data,
    getAuthHeaders(),
  ),
};

export default SubtaskAPI;
