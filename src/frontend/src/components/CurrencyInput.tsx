import React, { useState, useEffect, useRef } from 'react';
import { TextField, TextFieldProps } from '@mui/material';
import { CurrencyCode } from '../constants/currencyOptions';

interface CurrencyInputProps extends Omit<TextFieldProps, 'onChange' | 'value'> {
  value: number;
  onChange: (value: number) => void;
  currency: CurrencyCode;
}

// Configurações de formatação para cada moeda
const currencyConfig = {
  USD: { symbol: '$', decimals: 2, separator: '.', thousandSeparator: ',' },
  EUR: { symbol: '€', decimals: 2, separator: '.', thousandSeparator: ',' },
  CHF: { symbol: 'CHF', decimals: 2, separator: '.', thousandSeparator: "'" },
  GBP: { symbol: '£', decimals: 2, separator: '.', thousandSeparator: ',' },
  BRL: { symbol: 'R$', decimals: 2, separator: ',', thousandSeparator: '.' }
};

export const CurrencyInput: React.FC<CurrencyInputProps> = ({
  value,
  onChange,
  currency,
  ...textFieldProps
}) => {
  const [displayValue, setDisplayValue] = useState('');
  const [isFocused, setIsFocused] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);
  const config = currencyConfig[currency];

  // Função para formatar número para exibição
  const formatForDisplay = (num: number): string => {
    if (isNaN(num) || num === 0) return '';
    
    // Arredonda para o número correto de decimais
    const rounded = Math.round(num * Math.pow(10, config.decimals)) / Math.pow(10, config.decimals);
    
    // Formata com separadores de milhares e sempre mostra decimais
    const parts = rounded.toString().split('.');
    const integerPart = parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, config.thousandSeparator);
    
    // Sempre inclui os decimais, mesmo que sejam .00
    const decimalPart = parts[1] ? parts[1].padEnd(config.decimals, '0') : '00';
    
    return `${integerPart}${config.separator}${decimalPart}`;
  };

  // Função para converter string formatada de volta para número
  const parseFromDisplay = (str: string): number => {
    if (!str) return 0;
    
    // Remove símbolos e espaços
    let cleanStr = str.replace(/[^\d.,]/g, '');
    
    // Converte separadores para formato padrão
    if (config.separator === ',') {
      // Moedas que usam vírgula como separador decimal (BRL)
      cleanStr = cleanStr.replace(/\./g, '').replace(',', '.');
    } else {
      // Moedas que usam ponto como separador decimal (USD, EUR, CHF, GBP)
      cleanStr = cleanStr.replace(/,/g, '');
    }
    
    const parsed = parseFloat(cleanStr);
    return isNaN(parsed) ? 0 : Math.round(parsed * Math.pow(10, config.decimals)) / Math.pow(10, config.decimals);
  };

  // Função para limpar e preparar valor para digitação
  const prepareForEditing = (str: string): string => {
    if (!str) return '';
    
    // Remove separadores de milhares mas mantém o separador decimal
    let cleanStr = str.replace(/[^\d.,]/g, '');
    
    if (config.separator === ',') {
      // Para BRL: remove pontos (milhares) mas mantém vírgula (decimal)
      cleanStr = cleanStr.replace(/\./g, '');
    } else {
      // Para outras moedas: remove vírgulas (milhares) mas mantém ponto (decimal)
      cleanStr = cleanStr.replace(/,/g, '');
    }
    
    return cleanStr;
  };

  // Atualiza o valor de exibição quando o valor numérico muda (apenas quando não está focado)
  useEffect(() => {
    if (!isFocused) {
      setDisplayValue(formatForDisplay(value));
    }
  }, [value, currency, isFocused]);

  // Função para validar e limitar entrada durante digitação
  const validateAndLimitInput = (inputValue: string): string => {
    if (!inputValue) return '';
    
    // Remove símbolos e espaços
    let cleanStr = inputValue.replace(/[^\d.,]/g, '');
    
    // Para BRL (vírgula como separador decimal)
    if (config.separator === ',') {
      // Remove pontos (separadores de milhares) mas mantém vírgula
      cleanStr = cleanStr.replace(/\./g, '');
      
      // Limita a uma vírgula e máximo 2 dígitos após ela
      const parts = cleanStr.split(',');
      if (parts.length > 2) {
        // Se há mais de uma vírgula, mantém apenas a primeira
        cleanStr = parts[0] + ',' + parts.slice(1).join('');
      }
      if (parts.length === 2 && parts[1].length > 2) {
        // Limita a 2 casas decimais
        cleanStr = parts[0] + ',' + parts[1].substring(0, 2);
      }
    } else {
      // Para outras moedas (ponto como separador decimal)
      cleanStr = cleanStr.replace(/,/g, '');
      
      // Limita a um ponto e máximo 2 dígitos após ele
      const parts = cleanStr.split('.');
      if (parts.length > 2) {
        // Se há mais de um ponto, mantém apenas o primeiro
        cleanStr = parts[0] + '.' + parts.slice(1).join('');
      }
      if (parts.length === 2 && parts[1].length > 2) {
        // Limita a 2 casas decimais
        cleanStr = parts[0] + '.' + parts[1].substring(0, 2);
      }
    }
    
    return cleanStr;
  };

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const inputValue = event.target.value;
    
    // Valida e limita a entrada durante digitação
    const validatedValue = validateAndLimitInput(inputValue);
    setDisplayValue(validatedValue);
    
    // Converte para número e chama onChange
    const numericValue = parseFromDisplay(validatedValue);
    onChange(numericValue);
  };

  const handleBlur = () => {
    setIsFocused(false);
    
    // Formata o valor quando o campo perde o foco
    const numericValue = parseFromDisplay(displayValue);
    setDisplayValue(formatForDisplay(numericValue));
    onChange(numericValue);
  };

  const handleFocus = () => {
    setIsFocused(true);
    
    // Quando ganha foco, prepara o valor para edição
    if (value === 0) {
      setDisplayValue('');
    } else {
      // Remove formatação para facilitar edição
      setDisplayValue(prepareForEditing(formatForDisplay(value)));
    }
  };

  return (
    <TextField
      {...textFieldProps}
      ref={inputRef}
      value={displayValue}
      onChange={handleChange}
      onBlur={handleBlur}
      onFocus={handleFocus}
      InputProps={{
        ...textFieldProps.InputProps,
        startAdornment: (
          <span style={{ marginRight: 8, color: '#666', fontSize: '1rem' }}>
            {config.symbol}
          </span>
        ),
      }}
      placeholder={`0${config.separator}00`}
      inputProps={{
        ...textFieldProps.inputProps,
        inputMode: 'decimal',
        pattern: '[0-9]*[.,]?[0-9]*'
      }}
    />
  );
};

export default CurrencyInput;
