import React, { createContext, useCallback, useEffect, useMemo, useState } from 'react';

export interface AuthTokens {
  accessToken: string;
  refreshToken: string;
  idToken: string;
  expiresIn: number;
  tokenType: string;
}

interface AuthContextValue {
  tokens: AuthTokens | null;
  setTokens: (tokens: AuthTokens | null) => void;
  isAuthenticated: boolean;
  logout: () => void;
}

export const AuthContext = createContext<AuthContextValue | undefined>(undefined);

const STORAGE_KEY = 'financial-control-auth';

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [tokens, setTokensState] = useState<AuthTokens | null>(() => {
    if (typeof window === 'undefined') {
      return null;
    }
    const raw = window.localStorage.getItem(STORAGE_KEY);
    if (!raw) return null;
    try {
      return JSON.parse(raw) as AuthTokens;
    } catch (error) {
      console.error('Failed to parse auth tokens', error);
      return null;
    }
  });

  useEffect(() => {
    if (typeof window === 'undefined') {
      return;
    }
    if (tokens) {
      window.localStorage.setItem(STORAGE_KEY, JSON.stringify(tokens));
    } else {
      window.localStorage.removeItem(STORAGE_KEY);
    }
  }, [tokens]);

  const setTokens = useCallback((value: AuthTokens | null) => {
    setTokensState(value);
  }, []);

  const logout = useCallback(() => {
    setTokensState(null);
  }, []);

  const value = useMemo<AuthContextValue>(() => ({
    tokens,
    setTokens,
    isAuthenticated: Boolean(tokens?.accessToken),
    logout
  }), [tokens, setTokens, logout]);

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
