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
  const [cursorPosition, setCursorPosition] = useState(0);
  const inputRef = useRef<HTMLInputElement>(null);
  const config = currencyConfig[currency];

  // Função para formatar número para exibição (com separadores de milhares)
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

  // Função para formatar valor durante digitação (sem separadores de milhares)
  const formatForTyping = (str: string): string => {
    if (!str) return '';
    
    // Remove símbolos e espaços
    let cleanStr = str.replace(/[^\d.,]/g, '');
    
    // Para BRL (vírgula como separador decimal)
    if (config.separator === ',') {
      // Remove pontos (separadores de milhares) mas mantém vírgula
      cleanStr = cleanStr.replace(/\./g, '');
      
      // Limita a uma vírgula e máximo 2 dígitos após ela
      const parts = cleanStr.split(',');
      if (parts.length > 2) {
        cleanStr = parts[0] + ',' + parts.slice(1).join('');
      }
      if (parts.length === 2 && parts[1].length > 2) {
        cleanStr = parts[0] + ',' + parts[1].substring(0, 2);
      }
    } else {
      // Para outras moedas (ponto como separador decimal)
      cleanStr = cleanStr.replace(/,/g, '');
      
      // Limita a um ponto e máximo 2 dígitos após ele
      const parts = cleanStr.split('.');
      if (parts.length > 2) {
        cleanStr = parts[0] + '.' + parts.slice(1).join('');
      }
      if (parts.length === 2 && parts[1].length > 2) {
        cleanStr = parts[0] + '.' + parts[1].substring(0, 2);
      }
    }
    
    return cleanStr;
  };

  // Função para calcular nova posição do cursor após formatação
  const calculateCursorPosition = (oldValue: string, newValue: string, oldCursorPos: number): number => {
    const oldLength = oldValue.length;
    const newLength = newValue.length;
    const lengthDiff = newLength - oldLength;
    
    // Se o cursor estava no final, mantém no final
    if (oldCursorPos >= oldLength) {
      return newLength;
    }
    
    // Ajusta a posição baseado na diferença de tamanho
    return Math.max(0, Math.min(newLength, oldCursorPos + lengthDiff));
  };

  // Atualiza o valor de exibição quando o valor numérico muda (apenas quando não está focado)
  useEffect(() => {
    if (!isFocused) {
      setDisplayValue(formatForDisplay(value));
    }
  }, [value, currency, isFocused]);

  // Atualiza posição do cursor após mudança no displayValue
  useEffect(() => {
    if (inputRef.current && isFocused) {
      const input = inputRef.current.querySelector('input');
      if (input) {
        input.setSelectionRange(cursorPosition, cursorPosition);
      }
    }
  }, [displayValue, cursorPosition, isFocused]);

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const inputValue = event.target.value;
    const cursorPos = event.target.selectionStart || 0;
    
    // Formata o valor durante digitação (sem separadores de milhares)
    const formattedValue = formatForTyping(inputValue);
    
    // Calcula nova posição do cursor
    const newCursorPos = calculateCursorPosition(inputValue, formattedValue, cursorPos);
    
    setDisplayValue(formattedValue);
    setCursorPosition(newCursorPos);
    
    // Converte para número e chama onChange
    const numericValue = parseFromDisplay(formattedValue);
    onChange(numericValue);
  };

  const handleBlur = () => {
    setIsFocused(false);
    
    // Formata o valor com separadores de milhares quando perde o foco
    const numericValue = parseFromDisplay(displayValue);
    setDisplayValue(formatForDisplay(numericValue));
    onChange(numericValue);
  };

  const handleFocus = () => {
    setIsFocused(true);
    
    // Quando ganha foco, prepara o valor para edição (sem separadores de milhares)
    if (value === 0) {
      setDisplayValue('');
      setCursorPosition(0);
    } else {
      // Remove separadores de milhares para facilitar edição
      const formattedValue = formatForDisplay(value);
      const cleanValue = formattedValue.replace(new RegExp(`\\${config.thousandSeparator}`, 'g'), '');
      setDisplayValue(cleanValue);
      setCursorPosition(cleanValue.length);
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
