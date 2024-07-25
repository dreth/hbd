
import axios from 'axios';

let baseURL = 'http://localhost:8417';
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

// Add a request interceptor to include the Authorization header
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

/**
 * Generate a password
 */
export const generatePassword = async () => {
  try {
    const response = await api.get('/generate-password');
    return response.data.password; 
  } catch (error) {
    console.error('Error generating password', error);
    throw error;
  }
};

/**
 * Login user
 * @param {Object} loginData - The login data
 */
export const loginUser = async (userData) => {
  try {
    const response = await api.post('/login', userData);
    const { token, ...rest } = response.data;
    localStorage.setItem('token', token);  
    return { token, ...rest };
  } catch (error) {
    console.error('Error logging in', error);
    throw error;
  }
};



/**
 * Register user
 * @param {Object} userData - The user data
 */
export const registerUser = async (userData) => {
  try {
    const response = await api.post('/register', userData);
    const token = response.data.token;
    localStorage.setItem('token', token);
    return response.data;
  } catch (error) {
    console.error('Error registering user', error);
    throw error;
  }
};

/**
 * Get user data
 */
export const getUserData = async (token) => {
  try {
    const response = await api.get('/me');
    return response.data;
  } catch (error) {
    console.error('Error fetching user data', error);
    throw error;
  }
};

/**
 * Modify user details
 * @param {Object} userData - The new user data
 */
export const modifyUser = async (token, userData) => {
  try {
    const response = await api.put('/modify-user', userData);
    const token = response.data.token;
    if (token) {
      localStorage.setItem('token', token);
    }
    return response.data;
  } catch (error) {
    console.error('Error modifying user', error);
    throw error;
  }
};

/**
 * Add a new birthday
 * @param {Object} birthdayData - The birthday data to add
 */
export const addBirthday = async (token, birthdayData) => {
  try {
    const response = await api.post('/add-birthday', birthdayData);
    return response.data;
  } catch (error) {
    console.error('Error adding birthday', error);
    throw error;
  }
};

/**
 * Modify a birthday
 * @param {Object} birthdayData - The birthday data to modify
 */
export const modifyBirthday = async (token, birthdayData) => {
  try {
    const response = await api.put('/modify-birthday', birthdayData);
    return response.data;
  } catch (error) {
    console.error('Error modifying birthday', error);
    throw error;
  }
};

/**
 * Delete a birthday
 * @param {Object} birthdayData - The birthday data to delete
 */
export const deleteBirthday = async (token, birthdayData) => {
  try {
    const response = await api.delete('/delete-birthday', { data: birthdayData });
    return response.data;
  } catch (error) {
    console.error('Error deleting birthday', error);
    throw error;
  }
};

/**
 * Delete user
 */
export const deleteUser = async (p0) => {
  try {
    const response = await api.delete('/delete-user');
    localStorage.removeItem('token');
    return response.data;
  } catch (error) {
    console.error('Error deleting user', error);
    throw error;
  }
};
