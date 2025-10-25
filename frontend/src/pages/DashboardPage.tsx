import { useMemo } from 'react';
import { useQuery } from '@tanstack/react-query';
import {
  Grid,
  Paper,
  Typography,
  CircularProgress,
  Alert,
  Stack,
  LinearProgress
} from '@mui/material';

import { api } from '../services/api';

interface SummaryReportResponse {
  totalIncome: number;
  totalExpense: number;
  netBalance: number;
  spendingByCategory: Record<string, number>;
  budgetUsage: Record<string, number>;
  goalProgress: Record<string, number>;
}

const DashboardPage = () => {
  const { data, isLoading, isError } = useQuery<SummaryReportResponse>({
    queryKey: ['summary-report'],
    queryFn: async () => {
      const { data } = await api.get<SummaryReportResponse>('/reports/summary');
      return data;
    }
  });

  const spendingEntries = useMemo(
    () => Object.entries(data?.spendingByCategory ?? {}),
    [data?.spendingByCategory]
  );

  const budgetEntries = useMemo(
    () => Object.entries(data?.budgetUsage ?? {}),
    [data?.budgetUsage]
  );

  const goalEntries = useMemo(
    () => Object.entries(data?.goalProgress ?? {}),
    [data?.goalProgress]
  );

  if (isLoading) {
    return (
      <Stack alignItems="center" spacing={2}>
        <CircularProgress />
        <Typography>Loading financial overview...</Typography>
      </Stack>
    );
  }

  if (isError || !data) {
    return <Alert severity="error">Unable to load dashboard data.</Alert>;
  }

  return (
    <Grid container spacing={3}>
      <Grid item xs={12} md={4}>
        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle2" color="text.secondary">
            Total Income
          </Typography>
          <Typography variant="h5">${data.totalIncome.toFixed(2)}</Typography>
        </Paper>
      </Grid>
      <Grid item xs={12} md={4}>
        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle2" color="text.secondary">
            Total Expense
          </Typography>
          <Typography variant="h5">${data.totalExpense.toFixed(2)}</Typography>
        </Paper>
      </Grid>
      <Grid item xs={12} md={4}>
        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle2" color="text.secondary">
            Net Balance
          </Typography>
          <Typography variant="h5">${data.netBalance.toFixed(2)}</Typography>
        </Paper>
      </Grid>

      <Grid item xs={12} md={6}>
        <Paper sx={{ p: 3 }}>
          <Typography variant="h6" gutterBottom>
            Spending by Category
          </Typography>
          <Stack spacing={2}>
            {spendingEntries.length === 0 && (
              <Typography variant="body2" color="text.secondary">
                No spending data available.
              </Typography>
            )}
            {spendingEntries.map(([category, value]) => (
              <Stack key={category} spacing={1}>
                <Typography variant="body2">{category}</Typography>
                <LinearProgress variant="determinate" value={Math.min(100, value)} />
                <Typography variant="caption">${value.toFixed(2)}</Typography>
              </Stack>
            ))}
          </Stack>
        </Paper>
      </Grid>

      <Grid item xs={12} md={6}>
        <Paper sx={{ p: 3 }}>
          <Typography variant="h6" gutterBottom>
            Budget Usage
          </Typography>
          <Stack spacing={2}>
            {budgetEntries.length === 0 && (
              <Typography variant="body2" color="text.secondary">
                No budgets registered.
              </Typography>
            )}
            {budgetEntries.map(([category, value]) => (
              <Stack key={category} spacing={1}>
                <Typography variant="body2">{category}</Typography>
                <LinearProgress variant="determinate" value={Math.min(100, value)} />
                <Typography variant="caption">{value.toFixed(2)}%</Typography>
              </Stack>
            ))}
          </Stack>
        </Paper>
      </Grid>

      <Grid item xs={12}>
        <Paper sx={{ p: 3 }}>
          <Typography variant="h6" gutterBottom>
            Goal Progress
          </Typography>
          <Stack spacing={2}>
            {goalEntries.length === 0 && (
              <Typography variant="body2" color="text.secondary">
                No goals created yet.
              </Typography>
            )}
            {goalEntries.map(([goal, value]) => (
              <Stack key={goal} spacing={1}>
                <Typography variant="body2">{goal}</Typography>
                <LinearProgress variant="determinate" value={Math.min(100, value)} />
                <Typography variant="caption">{value.toFixed(2)}%</Typography>
              </Stack>
            ))}
          </Stack>
        </Paper>
      </Grid>
    </Grid>
  );
};

export default DashboardPage;
