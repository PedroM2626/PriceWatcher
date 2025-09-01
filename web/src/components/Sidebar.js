import React from 'react';
import { Drawer, List, ListItem, ListItemIcon, ListItemText, Divider, Toolbar, Typography } from '@mui/material';
import { useNavigate, useLocation } from 'react-router-dom';
import DashboardIcon from '@mui/icons-material/Dashboard';
import ShoppingCartIcon from '@mui/icons-material/ShoppingCart';
import NotificationsActiveIcon from '@mui/icons-material/NotificationsActive';
import SettingsIcon from '@mui/icons-material/Settings';
import { styled } from '@mui/material/styles';

const drawerWidth = 240;

const StyledDrawer = styled(Drawer)(({ theme }) => ({
  width: drawerWidth,
  flexShrink: 0,
  '& .MuiDrawer-paper': {
    width: drawerWidth,
    boxSizing: 'border-box',
    backgroundColor: theme.palette.background.paper,
    borderRight: `1px solid ${theme.palette.divider}`,
  },
}));

const Sidebar = ({ mobileOpen, onClose }) => {
  const navigate = useNavigate();
  const location = useLocation();

  const menuItems = [
    { text: 'Dashboard', icon: <DashboardIcon />, path: '/' },
    { text: 'Products', icon: <ShoppingCartIcon />, path: '/products' },
    { text: 'Alerts', icon: <NotificationsActiveIcon />, path: '/alerts' },
    { text: 'Settings', icon: <SettingsIcon />, path: '/settings' },
  ];

  const handleItemClick = (path) => {
    navigate(path);
    if (mobileOpen) {
      onClose();
    }
  };

  const isActive = (path) => {
    return location.pathname === path;
  };

  return (
    <>
      <StyledDrawer
        variant="temporary"
        open={mobileOpen}
        onClose={onClose}
        ModalProps={{
          keepMounted: true, // Better open performance on mobile
        }}
        sx={{
          display: { xs: 'block', sm: 'none' },
          '& .MuiDrawer-paper': { boxSizing: 'border-box', width: drawerWidth },
        }}
      >
        <DrawerContent />
      </StyledDrawer>
      <StyledDrawer
        variant="permanent"
        sx={{
          display: { xs: 'none', sm: 'block' },
          '& .MuiDrawer-paper': { boxSizing: 'border-box', width: drawerWidth },
        }}
        open
      >
        <DrawerContent />
      </StyledDrawer>
    </>
  );

  function DrawerContent() {
    return (
      <>
        <Toolbar>
          <Typography variant="h6" noWrap component="div" sx={{ fontWeight: 600 }}>
            PriceWatcher
          </Typography>
        </Toolbar>
        <Divider />
        <List>
          {menuItems.map((item) => (
            <ListItem
              button
              key={item.text}
              onClick={() => handleItemClick(item.path)}
              sx={{
                backgroundColor: isActive(item.path) ? 'action.selected' : 'transparent',
                '&:hover': {
                  backgroundColor: 'action.hover',
                },
                mb: 0.5,
                mx: 1,
                borderRadius: 1,
              }}
            >
              <ListItemIcon sx={{ minWidth: 40, color: isActive(item.path) ? 'primary.main' : 'inherit' }}>
                {item.icon}
              </ListItemIcon>
              <ListItemText 
                primary={item.text} 
                primaryTypographyProps={{
                  fontWeight: isActive(item.path) ? 'medium' : 'regular',
                }}
              />
            </ListItem>
          ))}
        </List>
      </>
    );
  }
};

export default Sidebar;
