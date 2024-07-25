import axios from 'axios';

var baseURL = 'https://hbd-api.lotiguere.com';
if (process.env.ENVIRONMENT === 'production') {
  baseURL = 'https://hbd-api.lotiguere.com';
}

// Create an Axios instance
const api = axios.create({
  baseURL: baseURL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Interceptor to handle requests
api.interceptors.request.use(
  (config) => {
    // Add authorization token or other headers if necessary
    // config.headers.Authorization = `Bearer ${token}`;
    return config;
  },
  (error) => {
    // Handle the error
    return Promise.reject(error);
  }
);

// Interceptor to handle responses
api.interceptors.response.use(
  (response) => response,
  (error) => {
    // Handle errors, e.g., display notifications or log errors
    return Promise.reject(error);
  }
);

export default api;
