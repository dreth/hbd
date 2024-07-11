// src/api/apiService.js
import api from './axiosConfig';

/**
 * Generate a new encryption key
 */
export const generateEncryptionKey = async () => {
  try {
    const response = await api.get('/generate-encryption-key');
    return response.data.encryption_key; // Return only the encryption_key
  } catch (error) {
    console.error('Error generating encryption key', error);
    throw error;
  }
};

/**
 * Register a new user
 * @param {Object} userData - The registration data
 */
export const registerUser = async (userData) => {
  try {
    const response = await api.post('/register', userData);
    return response.data;
  } catch (error) {
    console.error('Error registering user', error);
    throw error;
  }
};

/**
 * Login a user
 * @param {Object} loginData - The login data
 */
export const loginUser = async (loginData) => {
  try {
    const response = await api.post('/login', loginData);
    return response.data;
  } catch (error) {
    console.error('Error logging in user', error);
    throw error;
  }
};

/**
 * Modify user details
 * @param {Object} userData - The user data to modify
 */
export const modifyUser = async (userData) => {
  try {
    const response = await api.put('/modify-user', userData);
    return response.data;
  } catch (error) {
    console.error('Error modifying user', error);
    throw error;
  }
};

/**
 * Delete a user
 * @param {Object} userData - The user data to delete
 */
export const deleteUser = async (userData) => {
  try {
    const response = await api.delete('/delete-user', { data: userData });
    return response.data;
  } catch (error) {
    console.error('Error deleting user', error);
    throw error;
  }
};

/**
 * Add a new birthday
 * @param {Object} birthdayData - The birthday data to add
 */
export const addBirthday = async (birthdayData) => {
  try {
    const response = await api.post('/birthday', birthdayData);
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
export const modifyBirthday = async (birthdayData) => {
  try {
    const response = await api.put('/birthday', birthdayData);
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
export const deleteBirthday = async (birthdayData) => {
  try {
    const response = await api.delete('/birthday', { data: birthdayData });
    return response.data;
  } catch (error) {
    console.error('Error deleting birthday', error);
    throw error;
  }
};
