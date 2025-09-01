import axios from 'axios';
import { getAuthToken, removeAuthToken } from '../utils/auth';

// Create axios instance with base URL and common headers
const api = axios.create({
  baseURL: process.env.REACT_APP_API_URL || 'http://localhost:8080/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add a request interceptor to include the auth token
api.interceptors.request.use(
  (config) => {
    const token = getAuthToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Add a response interceptor to handle common errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response) {
      // Handle 401 Unauthorized
      if (error.response.status === 401) {
        // Clear auth data and redirect to login
        removeAuthToken();
        window.location.href = '/login';
      }
      
      // Extract error message from response
      const message = error.response.data?.message || 'An error occurred';
      return Promise.reject(new Error(message));
    }
    return Promise.reject(error);
  }
);

// Auth API
export const authAPI = {
  login: (email, password) => api.post('/auth/login', { email, password }),
  register: (userData) => api.post('/auth/register', userData),
  refreshToken: (refreshToken) => api.post('/auth/refresh-token', { refreshToken }),
  getMe: () => api.get('/auth/me'),
};

// Products API
export const productsAPI = {
  getAll: (params = {}) => api.get('/products', { params }),
  getById: (id) => api.get(`/products/${id}`),
  create: (productData) => api.post('/products', productData),
  update: (id, productData) => api.put(`/products/${id}`, productData),
  delete: (id) => api.delete(`/products/${id}`),
  checkPrice: (id) => api.post(`/products/${id}/check`),
  getPriceHistory: (id, params = {}) => 
    api.get(`/products/${id}/history`, { params }),
};

// Alerts API
export const alertsAPI = {
  getAll: (params = {}) => api.get('/alerts', { params }),
  getById: (id) => api.get(`/alerts/${id}`),
  create: (alertData) => api.post('/alerts', alertData),
  update: (id, alertData) => api.put(`/alerts/${id}`, alertData),
  delete: (id) => api.delete(`/alerts/${id}`),
  toggleStatus: (id) => api.patch(`/alerts/${id}/toggle`),
};

// User API
export const userAPI = {
  getProfile: () => api.get('/users/me'),
  updateProfile: (userData) => api.put('/users/me', userData),
  changePassword: (currentPassword, newPassword) => 
    api.post('/users/change-password', { currentPassword, newPassword }),
  getSessions: () => api.get('/users/me/sessions'),
  revokeSession: (sessionId) => api.delete(`/users/me/sessions/${sessionId}`),
};

export default api;
