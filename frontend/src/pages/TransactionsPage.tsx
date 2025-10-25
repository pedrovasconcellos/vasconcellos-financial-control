import { useMemo, useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  Alert,
  Box,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Grid,
  MenuItem,
  TextField,
  Typography,
  Paper,
  Table,
  TableHead,
  TableBody,
  TableRow,
  TableCell,
  Stack,
  Chip
} from '@mui/material';
import dayjs from 'dayjs';

import { api } from '../services/api';

interface Transaction {
  id: string;
  accountId: string;
  categoryId: string;
  amount: number;
  currency: string;
  description: string;
  occurredAt: string;
  status: string;
  tags: string[];
  notes: string;
  receiptUrl?: string | null;
}

interface AccountOption {
  id: string;
  name: string;
}

interface CategoryOption {
  id: string;
  name: string;
}

const TransactionsPage = () => {
  const queryClient = useQueryClient();
  const [open, setOpen] = useState(false);
  const [form, setForm] = useState({
    accountId: '',
    categoryId: '',
    amount: 0,
    currency: 'USD',
    description: '',
    occurredAt: dayjs().toISOString()
  });

  const transactionsQuery = useQuery<Transaction[]>({
    queryKey: ['transactions'],
    queryFn: async () => {
      const { data } = await api.get<Transaction[]>('/transactions');
      return data;
    }
  });

  const accountsQuery = useQuery<AccountOption[]>({
    queryKey: ['accounts-options'],
    queryFn: async () => {
      const { data } = await api.get<AccountOption[]>('/accounts');
      return data;
    }
  });

  const categoriesQuery = useQuery<CategoryOption[]>({
    queryKey: ['categories-options'],
    queryFn: async () => {
      const { data } = await api.get<CategoryOption[]>('/categories');
      return data;
    }
  });

  const createMutation = useMutation({
    mutationFn: async () => {
      await api.post('/transactions', {
        accountId: form.accountId,
        categoryId: form.categoryId,
        amount: Number(form.amount),
        currency: form.currency,
        description: form.description,
        occurredAt: form.occurredAt
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['transactions'] });
      setOpen(false);
      setForm({
        accountId: '',
        categoryId: '',
        amount: 0,
        currency: 'USD',
        description: '',
        occurredAt: dayjs().toISOString()
      });
    }
  });

  const accounts = accountsQuery.data ?? [];
  const categories = categoriesQuery.data ?? [];

  const accountMap = useMemo(
    () => Object.fromEntries(accounts.map((item) => [item.id, item.name])),
    [accounts]
  );
  const categoryMap = useMemo(
    () => Object.fromEntries(categories.map((item) => [item.id, item.name])),
    [categories]
  );

  return (
    <Box>
      <Stack direction="row" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Transactions</Typography>
        <Button variant="contained" onClick={() => setOpen(true)}>
          New Transaction
        </Button>
      </Stack>

      {(transactionsQuery.isError || accountsQuery.isError || categoriesQuery.isError) && (
        <Alert severity="error">Unable to load transaction data.</Alert>
      )}

      <Paper>
        <Table size="small">
          <TableHead>
            <TableRow>
              <TableCell>Date</TableCell>
              <TableCell>Description</TableCell>
              <TableCell>Account</TableCell>
              <TableCell>Category</TableCell>
              <TableCell align="right">Amount</TableCell>
              <TableCell>Status</TableCell>
              <TableCell>Receipt</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {transactionsQuery.isLoading && (
              <TableRow>
                <TableCell colSpan={7}>Loading...</TableCell>
              </TableRow>
            )}
            {transactionsQuery.data?.map((transaction) => (
              <TableRow key={transaction.id} hover>
                <TableCell>{dayjs(transaction.occurredAt).format('YYYY-MM-DD')}</TableCell>
                <TableCell>{transaction.description}</TableCell>
                <TableCell>{accountMap[transaction.accountId] ?? transaction.accountId}</TableCell>
                <TableCell>{categoryMap[transaction.categoryId] ?? transaction.categoryId}</TableCell>
                <TableCell align="right">
                  {transaction.currency} {transaction.amount.toFixed(2)}
                </TableCell>
                <TableCell>
                  <Chip label={transaction.status} size="small" color="primary" variant="outlined" />
                </TableCell>
                <TableCell>
                  {transaction.receiptUrl ? (
                    <Button href={transaction.receiptUrl} target="_blank" rel="noopener">
                      Receipt
                    </Button>
                  ) : (
                    '—'
                  )}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </Paper>

      <Dialog open={open} onClose={() => setOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Record Transaction</DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 0 }}>
            <Grid item xs={12}>
              <TextField
                label="Account"
                select
                value={form.accountId}
                onChange={(event) => setForm((prev) => ({ ...prev, accountId: event.target.value }))}
                fullWidth
                required
              >
                {accounts.map((account) => (
                  <MenuItem key={account.id} value={account.id}>
                    {account.name}
                  </MenuItem>
                ))}
              </TextField>
            </Grid>
            <Grid item xs={12}>
              <TextField
                label="Category"
                select
                value={form.categoryId}
                onChange={(event) => setForm((prev) => ({ ...prev, categoryId: event.target.value }))}
                fullWidth
                required
              >
                {categories.map((category) => (
                  <MenuItem key={category.id} value={category.id}>
                    {category.name}
                  </MenuItem>
                ))}
              </TextField>
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                label="Amount"
                type="number"
                value={form.amount}
                onChange={(event) => setForm((prev) => ({ ...prev, amount: Number(event.target.value) }))}
                fullWidth
                required
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                label="Currency"
                value={form.currency}
                onChange={(event) => setForm((prev) => ({ ...prev, currency: event.target.value }))}
                fullWidth
                inputProps={{ maxLength: 3 }}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                label="Date"
                type="date"
                value={dayjs(form.occurredAt).format('YYYY-MM-DD')}
                onChange={(event) =>
                  setForm((prev) => ({
                    ...prev,
                    occurredAt: dayjs(event.target.value).toISOString()
                  }))
                }
                fullWidth
                InputLabelProps={{ shrink: true }}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                label="Description"
                value={form.description}
                onChange={(event) => setForm((prev) => ({ ...prev, description: event.target.value }))}
                fullWidth
                multiline
              />
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpen(false)}>Cancel</Button>
          <Button onClick={() => createMutation.mutate()} disabled={createMutation.isPending}>
            {createMutation.isPending ? 'Saving...' : 'Save'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default TransactionsPage;
