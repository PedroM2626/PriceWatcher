import React, { useState, useEffect } from 'react';
import { useNavigate, Link as RouterLink, useLocation } from 'react-router-dom';
import {
  Container,
  Box,
  Typography,
  TextField,
  Button,
  InputAdornment,
  IconButton,
  Link,
  FormControlLabel,
  Checkbox,
  Paper,
  Alert,
  CircularProgress,
} from '@mui/material';
import { LockOutlined, Visibility, VisibilityOff } from '@mui/icons-material';
import useAuth from '../hooks/useAuth';

const Login = () => {
  const [showPassword, setShowPassword] = useState(false);
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    remember: false,
  });
  
  const { login, loading, error, isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const from = location.state?.from?.pathname || '/';

  // Redirect if already authenticated
  useEffect(() => {
    if (isAuthenticated()) {
      navigate(from, { replace: true });
    }
  }, [isAuthenticated, navigate, from]);

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!formData.email || !formData.password) {
      return; // Let form validation handle this
    }

    await login(formData.email, formData.password);
  };

  if (loading && !formData.email) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="100vh">
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Container component="main" maxWidth="xs">
      <Box
        sx={{
          marginTop: 8,
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
        }}
      >
        <Paper elevation={3} sx={{ p: 4, width: '100%', mt: 4 }}>
          <Box display="flex" flexDirection="column" alignItems="center" mb={3}>
            <Box
              sx={{
                width: 40,
                height: 40,
                borderRadius: '50%',
                backgroundColor: 'primary.main',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                mb: 1
              }}
            >
              <LockOutlined sx={{ color: 'white' }} />
            </Box>
            <Typography component="h1" variant="h5">
              Sign in to PriceWatcher
            </Typography>
          </Box>
          
          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}
          
          <Box component="form" onSubmit={handleSubmit} noValidate sx={{ mt: 1 }}>
            <TextField
              margin="normal"
              required
              fullWidth
              id="email"
              label="Email Address"
              name="email"
              autoComplete="email"
              autoFocus
              value={formData.email}
              onChange={handleChange}
              disabled={loading}
              inputProps={{
                'data-testid': 'email-input'
              }}
            />
            <TextField
              margin="normal"
              required
              fullWidth
              name="password"
              label="Password"
              type={showPassword ? 'text' : 'password'}
              id="password"
              autoComplete="current-password"
              value={formData.password}
              onChange={handleChange}
              disabled={loading}
              inputProps={{
                'data-testid': 'password-input'
              }}
              InputProps={{
                endAdornment: (
                  <InputAdornment position="end">
                    <IconButton
                      aria-label="toggle password visibility"
                      onClick={() => setShowPassword(!showPassword)}
                      onMouseDown={(e) => e.preventDefault()}
                      edge="end"
                      disabled={loading}
                    >
                      {showPassword ? <VisibilityOff /> : <Visibility />}
                    </IconButton>
                  </InputAdornment>
                ),
              }}
            />
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <FormControlLabel
                control={
                  <Checkbox 
                    value="remember" 
                    color="primary" 
                    name="remember"
                    checked={formData.remember}
                    onChange={handleChange}
                    disabled={loading}
                  />
                }
                label="Remember me"
              />
              <Link 
                component={RouterLink} 
                to="/forgot-password" 
                variant="body2"
                sx={{ textDecoration: 'none' }}
              >
                Forgot password?
              </Link>
            </Box>
            <Button
              type="submit"
              fullWidth
              variant="contained"
              sx={{ mt: 3, mb: 2, py: 1.5 }}
              disabled={loading || !formData.email || !formData.password}
              data-testid="login-button"
            >
              {loading ? (
                <CircularProgress size={24} color="inherit" />
              ) : (
                'Sign In'
              )}
            </Button>
            <Box sx={{ textAlign: 'center', mt: 2 }}>
              <Typography variant="body2" color="text.secondary">
                Don't have an account?{' '}
                <Link 
                  component={RouterLink} 
                  to="/register" 
                  variant="body2"
                  sx={{ textDecoration: 'none', fontWeight: 'medium' }}
                >
                  Sign up
                </Link>
              </Typography>
            </Box>
          </Box>
        </Paper>
      </Box>
        <Box sx={{ mt: 3, textAlign: 'center' }}>
          <Typography variant="caption" color="textSecondary">
            {new Date().getFullYear()} PriceWatcher. All rights reserved.
          </Typography>
        </Box>
      </Box>
    </Container>
  );
};

export default Login;
