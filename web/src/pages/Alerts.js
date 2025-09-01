import React, { useState, useEffect } from 'react';
import {
  Box,
  Card,
  CardContent,
  CardHeader,
  Divider,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TablePagination,
  Typography,
  Switch,
  FormControlLabel,
  TextField,
  MenuItem,
  Select,
  Button,
  Grid,
  Paper,
  Chip,
  IconButton,
  Tooltip,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  DialogContentText,
  Badge,
} from '@mui/material';
import { Add as AddIcon, FilterList, Delete as DeleteIcon, NotificationsActive } from '@mui/icons-material';

// Mock data - replace with API calls
const mockAlerts = [
  {
    id: '1',
    productId: '1',
    productName: 'Smartphone X',
    targetPrice: 849.99,
    currentPrice: 899.99,
    condition: 'below',
    isActive: true,
    createdAt: '2023-05-01T10:30:00Z',
    lastTriggered: null,
  },
  {
    id: '2',
    productId: '2',
    productName: 'Laptop Pro',
    targetPrice: 1199.99,
    currentPrice: 1299.99,
    condition: 'below',
    isActive: true,
    createdAt: '2023-05-05T14:20:00Z',
    lastTriggered: null,
  },
  {
    id: '3',
    productId: '3',
    productName: 'Wireless Earbuds',
    targetPrice: 139.99,
    currentPrice: 149.99,
    condition: 'below',
    isActive: false,
    createdAt: '2023-05-10T09:15:00Z',
    lastTriggered: '2023-05-12T11:45:00Z',
  },
];

const Alerts = () => {
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(5);
  const [alerts, setAlerts] = useState([]);
  const [filter, setFilter] = useState('all');
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [selectedAlert, setSelectedAlert] = useState(null);
  const [showAddDialog, setShowAddDialog] = useState(false);
  const [newAlert, setNewAlert] = useState({
    productId: '',
    targetPrice: '',
    condition: 'below',
  });

  // Load alerts (replace with API call)
  useEffect(() => {
    setAlerts(mockAlerts);
  }, []);

  const handleChangePage = (event, newPage) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const handleFilterChange = (event) => {
    setFilter(event.target.value);
    setPage(0);
  };

  const handleToggleAlert = async (alertId) => {
    // TODO: Implement API call to toggle alert status
    setAlerts(alerts.map(alert => 
      alert.id === alertId ? { ...alert, isActive: !alert.isActive } : alert
    ));
  };

  const handleDeleteClick = (alert) => {
    setSelectedAlert(alert);
    setDeleteDialogOpen(true);
  };

  const handleDeleteConfirm = () => {
    // TODO: Implement API call to delete alert
    setAlerts(alerts.filter(alert => alert.id !== selectedAlert.id));
    setDeleteDialogOpen(false);
    setSelectedAlert(null);
  };

  const handleAddAlert = () => {
    // TODO: Implement API call to add alert
    const newAlertObj = {
      id: `alert-${Date.now()}`,
      productId: newAlert.productId,
      productName: 'New Product', // This would come from the product lookup
      targetPrice: parseFloat(newAlert.targetPrice),
      currentPrice: 0, // This would come from the product lookup
      condition: newAlert.condition,
      isActive: true,
      createdAt: new Date().toISOString(),
      lastTriggered: null,
    };
    
    setAlerts([...alerts, newAlertObj]);
    setShowAddDialog(false);
    setNewAlert({
      productId: '',
      targetPrice: '',
      condition: 'below',
    });
  };

  const filteredAlerts = alerts.filter(alert => {
    if (filter === 'all') return true;
    if (filter === 'active') return alert.isActive;
    if (filter === 'inactive') return !alert.isActive;
    if (filter === 'triggered') return alert.lastTriggered !== null;
    return true;
  });

  const getStatusChip = (alert) => {
    if (alert.lastTriggered) {
      return <Chip label="Triggered" color="success" size="small" />;
    }
    return alert.isActive ? (
      <Chip label="Active" color="primary" size="small" />
    ) : (
      <Chip label="Inactive" variant="outlined" size="small" />
    );
  };

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Price Alerts</Typography>
        <Button
          variant="contained"
          color="primary"
          startIcon={<AddIcon />}
          onClick={() => setShowAddDialog(true)}
        >
          Create Alert
        </Button>
      </Box>

      <Card>
        <CardHeader
          title={
            <Box display="flex" alignItems="center">
              <FilterList sx={{ mr: 1, color: 'text.secondary' }} />
              <Typography variant="subtitle1" color="textSecondary">
                Filter Alerts
              </Typography>
              <Select
                value={filter}
                onChange={handleFilterChange}
                variant="outlined"
                size="small"
                sx={{ ml: 2, minWidth: 150 }}
              >
                <MenuItem value="all">All Alerts</MenuItem>
                <MenuItem value="active">Active</MenuItem>
                <MenuItem value="inactive">Inactive</MenuItem>
                <MenuItem value="triggered">Triggered</MenuItem>
              </Select>
            </Box>
          }
          sx={{ pb: 1 }}
        />
        <Divider />
        <CardContent>
          <TableContainer>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Product</TableCell>
                  <TableCell>Alert Condition</TableCell>
                  <TableCell>Current Price</TableCell>
                  <TableCell>Status</TableCell>
                  <TableCell>Created</TableCell>
                  <TableCell align="right">Actions</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {filteredAlerts.length > 0 ? (
                  filteredAlerts
                    .slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage)
                    .map((alert) => (
                      <TableRow key={alert.id} hover>
                        <TableCell>
                          <Typography variant="subtitle2">{alert.productName}</Typography>
                          <Typography variant="caption" color="textSecondary">
                            ID: {alert.productId}
                          </Typography>
                        </TableCell>
                        <TableCell>
                          <Box display="flex" alignItems="center">
                            <Typography variant="body2">
                              Price {alert.condition} ${alert.targetPrice.toFixed(2)}
                            </Typography>
                          </Box>
                        </TableCell>
                        <TableCell>
                          <Typography variant="body2" fontWeight="bold">
                            ${alert.currentPrice.toFixed(2)}
                          </Typography>
                        </TableCell>
                        <TableCell>
                          {getStatusChip(alert)}
                        </TableCell>
                        <TableCell>
                          <Typography variant="body2">
                            {new Date(alert.createdAt).toLocaleDateString()}
                          </Typography>
                          {alert.lastTriggered && (
                            <Typography variant="caption" color="success.main" display="block">
                              Triggered: {new Date(alert.lastTriggered).toLocaleString()}
                            </Typography>
                          )}
                        </TableCell>
                        <TableCell align="right">
                          <Tooltip title={alert.isActive ? 'Deactivate alert' : 'Activate alert'}>
                            <Switch
                              checked={alert.isActive}
                              onChange={() => handleToggleAlert(alert.id)}
                              color="primary"
                              size="small"
                            />
                          </Tooltip>
                          <Tooltip title="Delete alert">
                            <IconButton
                              size="small"
                              onClick={() => handleDeleteClick(alert)}
                              color="error"
                            >
                              <DeleteIcon fontSize="small" />
                            </IconButton>
                          </Tooltip>
                        </TableCell>
                      </TableRow>
                    ))
                ) : (
                  <TableRow>
                    <TableCell colSpan={6} align="center" sx={{ py: 4 }}>
                      <Box display="flex" flexDirection="column" alignItems="center">
                        <NotificationsActive sx={{ fontSize: 48, color: 'text.disabled', mb: 2 }} />
                        <Typography variant="subtitle1" color="textSecondary">
                          No price alerts found
                        </Typography>
                        <Typography variant="body2" color="textSecondary" sx={{ mt: 1, mb: 2 }}>
                          {filter === 'all' 
                            ? 'Create your first price alert to get started.'
                            : 'No alerts match the current filter.'}
                        </Typography>
                        {filter !== 'all' && (
                          <Button
                            variant="outlined"
                            color="primary"
                            onClick={() => setFilter('all')}
                            size="small"
                          >
                            Clear Filter
                          </Button>
                        )}
                      </Box>
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </TableContainer>
          
          {filteredAlerts.length > 0 && (
            <TablePagination
              rowsPerPageOptions={[5, 10, 25]}
              component="div"
              count={filteredAlerts.length}
              rowsPerPage={rowsPerPage}
              page={page}
              onPageChange={handleChangePage}
              onRowsPerPageChange={handleChangeRowsPerPage}
            />
          )}
        </CardContent>
      </Card>

      {/* Delete Confirmation Dialog */}
      <Dialog
        open={deleteDialogOpen}
        onClose={() => setDeleteDialogOpen(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>Delete Alert</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Are you sure you want to delete the price alert for "{selectedAlert?.productName}"?
            This action cannot be undone.
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialogOpen(false)}>Cancel</Button>
          <Button onClick={handleDeleteConfirm} color="error" variant="contained">
            Delete
          </Button>
        </DialogActions>
      </Dialog>

      {/* Add Alert Dialog */}
      <Dialog
        open={showAddDialog}
        onClose={() => setShowAddDialog(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>Create New Price Alert</DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Product ID"
                value={newAlert.productId}
                onChange={(e) => setNewAlert({ ...newAlert, productId: e.target.value })}
                margin="normal"
                variant="outlined"
                required
                helperText="Enter the product ID or search for a product"
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Target Price"
                type="number"
                value={newAlert.targetPrice}
                onChange={(e) => setNewAlert({ ...newAlert, targetPrice: e.target.value })}
                margin="normal"
                variant="outlined"
                required
                inputProps={{ min: 0, step: 0.01 }}
                InputProps={{
                  startAdornment: <span style={{ marginRight: 8 }}>$</span>,
                }}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                select
                label="Alert When Price Is"
                value={newAlert.condition}
                onChange={(e) => setNewAlert({ ...newAlert, condition: e.target.value })}
                margin="normal"
                variant="outlined"
                required
              >
                <MenuItem value="below">Below Target</MenuItem>
                <MenuItem value="above">Above Target</MenuItem>
                <MenuItem value="equal">Equal To Target</MenuItem>
              </TextField>
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowAddDialog(false)}>Cancel</Button>
          <Button 
            onClick={handleAddAlert} 
            variant="contained"
            disabled={!newAlert.productId || !newAlert.targetPrice}
          >
            Create Alert
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default Alerts;
