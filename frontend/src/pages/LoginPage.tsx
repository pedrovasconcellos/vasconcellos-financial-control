import { useState } from 'react';
import {
  Box,
  Button,
  Container,
  TextField,
  Typography,
  Alert,
  Paper,
  Stack
} from '@mui/material';
import { useNavigate } from 'react-router-dom';

import { api } from '../services/api';
import { useAuth } from '../hooks/useAuth';
import type { AuthTokens } from '../providers/AuthProvider';

const LoginPage = () => {
  const { setTokens } = useAuth();
  const navigate = useNavigate();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    setLoading(true);
    setError(null);
    try {
      const { data } = await api.post<AuthTokens>('/auth/login', {
        username: email,
        password
      });
      setTokens(data);
      navigate('/', { replace: true });
    } catch (err) {
      setError('Authentication failed. Please verify your credentials.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Container component="main" maxWidth="xs">
      <Paper elevation={3} sx={{ mt: 12, p: 4 }}>
        <Typography component="h1" variant="h5" align="center" gutterBottom>
          Finance Control
        </Typography>
        <Typography variant="body2" color="text.secondary" align="center" mb={2}>
          Secure access to your personal finance control panel.
        </Typography>
        <Box component="form" onSubmit={handleSubmit}>
          <Stack spacing={2}>
            <TextField
              label="Email"
              type="email"
              value={email}
              onChange={(event) => setEmail(event.target.value)}
              required
              fullWidth
            />
            <TextField
              label="Password"
              type="password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              required
              fullWidth
            />
            {error && <Alert severity="error">{error}</Alert>}
            <Button type="submit" variant="contained" fullWidth disabled={loading}>
              {loading ? 'Signing in...' : 'Sign in'}
            </Button>
          </Stack>
        </Box>
      </Paper>
    </Container>
  );
};

export default LoginPage;
