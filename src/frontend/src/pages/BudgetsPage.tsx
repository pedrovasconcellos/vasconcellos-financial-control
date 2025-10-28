import { useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
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
  Paper,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  TextField,
  Typography,
  LinearProgress,
  Stack
} from '@mui/material';
import dayjs from 'dayjs';

import { api } from '../services/api';
import { currencyOptions, defaultCurrency, CurrencyCode } from '../constants/currencyOptions';
import CurrencyInput from '../components/CurrencyInput';

interface Budget {
  id: string;
  categoryId: string;
  amount: number;
  currency: string;
  period: string;
  periodStart: string;
  periodEnd: string;
  spent: number;
  alertPercent: number;
}

interface CategoryOption {
  id: string;
  name: string;
}

interface BudgetFormState {
  categoryId: string;
  amount: number;
  currency: CurrencyCode;
  period: string;
  periodStart: string;
  periodEnd: string;
  alertPercent: number;
}

const periods = [
  { label: 'Monthly', value: 'monthly' },
  { label: 'Quarterly', value: 'quarterly' },
  { label: 'Yearly', value: 'yearly' }
];

const BudgetsPage = () => {
  const queryClient = useQueryClient();
  const [open, setOpen] = useState(false);
  const [form, setForm] = useState<BudgetFormState>({
    categoryId: '',
    amount: 0,
    currency: defaultCurrency,
    period: 'monthly',
    periodStart: dayjs().startOf('month').format('YYYY-MM-DD'),
    periodEnd: dayjs().endOf('month').format('YYYY-MM-DD'),
    alertPercent: 80
  });

  const budgetsQuery = useQuery<Budget[]>({
    queryKey: ['budgets'],
    queryFn: async () => {
      const { data } = await api.get<Budget[]>('/budgets');
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
      await api.post('/budgets', {
        categoryId: form.categoryId,
        amount: Number(form.amount),
        currency: form.currency,
        period: form.period,
        periodStart: new Date(form.periodStart).toISOString(),
        periodEnd: new Date(form.periodEnd).toISOString(),
        alertPercent: Number(form.alertPercent)
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['budgets'] });
      setOpen(false);
    }
  });

  return (
    <Box>
      <Stack direction="row" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Budgets</Typography>
        <Button variant="contained" onClick={() => setOpen(true)}>
          New Budget
        </Button>
      </Stack>

      {(budgetsQuery.isError || categoriesQuery.isError) && (
        <Alert severity="error">Unable to load budgets.</Alert>
      )}

      <Paper>
        <Table size="small">
          <TableHead>
            <TableRow>
              <TableCell>Category</TableCell>
              <TableCell align="right">Amount</TableCell>
              <TableCell align="right">Spent</TableCell>
              <TableCell>Period</TableCell>
              <TableCell>Usage</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {budgetsQuery.isLoading && (
              <TableRow>
                <TableCell colSpan={5}>Loading...</TableCell>
              </TableRow>
            )}
            {budgetsQuery.data?.map((budget) => {
              const usage = (budget.spent / budget.amount) * 100;
              const categoryName = categoriesQuery.data?.find((c) => c.id === budget.categoryId)?.name ??
                budget.categoryId;
              return (
                <TableRow key={budget.id} hover>
                  <TableCell>{categoryName}</TableCell>
                  <TableCell align="right">
                    {budget.currency} {budget.amount.toFixed(2)}
                  </TableCell>
                  <TableCell align="right">
                    {budget.currency} {budget.spent.toFixed(2)}
                  </TableCell>
                  <TableCell>
                    {dayjs(budget.periodStart).format('YYYY-MM-DD')} -{' '}
                    {dayjs(budget.periodEnd).format('YYYY-MM-DD')} ({budget.period})
                  </TableCell>
                  <TableCell>
                    <Stack spacing={1}>
                      <LinearProgress variant="determinate" value={Math.min(100, usage)} />
                      <Typography variant="caption">{usage.toFixed(1)}%</Typography>
                    </Stack>
                  </TableCell>
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
      </Paper>

      <Dialog open={open} onClose={() => setOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Create Budget</DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 0 }}>
            <Grid item xs={12}>
              <TextField
                label="Category"
                select
                value={form.categoryId}
                onChange={(event) => setForm((prev) => ({ ...prev, categoryId: event.target.value }))}
                fullWidth
                required
              >
                {categoriesQuery.data?.map((category) => (
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
                fullWidth
                required
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
            <Grid item xs={12} md={6}>
              <TextField
                label="Period"
                select
                value={form.period}
                onChange={(event) => setForm((prev) => ({ ...prev, period: event.target.value }))}
                fullWidth
              >
                {periods.map((period) => (
                  <MenuItem key={period.value} value={period.value}>
                    {period.label}
                  </MenuItem>
                ))}
              </TextField>
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                label="Alert %"
                type="number"
                value={form.alertPercent}
                onChange={(event) => setForm((prev) => ({ ...prev, alertPercent: Number(event.target.value) }))}
                fullWidth
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                label="Start"
                type="date"
                value={form.periodStart}
                onChange={(event) => setForm((prev) => ({ ...prev, periodStart: event.target.value }))}
                fullWidth
                InputLabelProps={{ shrink: true }}
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                label="End"
                type="date"
                value={form.periodEnd}
                onChange={(event) => setForm((prev) => ({ ...prev, periodEnd: event.target.value }))}
                fullWidth
                InputLabelProps={{ shrink: true }}
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

export default BudgetsPage;
