import axios from 'axios';

const api = axios.create({
  // Point this to your Go server's port
  baseURL: 'http://localhost:8080/api',
});

// Interceptor to automatically attach your JWT token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default api;