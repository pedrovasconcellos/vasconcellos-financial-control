export const currencyOptions = [
  { value: 'USD', label: 'United States Dollar' },
  { value: 'EUR', label: 'Euro' },
  { value: 'CHF', label: 'Swiss Franc' },
  { value: 'GBP', label: 'Pound Sterling' },
  { value: 'BRL', label: 'Brazilian Real' }
] as const;

export type CurrencyCode = (typeof currencyOptions)[number]['value'];

export const defaultCurrency: CurrencyCode = currencyOptions[0].value;
