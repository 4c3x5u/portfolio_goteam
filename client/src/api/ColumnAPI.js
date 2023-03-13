import axios from 'axios';
import { getAuthHeaders } from '../misc/util';

const ColumnAPI = {
  patch: (columnId, data) => axios.patch(
    `${process.env.SERVER_URL}/column?id=${columnId}`,
    data,
    getAuthHeaders(),
  ),
};

export default ColumnAPI;
