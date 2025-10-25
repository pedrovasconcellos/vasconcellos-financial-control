import React, { useMemo } from 'react';
import { NavLink, useLocation } from 'react-router-dom';
import {
  AppBar,
  Toolbar,
  Typography,
  Box,
  IconButton,
  Drawer,
  List,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Button
} from '@mui/material';
import DashboardIcon from '@mui/icons-material/Dashboard';
import AccountBalanceIcon from '@mui/icons-material/AccountBalance';
import ReceiptLongIcon from '@mui/icons-material/ReceiptLong';
import SavingsIcon from '@mui/icons-material/Savings';
import FlagIcon from '@mui/icons-material/Flag';

const drawerWidth = 240;

interface DrawerItem {
  label: string;
  path: string;
  icon: React.ReactNode;
}

interface AppLayoutProps {
  children: React.ReactNode;
  onLogout: () => void;
}

const AppLayout: React.FC<AppLayoutProps> = ({ children, onLogout }) => {
  const location = useLocation();

  const items = useMemo<DrawerItem[]>(
    () => [
      { label: 'Dashboard', path: '/', icon: <DashboardIcon /> },
      { label: 'Accounts', path: '/accounts', icon: <AccountBalanceIcon /> },
      { label: 'Transactions', path: '/transactions', icon: <ReceiptLongIcon /> },
      { label: 'Budgets', path: '/budgets', icon: <SavingsIcon /> },
      { label: 'Goals', path: '/goals', icon: <FlagIcon /> }
    ],
    []
  );

  return (
    <Box sx={{ display: 'flex' }}>
      <AppBar position="fixed" sx={{ zIndex: (theme) => theme.zIndex.drawer + 1 }}>
        <Toolbar>
          <Typography variant="h6" sx={{ flexGrow: 1 }}>
            Finance Control
          </Typography>
          <Button color="inherit" onClick={onLogout}>
            Logout
          </Button>
        </Toolbar>
      </AppBar>
      <Drawer
        variant="permanent"
        sx={{
          width: drawerWidth,
          flexShrink: 0,
          [`& .MuiDrawer-paper`]: { width: drawerWidth, boxSizing: 'border-box' }
        }}
      >
        <Toolbar />
        <Box sx={{ overflow: 'auto' }}>
          <List>
            {items.map((item) => (
              <ListItemButton
                key={item.path}
                component={NavLink}
                to={item.path}
                selected={location.pathname === item.path}
              >
                <ListItemIcon>{item.icon}</ListItemIcon>
                <ListItemText primary={item.label} />
              </ListItemButton>
            ))}
          </List>
        </Box>
      </Drawer>
      <Box component="main" sx={{ flexGrow: 1, p: 3 }}>
        <Toolbar />
        {children}
      </Box>
    </Box>
  );
};

export default AppLayout;
