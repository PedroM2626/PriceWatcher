import React, { useState } from 'react';
import { useNavigate, Link as RouterLink } from 'react-router-dom';
import { Box, Container, Paper, Typography, TextField, Button, Alert, CircularProgress, Link } from '@mui/material';
import { supabase } from '../services/supabase';

const Register = () => {
  const navigate = useNavigate();
  const [form, setForm] = useState({ email: '', password: '' });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [info, setInfo] = useState(null);

  const onChange = (e) => setForm({ ...form, [e.target.name]: e.target.value });

  const onSubmit = async (e) => {
    e.preventDefault();
    setError(null); setInfo(null); setLoading(true);
    try {
      const { error: err } = await supabase.auth.signUp({ email: form.email, password: form.password });
      if (err) throw err;
      setInfo('Check your email to confirm your account.');
      setTimeout(() => navigate('/login'), 1500);
    } catch (err) {
      setError(err.message || 'Failed to register');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Container component="main" maxWidth="xs">
      <Box sx={{ mt: 8, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
        <Paper elevation={3} sx={{ p: 4, width: '100%', mt: 4 }}>
          <Typography component="h1" variant="h5" gutterBottom>
            Create your account
          </Typography>
          {error && <Alert severity="error" sx={{ mb:2 }}>{error}</Alert>}
          {info && <Alert severity="success" sx={{ mb:2 }}>{info}</Alert>}
          <Box component="form" onSubmit={onSubmit} noValidate sx={{ mt: 1 }}>
            <TextField name="email" id="email" label="Email Address" fullWidth required margin="normal" value={form.email} onChange={onChange} />
            <TextField name="password" id="password" label="Password" type="password" fullWidth required margin="normal" value={form.password} onChange={onChange} />
            <Button type="submit" fullWidth variant="contained" sx={{ mt: 3, mb: 2 }} disabled={loading || !form.email || !form.password}>
              {loading ? <CircularProgress size={24} color="inherit" /> : 'Sign Up'}
            </Button>
            <Box sx={{ textAlign:'center' }}>
              <Typography variant="body2" color="text.secondary">
                Already have an account? <Link component={RouterLink} to="/login">Sign in</Link>
              </Typography>
            </Box>
          </Box>
        </Paper>
      </Box>
    </Container>
  );
};

export default Register;
