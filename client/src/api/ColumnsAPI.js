import axios from 'axios';
import { getAuthHeaders } from '../misc/util';

const ColumnsAPI = {
  patch: (columnId, data) => axios.patch(
    `${process.env.REACT_APP_BACKEND_URL}/columns/?id=${columnId}`,
    data,
    getAuthHeaders(),
  ),
};

export default ColumnsAPI;
