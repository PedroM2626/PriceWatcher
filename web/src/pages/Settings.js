import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Card,
  CardContent,
  CardHeader,
  Divider,
  Grid,
  TextField,
  Button,
  FormControlLabel,
  Switch,
  Tabs,
  Tab,
  Alert,
  Snackbar,
  CircularProgress,
  InputAdornment,
  IconButton,
  List,
  ListItem,
  ListItemText,
  ListItemSecondaryAction,
  MenuItem,
  Select,
  FormControl,
  InputLabel,
  Paper,
} from '@mui/material';
import {
  Email as EmailIcon,
  Telegram as TelegramIcon,
  Save as SaveIcon,
  Visibility,
  VisibilityOff,
  Delete as DeleteIcon,
  Add as AddIcon,
} from '@mui/icons-material';

// Mock data - replace with API calls
const mockUser = {
  email: 'user@example.com',
  name: 'John Doe',
  notificationPreferences: {
    email: true,
    telegram: false,
    priceDrop: true,
    backInStock: true,
    weeklyDigest: true,
  },
  telegramChatId: null,
  currency: 'USD',
  timezone: 'America/New_York',
  theme: 'light',
};

const currencies = [
  { value: 'USD', label: 'US Dollar ($)' },
  { value: 'EUR', label: 'Euro (€)' },
  { value: 'GBP', label: 'British Pound (£)' },
  { value: 'BRL', label: 'Brazilian Real (R$)' },
  { value: 'JPY', label: 'Japanese Yen (¥)' },
];

const timezones = [
  'America/New_York',
  'America/Chicago',
  'America/Denver',
  'America/Los_Angeles',
  'Europe/London',
  'Europe/Paris',
  'Asia/Tokyo',
  'Australia/Sydney',
];

const Settings = () => {
  const [tabValue, setTabValue] = useState(0);
  const [user, setUser] = useState(mockUser);
  const [loading, setLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [snackbar, setSnackbar] = useState({
    open: false,
    message: '',
    severity: 'success',
  });
  const [newPassword, setNewPassword] = useState({
    current: '',
    new: '',
    confirm: '',
  });
  const [telegramConnected, setTelegramConnected] = useState(!!mockUser.telegramChatId);

  const handleTabChange = (event, newValue) => {
    setTabValue(newValue);
  };

  const handleNotificationChange = (event) => {
    const { name, checked } = event.target;
    setUser({
      ...user,
      notificationPreferences: {
        ...user.notificationPreferences,
        [name]: checked,
      },
    });
  };

  const handleInputChange = (event) => {
    const { name, value } = event.target;
    setUser({
      ...user,
      [name]: value,
    });
  };

  const handlePasswordChange = (event) => {
    const { name, value } = event.target;
    setNewPassword({
      ...newPassword,
      [name]: value,
    });
  };

  const handleSave = async () => {
    setLoading(true);
    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      setSnackbar({
        open: true,
        message: 'Settings saved successfully!',
        severity: 'success',
      });
    } catch (error) {
      setSnackbar({
        open: true,
        message: 'Failed to save settings. Please try again.',
        severity: 'error',
      });
    } finally {
      setLoading(false);
    }
  };

  const handleChangePassword = async () => {
    if (newPassword.new !== newPassword.confirm) {
      setSnackbar({
        open: true,
        message: 'New passwords do not match!',
        severity: 'error',
      });
      return;
    }
    
    setLoading(true);
    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      setSnackbar({
        open: true,
        message: 'Password changed successfully!',
        severity: 'success',
      });
      
      // Clear password fields
      setNewPassword({
        current: '',
        new: '',
        confirm: '',
      });
    } catch (error) {
      setSnackbar({
        open: true,
        message: 'Failed to change password. Please try again.',
        severity: 'error',
      });
    } finally {
      setLoading(false);
    }
  };

  const handleConnectTelegram = () => {
    // In a real app, this would open a Telegram bot link
    setSnackbar({
      open: true,
      message: 'Please check your Telegram app to complete the connection.',
      severity: 'info',
    });
    
    // Simulate successful connection
    setTimeout(() => {
      setTelegramConnected(true);
      setSnackbar({
        open: true,
        message: 'Telegram connected successfully!',
        severity: 'success',
      });
    }, 2000);
  };

  const handleDisconnectTelegram = () => {
    setTelegramConnected(false);
    setSnackbar({
      open: true,
      message: 'Telegram disconnected successfully!',
      severity: 'success',
    });
  };

  const handleCloseSnackbar = () => {
    setSnackbar({
      ...snackbar,
      open: false,
    });
  };

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Settings
      </Typography>
      
      <Tabs
        value={tabValue}
        onChange={handleTabChange}
        indicatorColor="primary"
        textColor="primary"
        variant="scrollable"
        scrollButtons="auto"
        sx={{ mb: 3 }}
      >
        <Tab label="Account" />
        <Tab label="Notifications" />
        <Tab label="Security" />
        <Tab label="Preferences" />
      </Tabs>
      
      {/* Account Settings */}
      {tabValue === 0 && (
        <Card>
          <CardHeader title="Account Information" />
          <Divider />
          <CardContent>
            <Grid container spacing={3}>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Name"
                  name="name"
                  value={user.name}
                  onChange={handleInputChange}
                  margin="normal"
                  variant="outlined"
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Email Address"
                  name="email"
                  type="email"
                  value={user.email}
                  onChange={handleInputChange}
                  margin="normal"
                  variant="outlined"
                  disabled
                />
              </Grid>
              <Grid item xs={12}>
                <Button
                  variant="contained"
                  color="primary"
                  startIcon={loading ? <CircularProgress size={20} color="inherit" /> : <SaveIcon />}
                  onClick={handleSave}
                  disabled={loading}
                >
                  Save Changes
                </Button>
              </Grid>
            </Grid>
          </CardContent>
        </Card>
      )}
      
      {/* Notification Settings */}
      {tabValue === 1 && (
        <Card>
          <CardHeader title="Notification Preferences" />
          <Divider />
          <CardContent>
            <Typography variant="h6" gutterBottom>
              Notification Methods
            </Typography>
            <List>
              <ListItem>
                <ListItemText
                  primary="Email Notifications"
                  secondary="Receive notifications via email"
                  primaryTypographyProps={{
                    variant: 'subtitle1',
                  }}
                />
                <ListItemSecondaryAction>
                  <FormControlLabel
                    control={
                      <Switch
                        checked={user.notificationPreferences.email}
                        onChange={handleNotificationChange}
                        name="email"
                        color="primary"
                      />
                    }
                    label={user.notificationPreferences.email ? 'On' : 'Off'}
                  />
                </ListItemSecondaryAction>
              </ListItem>
              
              <Divider component="li" />
              
              <ListItem>
                <ListItemText
                  primary="Telegram Notifications"
                  secondary={
                    telegramConnected 
                      ? `Connected to Telegram (${user.telegramChatId || 'Chat ID: 123456789'})` 
                      : "Connect your Telegram account to receive instant notifications"
                  }
                  primaryTypographyProps={{
                    variant: 'subtitle1',
                  }}
                />
                <ListItemSecondaryAction>
                  {telegramConnected ? (
                    <Button
                      variant="outlined"
                      color="error"
                      size="small"
                      startIcon={<DeleteIcon />}
                      onClick={handleDisconnectTelegram}
                    >
                      Disconnect
                    </Button>
                  ) : (
                    <Button
                      variant="contained"
                      color="primary"
                      size="small"
                      startIcon={<TelegramIcon />}
                      onClick={handleConnectTelegram}
                    >
                      Connect Telegram
                    </Button>
                  )}
                </ListItemSecondaryAction>
              </ListItem>
            </List>
            
            <Typography variant="h6" gutterBottom sx={{ mt: 4, mb: 2 }}>
              Notification Types
            </Typography>
            
            <Paper variant="outlined" sx={{ p: 2, mb: 2 }}>
              <FormControlLabel
                control={
                  <Switch
                    checked={user.notificationPreferences.priceDrop}
                    onChange={handleNotificationChange}
                    name="priceDrop"
                    color="primary"
                  />
                }
                label="Price Drop Alerts"
                sx={{ mb: 1, display: 'block' }}
              />
              <Typography variant="body2" color="textSecondary" sx={{ ml: 4, mb: 2 }}>
                Get notified when a product's price drops below your target price
              </Typography>
              
              <FormControlLabel
                control={
                  <Switch
                    checked={user.notificationPreferences.backInStock}
                    onChange={handleNotificationChange}
                    name="backInStock"
                    color="primary"
                  />
                }
                label="Back in Stock Alerts"
                sx={{ mb: 1, display: 'block' }}
              />
              <Typography variant="body2" color="textSecondary" sx={{ ml: 4, mb: 2 }}>
                Get notified when an out-of-stock product becomes available again
              </Typography>
              
              <FormControlLabel
                control={
                  <Switch
                    checked={user.notificationPreferences.weeklyDigest}
                    onChange={handleNotificationChange}
                    name="weeklyDigest"
                    color="primary"
                  />
                }
                label="Weekly Digest"
                sx={{ mb: 1, display: 'block' }}
              />
              <Typography variant="body2" color="textSecondary" sx={{ ml: 4 }}>
                Receive a weekly summary of price changes for your tracked products
              </Typography>
            </Paper>
            
            <Box sx={{ mt: 3 }}>
              <Button
                variant="contained"
                color="primary"
                startIcon={loading ? <CircularProgress size={20} color="inherit" /> : <SaveIcon />}
                onClick={handleSave}
                disabled={loading}
              >
                Save Notification Settings
              </Button>
            </Box>
          </CardContent>
        </Card>
      )}
      
      {/* Security Settings */}
      {tabValue === 2 && (
        <Card>
          <CardHeader title="Security" />
          <Divider />
          <CardContent>
            <Typography variant="h6" gutterBottom>
              Change Password
            </Typography>
            
            <Grid container spacing={3}>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Current Password"
                  name="current"
                  type={showPassword ? 'text' : 'password'}
                  value={newPassword.current}
                  onChange={handlePasswordChange}
                  margin="normal"
                  variant="outlined"
                  InputProps={{
                    endAdornment: (
                      <InputAdornment position="end">
                        <IconButton
                          aria-label="toggle password visibility"
                          onClick={() => setShowPassword(!showPassword)}
                          edge="end"
                        >
                          {showPassword ? <VisibilityOff /> : <Visibility />}
                        </IconButton>
                      </InputAdornment>
                    ),
                  }}
                />
                
                <TextField
                  fullWidth
                  label="New Password"
                  name="new"
                  type={showPassword ? 'text' : 'password'}
                  value={newPassword.new}
                  onChange={handlePasswordChange}
                  margin="normal"
                  variant="outlined"
                  helperText="Must be at least 8 characters long"
                />
                
                <TextField
                  fullWidth
                  label="Confirm New Password"
                  name="confirm"
                  type={showPassword ? 'text' : 'password'}
                  value={newPassword.confirm}
                  onChange={handlePasswordChange}
                  margin="normal"
                  variant="outlined"
                />
                
                <Box sx={{ mt: 2 }}>
                  <Button
                    variant="contained"
                    color="primary"
                    onClick={handleChangePassword}
                    disabled={loading || !newPassword.current || !newPassword.new || !newPassword.confirm}
                    startIcon={loading ? <CircularProgress size={20} color="inherit" /> : <SaveIcon />}
                  >
                    Change Password
                  </Button>
                </Box>
              </Grid>
              
              <Grid item xs={12} md={6}>
                <Paper variant="outlined" sx={{ p: 2, height: '100%' }}>
                  <Typography variant="subtitle1" gutterBottom>
                    Password Requirements
                  </Typography>
                  <Typography variant="body2" color="textSecondary" paragraph>
                    To create a strong password, make sure it includes:
                  </Typography>
                  <ul style={{ paddingLeft: 20, margin: 0 }}>
                    <li>At least 8 characters</li>
                    <li>At least one uppercase letter</li>
                    <li>At least one lowercase letter</li>
                    <li>At least one number</li>
                    <li>At least one special character (e.g., !@#$%^&*)</li>
                  </ul>
                </Paper>
              </Grid>
            </Grid>
            
            <Divider sx={{ my: 4 }} />
            
            <Typography variant="h6" gutterBottom>
              Active Sessions
            </Typography>
            <Paper variant="outlined" sx={{ p: 2, mb: 2 }}>
              <Box display="flex" justifyContent="space-between" alignItems="center">
                <Box>
                  <Typography variant="subtitle2">Current Session</Typography>
                  <Typography variant="body2" color="textSecondary">
                    {new Date().toLocaleString()} • Chrome on Windows 10
                  </Typography>
                  <Typography variant="caption" color="textSecondary">
                    IP: 192.168.1.1 • {navigator.userAgent}
                  </Typography>
                </Box>
                <Button color="error" size="small" disabled>
                  Logout from this device
                </Button>
              </Box>
            </Paper>
            
            <Paper variant="outlined" sx={{ p: 2 }}>
              <Box display="flex" justifyContent="space-between" alignItems="center">
                <Box>
                  <Typography variant="subtitle2">Previous Session</Typography>
                  <Typography variant="body2" color="textSecondary">
                    {new Date(Date.now() - 86400000).toLocaleString()} • Safari on iPhone
                  </Typography>
                  <Typography variant="caption" color="textSecondary">
                    IP: 192.168.1.2 • {navigator.userAgent}
                  </Typography>
                </Box>
                <Button color="error" size="small">
                  Logout from this device
                </Button>
              </Box>
            </Paper>
            
            <Box sx={{ mt: 3 }}>
              <Button
                variant="outlined"
                color="error"
                startIcon={<DeleteIcon />}
                onClick={() => {
                  setSnackbar({
                    open: true,
                    message: 'All other sessions have been logged out.',
                    severity: 'success',
                  });
                }}
              >
                Logout from all other devices
              </Button>
            </Box>
          </CardContent>
        </Card>
      )}
      
      {/* Preferences */}
      {tabValue === 3 && (
        <Card>
          <CardHeader title="Preferences" />
          <Divider />
          <CardContent>
            <Grid container spacing={3}>
              <Grid item xs={12} md={6}>
                <FormControl fullWidth margin="normal">
                  <InputLabel id="currency-label">Currency</InputLabel>
                  <Select
                    labelId="currency-label"
                    id="currency"
                    name="currency"
                    value={user.currency}
                    onChange={handleInputChange}
                    label="Currency"
                  >
                    {currencies.map((currency) => (
                      <MenuItem key={currency.value} value={currency.value}>
                        {currency.label}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
                
                <FormControl fullWidth margin="normal">
                  <InputLabel id="timezone-label">Timezone</InputLabel>
                  <Select
                    labelId="timezone-label"
                    id="timezone"
                    name="timezone"
                    value={user.timezone}
                    onChange={handleInputChange}
                    label="Timezone"
                  >
                    {timezones.map((timezone) => (
                      <MenuItem key={timezone} value={timezone}>
                        {timezone}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
                
                <FormControl fullWidth margin="normal">
                  <InputLabel id="theme-label">Theme</InputLabel>
                  <Select
                    labelId="theme-label"
                    id="theme"
                    name="theme"
                    value={user.theme}
                    onChange={handleInputChange}
                    label="Theme"
                  >
                    <MenuItem value="light">Light</MenuItem>
                    <MenuItem value="dark">Dark</MenuItem>
                    <MenuItem value="system">System Default</MenuItem>
                  </Select>
                </FormControl>
                
                <Box sx={{ mt: 3 }}>
                  <Button
                    variant="contained"
                    color="primary"
                    startIcon={loading ? <CircularProgress size={20} color="inherit" /> : <SaveIcon />}
                    onClick={handleSave}
                    disabled={loading}
                  >
                    Save Preferences
                  </Button>
                </Box>
              </Grid>
              
              <Grid item xs={12} md={6}>
                <Paper variant="outlined" sx={{ p: 2, height: '100%' }}>
                  <Typography variant="subtitle1" gutterBottom>
                    Data & Privacy
                  </Typography>
                  <Typography variant="body2" color="textSecondary" paragraph>
                    We respect your privacy and are committed to protecting your personal data. 
                    You can manage your data and privacy settings here.
                  </Typography>
                  
                  <Button
                    variant="outlined"
                    color="primary"
                    size="small"
                    sx={{ mr: 1, mb: 1 }}
                  >
                    Download My Data
                  </Button>
                  
                  <Button
                    variant="outlined"
                    color="error"
                    size="small"
                    sx={{ mb: 1 }}
                    onClick={() => {
                      setSnackbar({
                        open: true,
                        message: 'Account deletion requested. Please check your email to confirm.',
                        severity: 'warning',
                      });
                    }}
                  >
                    Delete My Account
                  </Button>
                  
                  <Typography variant="caption" color="textSecondary" display="block" sx={{ mt: 2 }}>
                    Last updated: {new Date().toLocaleDateString()}
                  </Typography>
                </Paper>
              </Grid>
            </Grid>
          </CardContent>
        </Card>
      )}
      
      <Snackbar
        open={snackbar.open}
        autoHideDuration={6000}
        onClose={handleCloseSnackbar}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
      >
        <Alert onClose={handleCloseSnackbar} severity={snackbar.severity} sx={{ width: '100%' }}>
          {snackbar.message}
        </Alert>
      </Snackbar>
    </Box>
  );
};

export default Settings;
