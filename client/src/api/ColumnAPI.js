import axios from 'axios';

const ColumnAPI = {
  patch: (columnId, data) => axios.patch(
    `${process.env.REACT_APP_SERVER_URL}/column?id=${columnId}`,
    data,
    { withCredentials: true },
  ),
};

export default ColumnAPI;
