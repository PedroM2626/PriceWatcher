import React, { useState, useEffect } from 'react';
import { 
  Box, 
  Grid, 
  Paper, 
  Typography, 
  Card, 
  CardContent, 
  CardHeader, 
  Divider, 
  CircularProgress,
  Alert
} from '@mui/material';
import { Line } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';
import { productsAPI } from '../services/api';

// Register ChartJS components
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
);

// Default empty data structure
const defaultData = {
  totalProducts: 0,
  activeAlerts: 0,
  priceDrops: 0,
  priceHistory: {
    labels: [],
    datasets: [
      {
        label: 'Average Price',
        data: [],
        borderColor: 'rgb(25, 118, 210)',
        backgroundColor: 'rgba(25, 118, 210, 0.1)',
        tension: 0.3,
        fill: true,
      },
    ],
  },
  recentProducts: [],
};

const StatCard = ({ title, value, subtitle, icon, color }) => (
  <Card sx={{ height: '100%' }}>
    <CardContent>
      <Box display="flex" justifyContent="space-between" alignItems="center">
        <Box>
          <Typography color="textSecondary" variant="subtitle2" gutterBottom>
            {title}
          </Typography>
          <Typography variant="h4">{value}</Typography>
          {subtitle && (
            <Typography variant="caption" color="textSecondary">
              {subtitle}
            </Typography>
          )}
        </Box>
        <Box
          sx={{
            backgroundColor: `${color}.light`,
            borderRadius: '50%',
            width: 56,
            height: 56,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            color: 'white',
          }}
        >
          {icon}
        </Box>
      </Box>
    </CardContent>
  </Card>
);

const Dashboard = () => {
  const [data, setData] = useState(defaultData);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchDashboardData = async () => {
      try {
        setLoading(true);
        
        // Fetch products with their price history
        const productsResponse = await productsAPI.getAll({ limit: 5, includeHistory: true });
        const products = productsResponse.data.items || [];
        
        // Calculate stats
        const totalProducts = products.length;
        const activeAlerts = products.reduce((count, product) => 
          count + (product.alerts?.length || 0), 0);
        
        const priceDrops = products.filter(
          product => product.priceHistory && product.priceHistory.length >= 2 &&
          product.priceHistory[0].price < product.priceHistory[1].price
        ).length;
        
        // Prepare price history data
        const priceHistory = {
          labels: [],
          datasets: [{
            label: 'Average Price',
            data: [],
            borderColor: 'rgb(25, 118, 210)',
            backgroundColor: 'rgba(25, 118, 210, 0.1)',
            tension: 0.3,
            fill: true,
          }]
        };
        
        // Process price history (example: last 7 days)
        if (products.length > 0) {
          const last7Days = [...Array(7)].map((_, i) => {
            const d = new Date();
            d.setDate(d.getDate() - (6 - i));
            return d.toISOString().split('T')[0];
          });
          
          priceHistory.labels = last7Days.map(date => 
            new Date(date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
          );
          
          // This is a simplified example - in a real app, you'd aggregate prices by date
          priceHistory.datasets[0].data = last7Days.map((_, i) => 
            Math.random() * 200 + 100 // Replace with actual data
          );
        }
        
        // Get recent products with price changes
        const recentProducts = products.map(product => {
          const priceHistory = product.priceHistory || [];
          const currentPrice = priceHistory[0]?.price || 0;
          const previousPrice = priceHistory[1]?.price || currentPrice;
          const change = currentPrice - previousPrice;
          const changePercent = previousPrice > 0 ? 
            ((change / previousPrice) * 100).toFixed(1) : 0;
            
          return {
            id: product.id,
            name: product.name,
            price: currentPrice,
            change,
            changePercent: parseFloat(changePercent)
          };
        });
        
        setData({
          totalProducts,
          activeAlerts,
          priceDrops,
          priceHistory,
          recentProducts
        });
        
      } catch (err) {
        console.error('Failed to fetch dashboard data:', err);
        setError(err.message || 'Failed to load dashboard data');
      } finally {
        setLoading(false);
      }
    };
    
    fetchDashboardData();
  }, []);
  
  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="60vh">
        <CircularProgress />
      </Box>
    );
  }
  
  if (error) {
    return (
      <Box my={4}>
        <Alert severity="error">
          {error}
        </Alert>
      </Box>
    );
  }
  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Dashboard
      </Typography>
      
      {/* Stats Cards */}
      <Grid container spacing={3} sx={{ mb: 3 }}>
        <Grid item xs={12} sm={6} md={3}>
          <StatCard
            title="Total Products"
            value={data.totalProducts}
            icon={<i className="fas fa-box" style={{ fontSize: 24 }} />}
            color="primary"
          />
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <StatCard
            title="Active Alerts"
            value={data.activeAlerts}
            icon={<i className="fas fa-bell" style={{ fontSize: 24 }} />}
            color="warning"
          />
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <StatCard
            title="Price Drops"
            value={data.priceDrops}
            subtitle="This month"
            icon={<i className="fas fa-arrow-down" style={{ fontSize: 24 }} />}
            color="success"
          />
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <StatCard
            title="Avg. Savings"
            value="$24.50"
            subtitle="Per product"
            icon={<i className="fas fa-piggy-bank" style={{ fontSize: 24 }} />}
            color="info"
          />
        </Grid>
      </Grid>

      {/* Price History Chart */}
      <Grid container spacing={3}>
        <Grid item xs={12} md={8}>
          <Paper sx={{ p: 2, height: '100%' }}>
            <Typography variant="h6" gutterBottom>
              Price History (Last 6 Months)
            </Typography>
            <Box sx={{ height: 300 }}>
              {data.priceHistory.labels.length > 0 ? (
                <Line
                  data={data.priceHistory}
                  options={{
                    responsive: true,
                    maintainAspectRatio: false,
                    plugins: {
                      legend: {
                        position: 'top',
                      },
                      tooltip: {
                        callbacks: {
                          label: (context) => {
                            return `$${context.parsed.y.toFixed(2)}`;
                          }
                        }
                      }
                    },
                    scales: {
                      y: {
                        beginAtZero: false,
                        grid: {
                          drawBorder: false,
                        },
                        ticks: {
                          callback: (value) => `$${value}`
                        }
                      },
                      x: {
                        grid: {
                          display: false,
                        },
                      },
                    },
                  }}
                />
              ) : (
                <Box display="flex" justifyContent="center" alignItems="center" height={300}>
                  <Typography color="textSecondary">No price history available</Typography>
                </Box>
              )}
            </Box>
          </Paper>
        </Grid>
        
        {/* Recent Price Changes */}
        <Grid item xs={12} md={4}>
          <Paper sx={{ p: 2, height: '100%' }}>
            <Typography variant="h6" gutterBottom>
              Recent Price Changes
            </Typography>
            <Box>
              {data.recentProducts.length > 0 ? (
            data.recentProducts.map((product) => (
                <Box key={product.id} sx={{ mb: 2, pb: 2, borderBottom: '1px solid', borderColor: 'divider' }}>
                  <Box display="flex" justifyContent="space-between" alignItems="center">
                    <Typography variant="subtitle1">{product.name}</Typography>
                    <Typography 
                      variant="body2" 
                      color={product.change < 0 ? 'success.main' : product.change > 0 ? 'error.main' : 'text.primary'}
                    >
                      {product.change < 0 ? '↓ ' : product.change > 0 ? '↑ ' : ''}
                      ${Math.abs(product.change).toFixed(2)} ({product.changePercent}%)
                    </Typography>
                  </Box>
                  <Typography variant="h6">
                    ${product.price.toFixed(2)}
                  </Typography>
                </Box>
                ))
          ) : (
            <Typography color="textSecondary" align="center" sx={{ mt: 2 }}>
              No recent products with price changes
            </Typography>
          )}
            </Box>
          </Paper>
        </Grid>
      </Grid>
    </Box>
  );
};

export default Dashboard;
