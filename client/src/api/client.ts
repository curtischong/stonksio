import axios from 'axios';

const apiClient = () => {
  return axios.create({
    baseURL: 'https://stonk.st/api/',
  });
};

export default apiClient;
