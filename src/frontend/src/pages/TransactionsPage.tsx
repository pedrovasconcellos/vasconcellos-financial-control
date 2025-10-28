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
  Chip,
  Skeleton
} from '@mui/material';
import dayjs from 'dayjs';

import { api } from '../services/api';
import { currencyOptions, defaultCurrency, CurrencyCode } from '../constants/currencyOptions';
import CurrencyInput from '../components/CurrencyInput';

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

interface TransactionFormState {
  accountId: string;
  categoryId: string;
  amount: number;
  currency: CurrencyCode;
  description: string;
  occurredAt: string;
}

const TransactionsPage = () => {
  const queryClient = useQueryClient();
  const [open, setOpen] = useState(false);
  const [form, setForm] = useState<TransactionFormState>({
    accountId: '',
    categoryId: '',
    amount: 0,
    currency: defaultCurrency,
    description: '',
    occurredAt: dayjs().toISOString()
  });
  const [formErrors, setFormErrors] = useState<{ accountId?: string; categoryId?: string; amount?: string }>({});
  const [submitError, setSubmitError] = useState<string | null>(null);

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
        currency: defaultCurrency,
        description: '',
        occurredAt: dayjs().toISOString()
      });
      setFormErrors({});
      setSubmitError(null);
    },
    onError: (error) => {
      const message = error instanceof Error ? error.message : 'Failed to save transaction.';
      setSubmitError(message);
    }
  });

  const validateForm = () => {
    const errors: { accountId?: string; categoryId?: string; amount?: string } = {};
    if (!form.accountId) {
      errors.accountId = 'Account is required.';
    }
    if (!form.categoryId) {
      errors.categoryId = 'Category is required.';
    }
    if (!form.amount || Number.isNaN(form.amount) || form.amount <= 0) {
      errors.amount = 'Amount must be greater than zero.';
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

  const showTransactionsError = transactionsQuery.isError;
  const showAccountsError = accountsQuery.isError;
  const showCategoriesError = categoriesQuery.isError;

  return (
    <Box>
      <Stack direction="row" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Transactions</Typography>
        <Button variant="contained" onClick={() => setOpen(true)}>
          New Transaction
        </Button>
      </Stack>

      {(showTransactionsError || showAccountsError || showCategoriesError) && (
        <Stack spacing={2} mb={2}>
          {showTransactionsError && <Alert severity="error">Failed to load transactions.</Alert>}
          {showAccountsError && <Alert severity="error">Failed to load accounts list.</Alert>}
          {showCategoriesError && <Alert severity="error">Failed to load categories list.</Alert>}
        </Stack>
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
            {transactionsQuery.isLoading &&
              Array.from({ length: 5 }).map((_, index) => (
                <TableRow key={`transaction-skeleton-${index}`}>
                  <TableCell><Skeleton variant="text" width={120} /></TableCell>
                  <TableCell><Skeleton variant="text" /></TableCell>
                  <TableCell><Skeleton variant="text" width={140} /></TableCell>
                  <TableCell><Skeleton variant="text" width={140} /></TableCell>
                  <TableCell align="right"><Skeleton variant="text" width={100} /></TableCell>
                  <TableCell><Skeleton variant="rectangular" width={80} height={24} /></TableCell>
                  <TableCell><Skeleton variant="text" width={80} /></TableCell>
                </TableRow>
              ))}
            {!transactionsQuery.isLoading &&
              transactionsQuery.data?.map((transaction) => (
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
            {!transactionsQuery.isLoading && !showTransactionsError && (transactionsQuery.data?.length ?? 0) === 0 && (
              <TableRow>
                <TableCell colSpan={7} align="center">
                  No transactions recorded yet.
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
        <DialogTitle>Record Transaction</DialogTitle>
        <DialogContent>
          {submitError && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {submitError}
            </Alert>
          )}
          <Grid container spacing={2} sx={{ mt: 0 }}>
            <Grid item xs={12}>
                <TextField
                  label="Account"
                  select
                  value={form.accountId}
                  onChange={(event) => setForm((prev) => ({ ...prev, accountId: event.target.value }))}
                  onBlur={validateForm}
                  onFocus={() => setFormErrors((prev) => ({ ...prev, accountId: undefined }))}
                  fullWidth
                  required
                  error={Boolean(formErrors.accountId)}
                  helperText={formErrors.accountId}
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
                  onBlur={validateForm}
                  onFocus={() => setFormErrors((prev) => ({ ...prev, categoryId: undefined }))}
                  fullWidth
                  required
                  error={Boolean(formErrors.categoryId)}
                  helperText={formErrors.categoryId}
                >
                  {categories.map((category) => (
                    <MenuItem key={category.id} value={category.id}>
                      {category.name}
                    </MenuItem>
                ))}
              </TextField>
            </Grid>
            <Grid item xs={12} md={6}>
                <CurrencyInput
                  label="Amount"
                  value={form.amount}
                  onChange={(value) => setForm((prev) => ({ ...prev, amount: value }))}
                  currency={form.currency}
                  onBlur={validateForm}
                  onFocus={() => setFormErrors((prev) => ({ ...prev, amount: undefined }))}
                  fullWidth
                  required
                  error={Boolean(formErrors.amount)}
                  helperText={formErrors.amount ?? ''}
                />
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
          <Button
            onClick={() => {
              setOpen(false);
              setSubmitError(null);
              setFormErrors({});
            }}
          >
            Cancel
          </Button>
          <Button
            onClick={handleSubmit}
            disabled={createMutation.isPending}
          >
            {createMutation.isPending ? 'Saving...' : 'Save'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default TransactionsPage;
