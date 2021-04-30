import axios from 'axios';
import { getAuthHeaders } from '../misc/util';

const ColumnsAPI = {
  // not in use – maintained for demonstration purposes
  get: (boardId) => axios.get(
    `${process.env.REACT_APP_BACKEND_URL}/columns/?board_id=${boardId}`,
    getAuthHeaders(),
  ),

  patch: (columnId, data) => axios.patch(
    `${process.env.REACT_APP_BACKEND_URL}/columns/?id=${columnId}`,
    data,
    getAuthHeaders(),
  ),
};

export default ColumnsAPI;
