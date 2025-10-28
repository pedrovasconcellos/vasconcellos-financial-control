import { createTheme } from '@mui/material/styles';

const theme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: '#1976d2'
    },
    secondary: {
      main: '#ff9800'
    },
    background: {
      default: '#f4f6f8',
      paper: '#ffffff'
    }
  },
  typography: {
    fontFamily: 'Roboto, sans-serif'
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          textTransform: 'none',
          borderRadius: 8
        }
      }
    }
  }
});

export default theme;
