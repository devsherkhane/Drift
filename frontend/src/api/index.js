import axios from 'axios';
import router from '../router';

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
});

// REQUEST Interceptor: Attach token to every outgoing request
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// RESPONSE Interceptor: Handle expired tokens
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response && error.response.status === 401) {
      // If the backend says the token is invalid/expired
      localStorage.removeItem('token');
      router.push('/login'); 
    }
    return Promise.reject(error);
  }
);

export default api;