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
  LinearProgress,
  MenuItem,
  Paper,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  TextField,
  Typography
} from '@mui/material';
import dayjs from 'dayjs';

import { api } from '../services/api';
import { currencyOptions, defaultCurrency, CurrencyCode } from '../constants/currencyOptions';

interface Goal {
  id: string;
  name: string;
  targetAmount: number;
  currentAmount: number;
  currency: string;
  deadline: string;
  status: string;
  description: string;
}

const GoalsPage = () => {
  const queryClient = useQueryClient();
  const [open, setOpen] = useState(false);
  const [progressOpen, setProgressOpen] = useState(false);
  const [selectedGoal, setSelectedGoal] = useState<Goal | null>(null);
  const [form, setForm] = useState<{
    name: string;
    targetAmount: number;
    currency: CurrencyCode;
    deadline: string;
    description: string;
  }>({
    name: '',
    targetAmount: 0,
    currency: defaultCurrency,
    deadline: dayjs().add(6, 'month').format('YYYY-MM-DD'),
    description: ''
  });
  const [progressAmount, setProgressAmount] = useState(0);

  const goalsQuery = useQuery<Goal[]>({
    queryKey: ['goals'],
    queryFn: async () => {
      const { data } = await api.get<Goal[]>('/goals');
      return data;
    }
  });

  const createMutation = useMutation({
    mutationFn: async () => {
      await api.post('/goals', {
        name: form.name,
        targetAmount: Number(form.targetAmount),
        currency: form.currency,
        deadline: new Date(form.deadline).toISOString(),
        description: form.description
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['goals'] });
      setOpen(false);
    }
  });

  const progressMutation = useMutation({
    mutationFn: async () => {
      if (!selectedGoal) return;
      await api.post(`/goals/${selectedGoal.id}/progress`, { amount: Number(progressAmount) });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['goals'] });
      setProgressOpen(false);
      setProgressAmount(0);
      setSelectedGoal(null);
    }
  });

  const openProgressDialog = (goal: Goal) => {
    setSelectedGoal(goal);
    setProgressAmount(0);
    setProgressOpen(true);
  };

  return (
    <Box>
      <Stack direction="row" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Goals</Typography>
        <Button variant="contained" onClick={() => setOpen(true)}>
          New Goal
        </Button>
      </Stack>

      {goalsQuery.isError && <Alert severity="error">Unable to load goals.</Alert>}

      <Paper>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Name</TableCell>
              <TableCell align="right">Target</TableCell>
              <TableCell align="right">Current</TableCell>
              <TableCell>Deadline</TableCell>
              <TableCell>Status</TableCell>
              <TableCell>Progress</TableCell>
              <TableCell align="right">Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {goalsQuery.isLoading && (
              <TableRow>
                <TableCell colSpan={7}>Loading...</TableCell>
              </TableRow>
            )}
            {goalsQuery.data?.map((goal) => {
              const progress = (goal.currentAmount / goal.targetAmount) * 100;
              return (
                <TableRow key={goal.id} hover>
                  <TableCell>{goal.name}</TableCell>
                  <TableCell align="right">
                    {goal.currency} {goal.targetAmount.toFixed(2)}
                  </TableCell>
                  <TableCell align="right">
                    {goal.currency} {goal.currentAmount.toFixed(2)}
                  </TableCell>
                  <TableCell>{dayjs(goal.deadline).format('YYYY-MM-DD')}</TableCell>
                  <TableCell>{goal.status}</TableCell>
                  <TableCell>
                    <Stack spacing={1}>
                      <LinearProgress variant="determinate" value={Math.min(100, progress)} />
                      <Typography variant="caption">{progress.toFixed(1)}%</Typography>
                    </Stack>
                  </TableCell>
                  <TableCell align="right">
                    <Button onClick={() => openProgressDialog(goal)} size="small">
                      Update
                    </Button>
                  </TableCell>
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
      </Paper>

      <Dialog open={open} onClose={() => setOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Create Goal</DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 0 }}>
            <Grid item xs={12}>
              <TextField
                label="Name"
                value={form.name}
                onChange={(event) => setForm((prev) => ({ ...prev, name: event.target.value }))}
                fullWidth
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                label="Target Amount"
                type="number"
                value={form.targetAmount}
                onChange={(event) =>
                  setForm((prev) => ({ ...prev, targetAmount: Number(event.target.value) }))
                }
                fullWidth
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
                label="Deadline"
                type="date"
                value={form.deadline}
                onChange={(event) => setForm((prev) => ({ ...prev, deadline: event.target.value }))}
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
                minRows={3}
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

      <Dialog open={progressOpen} onClose={() => setProgressOpen(false)} maxWidth="xs" fullWidth>
        <DialogTitle>Update Progress</DialogTitle>
        <DialogContent>
          <Stack spacing={2}>
            <Typography>Goal: {selectedGoal?.name}</Typography>
            <TextField
              label="Amount"
              type="number"
              value={progressAmount}
              onChange={(event) => setProgressAmount(Number(event.target.value))}
              fullWidth
            />
          </Stack>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setProgressOpen(false)}>Cancel</Button>
          <Button onClick={() => progressMutation.mutate()} disabled={progressMutation.isPending}>
            {progressMutation.isPending ? 'Updating...' : 'Update'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default GoalsPage;
