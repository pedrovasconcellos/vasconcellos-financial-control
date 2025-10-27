import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  Box,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Grid,
  TextField,
  MenuItem,
  Paper,
  Typography,
  IconButton,
  Table,
  TableHead,
  TableRow,
  TableCell,
  TableBody,
  Stack,
  Alert,
  Skeleton
} from '@mui/material';
import DeleteIcon from '@mui/icons-material/Delete';

import { api } from '../services/api';
import { currencyOptions, defaultCurrency, CurrencyCode } from '../constants/currencyOptions';

interface Account {
  id: string;
  name: string;
  type: string;
  currency: string;
  description: string;
  balance: number;
}

interface AccountFormState {
  name: string;
  type: string;
  currency: CurrencyCode;
  description: string;
  balance: number;
}

const accountTypes = [
  { label: 'Checking', value: 'checking' },
  { label: 'Savings', value: 'savings' },
  { label: 'Credit', value: 'credit' },
  { label: 'Cash', value: 'cash' }
];

const AccountsPage = () => {
  const queryClient = useQueryClient();
  const [open, setOpen] = useState(false);
  const [form, setForm] = useState<AccountFormState>({
    name: '',
    type: 'checking',
    currency: defaultCurrency,
    description: '',
    balance: 0
  });
  const [formErrors, setFormErrors] = useState<{ name?: string; balance?: string }>({});
  const [submitError, setSubmitError] = useState<string | null>(null);

  const { data, isLoading, isError } = useQuery<Account[]>({
    queryKey: ['accounts'],
    queryFn: async () => {
      const { data } = await api.get<Account[]>('/accounts');
      return data;
    }
  });

  const createMutation = useMutation({
    mutationFn: async () => {
      await api.post('/accounts', form);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['accounts'] });
      setOpen(false);
      setForm({ name: '', type: 'checking', currency: defaultCurrency, description: '', balance: 0 });
      setFormErrors({});
      setSubmitError(null);
    },
    onError: (error) => {
      const message = error instanceof Error ? error.message : 'Failed to create account.';
      setSubmitError(message);
    }
  });

  const validateForm = () => {
    const errors: { name?: string; balance?: string } = {};
    if (!form.name.trim()) {
      errors.name = 'Name is required.';
    }
    if (Number.isNaN(form.balance)) {
      errors.balance = 'Balance must be a number.';
    }
    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = () => {
    if (!validateForm()) {
      setSubmitError('Please fix the highlighted fields before saving.');
      return;
    }
    setSubmitError(null);
    createMutation.mutate();
  };

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await api.delete(`/accounts/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['accounts'] });
    }
  });

  return (
    <Box>
      <Stack direction="row" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Accounts</Typography>
        <Button variant="contained" onClick={() => setOpen(true)}>
          New Account
        </Button>
      </Stack>

      {isError && <Alert severity="error">Failed to load accounts list.</Alert>}

      <Paper>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Name</TableCell>
              <TableCell>Type</TableCell>
              <TableCell>Currency</TableCell>
              <TableCell align="right">Balance</TableCell>
              <TableCell>Description</TableCell>
              <TableCell align="right">Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {isLoading &&
              Array.from({ length: 5 }).map((_, index) => (
                <TableRow key={`account-skeleton-${index}`}>
                  <TableCell><Skeleton variant="text" width={160} /></TableCell>
                  <TableCell><Skeleton variant="text" width={120} /></TableCell>
                  <TableCell><Skeleton variant="text" width={80} /></TableCell>
                  <TableCell align="right"><Skeleton variant="text" width={100} /></TableCell>
                  <TableCell><Skeleton variant="text" /></TableCell>
                  <TableCell align="right"><Skeleton variant="circular" width={24} height={24} /></TableCell>
                </TableRow>
              ))}
            {!isLoading && data?.map((account) => (
              <TableRow key={account.id} hover>
                <TableCell>{account.name}</TableCell>
                <TableCell>{account.type}</TableCell>
                <TableCell>{account.currency}</TableCell>
                <TableCell align="right">${account.balance.toFixed(2)}</TableCell>
                <TableCell>{account.description}</TableCell>
                <TableCell align="right">
                  <IconButton
                    aria-label="delete"
                    onClick={() => deleteMutation.mutate(account.id)}
                    size="small"
                  >
                    <DeleteIcon fontSize="small" />
                  </IconButton>
                </TableCell>
              </TableRow>
            ))}
            {!isLoading && !isError && (data?.length ?? 0) === 0 && (
              <TableRow>
                <TableCell colSpan={6} align="center">
                  No accounts registered yet.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </Paper>

      <Dialog
        open={open}
        onClose={() => {
          setOpen(false);
          setSubmitError(null);
          setFormErrors({});
        }}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>Create Account</DialogTitle>
        <DialogContent>
          {submitError && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {submitError}
            </Alert>
          )}
          <Grid container spacing={2} sx={{ mt: 0 }}>
            <Grid item xs={12}>
              <TextField
                label="Name"
                value={form.name}
                onChange={(event) => setForm((prev) => ({ ...prev, name: event.target.value }))}
                onBlur={validateForm}
                onFocus={() => setFormErrors((prev) => ({ ...prev, name: undefined }))}
                fullWidth
                required
                error={Boolean(formErrors.name)}
                helperText={formErrors.name}
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                label="Type"
                select
                value={form.type}
                onChange={(event) => setForm((prev) => ({ ...prev, type: event.target.value }))}
                fullWidth
              >
                {accountTypes.map((option) => (
                  <MenuItem key={option.value} value={option.value}>
                    {option.label}
                  </MenuItem>
                ))}
              </TextField>
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                label="Currency"
                select
                value={form.currency}
                onChange={(event) =>
                  setForm((prev) => ({ ...prev, currency: event.target.value as CurrencyCode }))
                }
                fullWidth
              >
                {currencyOptions.map((option) => (
                  <MenuItem key={option.value} value={option.value}>
                    {option.label}
                  </MenuItem>
                ))}
              </TextField>
            </Grid>
            <Grid item xs={12}>
              <TextField
                label="Initial Balance"
                type="number"
                value={form.balance}
                onChange={(event) =>
                  setForm((prev) => ({ ...prev, balance: Number(event.target.value) }))
                }
                onBlur={validateForm}
                onFocus={() => setFormErrors((prev) => ({ ...prev, balance: undefined }))}
                fullWidth
                error={Boolean(formErrors.balance)}
                helperText={formErrors.balance}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                label="Description"
                value={form.description}
                onChange={(event) =>
                  setForm((prev) => ({ ...prev, description: event.target.value }))
                }
                fullWidth
                multiline
                minRows={2}
              />
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button
            onClick={() => {
              setOpen(false);
              setSubmitError(null);
              setFormErrors({});
            }}
          >
            Cancel
          </Button>
          <Button onClick={handleSubmit} disabled={createMutation.isPending}>
            {createMutation.isPending ? 'Saving...' : 'Save'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default AccountsPage;
