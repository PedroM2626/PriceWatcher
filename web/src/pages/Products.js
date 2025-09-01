import React, { useState, useEffect } from 'react';
import {
  Box,
  Button,
  Card,
  CardContent,
  CardHeader,
  Divider,
  Grid,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TablePagination,
  TextField,
  Typography,
  IconButton,
  Menu,
  MenuItem,
  Chip,
  Avatar,
  Tooltip,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  DialogContentText,
} from '@mui/material';
import { Add as AddIcon, MoreVert as MoreVertIcon, Link as LinkIcon } from '@mui/icons-material';
import { format } from 'date-fns';
import { Link } from 'react-router-dom';

// Mock data - replace with API calls
const mockProducts = [
  {
    id: '1',
    name: 'Smartphone X',
    currentPrice: 899.99,
    originalPrice: 949.99,
    image: 'https://via.placeholder.com/60',
    url: 'https://example.com/product/smartphone-x',
    website: 'example.com',
    lastChecked: '2023-05-15T14:30:00Z',
    availability: 'in_stock',
  },
  {
    id: '2',
    name: 'Laptop Pro',
    currentPrice: 1299.99,
    originalPrice: 1299.99,
    image: 'https://via.placeholder.com/60',
    url: 'https://example.com/product/laptop-pro',
    website: 'example.com',
    lastChecked: '2023-05-16T09:15:00Z',
    availability: 'in_stock',
  },
  {
    id: '3',
    name: 'Wireless Earbuds',
    currentPrice: 149.99,
    originalPrice: 169.99,
    image: 'https://via.placeholder.com/60',
    url: 'https://example.com/product/wireless-earbuds',
    website: 'example.com',
    lastChecked: '2023-05-14T16:45:00Z',
    availability: 'low_stock',
  },
];

const Products = () => {
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [searchTerm, setSearchTerm] = useState('');
  const [anchorEl, setAnchorEl] = useState(null);
  const [selectedProduct, setSelectedProduct] = useState(null);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [addDialogOpen, setAddDialogOpen] = useState(false);
  const [newProductUrl, setNewProductUrl] = useState('');
  const [products, setProducts] = useState([]);

  // Load products (replace with API call)
  useEffect(() => {
    // Simulate API call
    setProducts(mockProducts);
  }, []);

  const handleMenuOpen = (event, product) => {
    setAnchorEl(event.currentTarget);
    setSelectedProduct(product);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
    setSelectedProduct(null);
  };

  const handleDeleteClick = () => {
    setDeleteDialogOpen(true);
    handleMenuClose();
  };

  const handleDeleteConfirm = () => {
    // TODO: Implement delete functionality
    console.log('Deleting product:', selectedProduct.id);
    setDeleteDialogOpen(false);
  };

  const handleAddProduct = () => {
    // TODO: Implement add product functionality
    console.log('Adding product with URL:', newProductUrl);
    setAddDialogOpen(false);
    setNewProductUrl('');
  };

  const handleChangePage = (event, newPage) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const filteredProducts = products.filter((product) =>
    product.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    product.website.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const getStatusColor = (status) => {
    switch (status) {
      case 'in_stock':
        return 'success';
      case 'low_stock':
        return 'warning';
      case 'out_of_stock':
        return 'error';
      default:
        return 'default';
    }
  };

  const getStatusLabel = (status) => {
    switch (status) {
      case 'in_stock':
        return 'In Stock';
      case 'low_stock':
        return 'Low Stock';
      case 'out_of_stock':
        return 'Out of Stock';
      default:
        return 'Unknown';
    }
  };

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Tracked Products</Typography>
        <Button
          variant="contained"
          color="primary"
          startIcon={<AddIcon />}
          onClick={() => setAddDialogOpen(true)}
        >
          Add Product
        </Button>
      </Box>

      <Card>
        <CardHeader
          title={
            <TextField
              variant="outlined"
              placeholder="Search products..."
              size="small"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              sx={{ width: 300 }}
              InputProps={{
                startAdornment: <i className="fas fa-search" style={{ marginRight: 8, color: 'rgba(0, 0, 0, 0.54)' }} />,
              }}
            />
          }
        />
        <Divider />
        <CardContent>
          <TableContainer>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Product</TableCell>
                  <TableCell>Price</TableCell>
                  <TableCell>Status</TableCell>
                  <TableCell>Last Checked</TableCell>
                  <TableCell>Website</TableCell>
                  <TableCell align="right">Actions</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {filteredProducts.length > 0 ? (
                  filteredProducts
                    .slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage)
                    .map((product) => (
                      <TableRow key={product.id} hover>
                        <TableCell>
                          <Box display="flex" alignItems="center">
                            <Avatar
                              src={product.image}
                              alt={product.name}
                              sx={{ width: 40, height: 40, mr: 2 }}
                            />
                            <Box>
                              <Typography variant="subtitle2">{product.name}</Typography>
                              <Typography variant="caption" color="textSecondary">
                                ID: {product.id}
                              </Typography>
                            </Box>
                          </Box>
                        </TableCell>
                        <TableCell>
                          <Typography variant="body2" fontWeight="bold">
                            ${product.currentPrice.toFixed(2)}
                          </Typography>
                          {product.currentPrice < product.originalPrice && (
                            <Typography variant="caption" color="success.main">
                              ${(product.originalPrice - product.currentPrice).toFixed(2)} off
                            </Typography>
                          )}
                        </TableCell>
                        <TableCell>
                          <Chip
                            label={getStatusLabel(product.availability)}
                            color={getStatusColor(product.availability)}
                            size="small"
                          />
                        </TableCell>
                        <TableCell>
                          {format(new Date(product.lastChecked), 'MMM d, yyyy HH:mm')}
                        </TableCell>
                        <TableCell>
                          <Chip
                            label={product.website}
                            variant="outlined"
                            size="small"
                            avatar={<Avatar src={`https://www.google.com/s2/favicons?domain=${product.website}`} />}
                          />
                        </TableCell>
                        <TableCell align="right">
                          <IconButton
                            size="small"
                            onClick={(e) => handleMenuOpen(e, product)}
                          >
                            <MoreVertIcon />
                          </IconButton>
                        </TableCell>
                      </TableRow>
                    ))
                ) : (
                  <TableRow>
                    <TableCell colSpan={6} align="center" sx={{ py: 4 }}>
                      <Typography variant="subtitle1" color="textSecondary">
                        No products found. Try adjusting your search or add a new product.
                      </Typography>
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </TableContainer>
          
          <TablePagination
            rowsPerPageOptions={[5, 10, 25]}
            component="div"
            count={filteredProducts.length}
            rowsPerPage={rowsPerPage}
            page={page}
            onPageChange={handleChangePage}
            onRowsPerPageChange={handleChangeRowsPerPage}
          />
        </CardContent>
      </Card>

      {/* Product Actions Menu */}
      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={handleMenuClose}
        anchorOrigin={{
          vertical: 'top',
          horizontal: 'right',
        }}
        transformOrigin={{
          vertical: 'top',
          horizontal: 'right',
        }}
      >
        <MenuItem component={Link} to={`/products/${selectedProduct?.id}`}>
          View Details
        </MenuItem>
        <MenuItem onClick={handleMenuClose}>
          Check Price Now
        </MenuItem>
        <MenuItem onClick={handleMenuClose}>
          Set Price Alert
        </MenuItem>
        <Divider />
        <MenuItem onClick={handleDeleteClick} sx={{ color: 'error.main' }}>
          Remove Product
        </MenuItem>
      </Menu>

      {/* Delete Confirmation Dialog */}
      <Dialog
        open={deleteDialogOpen}
        onClose={() => setDeleteDialogOpen(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>Remove Product</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Are you sure you want to remove "{selectedProduct?.name}" from your tracked products?
            This action cannot be undone.
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialogOpen(false)}>Cancel</Button>
          <Button onClick={handleDeleteConfirm} color="error" variant="contained">
            Remove
          </Button>
        </DialogActions>
      </Dialog>

      {/* Add Product Dialog */}
      <Dialog
        open={addDialogOpen}
        onClose={() => setAddDialogOpen(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>Add New Product</DialogTitle>
        <DialogContent>
          <DialogContentText sx={{ mb: 2 }}>
            Enter the product URL from a supported retailer to start tracking its price.
          </DialogContentText>
          <TextField
            autoFocus
            margin="dense"
            label="Product URL"
            type="url"
            fullWidth
            variant="outlined"
            value={newProductUrl}
            onChange={(e) => setNewProductUrl(e.target.value)}
            placeholder="https://www.example.com/product/123"
            InputProps={{
              startAdornment: <LinkIcon sx={{ mr: 1, color: 'text.secondary' }} />,
            }}
          />
          <Box mt={2}>
            <Typography variant="caption" color="textSecondary">
              Supported retailers: Amazon, Mercado Livre, and more.
            </Typography>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setAddDialogOpen(false)}>Cancel</Button>
          <Button 
            onClick={handleAddProduct} 
            variant="contained"
            disabled={!newProductUrl}
          >
            Track Product
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default Products;
