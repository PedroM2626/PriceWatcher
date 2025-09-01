// Store JWT token in localStorage
export const setAuthToken = (token) => {
  if (token) {
    localStorage.setItem('token', token);
  } else {
    localStorage.removeItem('token');
  }
};

// Get JWT token from localStorage
export const getAuthToken = () => {
  return localStorage.getItem('token');
};

// Remove JWT token from localStorage
export const removeAuthToken = () => {
  localStorage.removeItem('token');
};

// Check if user is authenticated
export const isAuthenticated = () => {
  const token = getAuthToken();
  if (!token) return false;
  
  // Check if token is expired
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    return payload.exp * 1000 > Date.now();
  } catch (error) {
    return false;
  }
};

// Get user info from token
export const getUserFromToken = () => {
  const token = getAuthToken();
  if (!token) return null;
  
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    return {
      id: payload.sub,
      email: payload.email,
      name: payload.name,
      role: payload.role,
    };
  } catch (error) {
    return null;
  }
};

// Handle successful authentication
export const handleAuthentication = (token) => {
  setAuthToken(token);
  // You can add additional logic here, like redirecting to dashboard
  window.location.href = '/';
};

// Logout user
export const logout = () => {
  removeAuthToken();
  // You can add additional cleanup here
  window.location.href = '/login';
};
