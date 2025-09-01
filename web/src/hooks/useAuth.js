import { useState, useEffect, useCallback } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { authAPI } from '../services/api';
import { getToken, setToken, removeToken, getUserFromToken } from '../utils/auth';

export function useAuth() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const navigate = useNavigate();
  const location = useLocation();

  // Check if user is authenticated
  const isAuthenticated = useCallback(() => {
    const token = getToken();
    if (!token) return false;
    
    // Additional checks could be added here, like token expiration
    return true;
  }, []);

  // Login function
  const login = async (email, password) => {
    try {
      setLoading(true);
      setError(null);
      
      const response = await authAPI.login({ email, password });
      const { token, user: userData } = response.data;
      
      // Store token and update state
      setToken(token);
      setUser(userData);
      
      // Redirect to the originally requested page or home
      const from = location.state?.from?.pathname || '/';
      navigate(from, { replace: true });
      
      return { success: true };
    } catch (err) {
      const errorMessage = err.response?.data?.message || 'Login failed';
      setError(errorMessage);
      return { success: false, error: errorMessage };
    } finally {
      setLoading(false);
    }
  };

  // Logout function
  const logout = useCallback(() => {
    removeToken();
    setUser(null);
    navigate('/login');
  }, [navigate]);

  // Check auth status on mount
  useEffect(() => {
    const checkAuth = async () => {
      try {
        const token = getToken();
        if (token) {
          // If token exists, get user data
          const userData = getUserFromToken(token);
          if (userData) {
            setUser(userData);
          } else {
            // If token is invalid, clear it
            removeToken();
          }
        }
      } catch (err) {
        console.error('Auth check failed:', err);
        removeToken();
      } finally {
        setLoading(false);
      }
    };

    checkAuth();
  }, []);

  return {
    user,
    loading,
    error,
    isAuthenticated,
    login,
    logout,
  };
}

export default useAuth;
