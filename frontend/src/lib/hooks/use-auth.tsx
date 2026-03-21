'use client';

import { createContext, useCallback, useContext, useState, useEffect, ReactNode } from 'react';
import { Tenant, Store, LoginResponse } from '@/lib/types';
import { apiClient } from '@/lib/api/client';

interface AuthContextType {
  user: Tenant | null;
  token: string | null;
  currentStore: Store | null;
  stores: Store[];
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (response: LoginResponse) => void;
  logout: () => void;
  setCurrentStore: (store: Store) => void;
  refreshStores: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<Tenant | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [currentStore, setCurrentStore] = useState<Store | null>(null);
  const [stores, setStores] = useState<Store[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const isAuthenticated = !!user && !!token;

  const logout = useCallback(() => {
    setUser(null);
    setToken(null);
    setCurrentStore(null);
    setStores([]);

    localStorage.removeItem('auth_token');
    localStorage.removeItem('auth_user');

    setIsLoading(false);
  }, []);

  const loadStores = useCallback(async (authToken: string) => {
    try {
      const originalToken = apiClient.getToken();
      apiClient.setToken(authToken);

      const userStores = await apiClient.getStores();

      setStores(userStores);

      if (userStores.length > 0 && !currentStore) {
        setCurrentStore(userStores[0]);
      }

      if (originalToken) {
        apiClient.setToken(originalToken);
      }
    } catch (error) {
      console.error('Failed to load stores:', error);
    } finally {
      setIsLoading(false);
    }
  }, [currentStore]);

  useEffect(() => {
    const storedToken = localStorage.getItem('auth_token');
    const storedUser = localStorage.getItem('auth_user');

    if (storedToken && storedUser) {
      try {
        const parsedUser = JSON.parse(storedUser) as Tenant;
        setToken(storedToken);
        setUser(parsedUser);
        loadStores(storedToken);
      } catch (error) {
        console.error('Failed to parse stored auth data:', error);
        logout();
      }
    } else {
      setIsLoading(false);
    }
  }, [loadStores, logout]);

  const login = (response: LoginResponse) => {
    setUser(response.tenant);
    setToken(response.token);

    localStorage.setItem('auth_token', response.token);
    localStorage.setItem('auth_user', JSON.stringify(response.tenant));

    loadStores(response.token);
  };

  const refreshStores = async () => {
    if (user && token) {
      await loadStores(token);
    }
  };

  const value: AuthContextType = {
    user,
    token,
    currentStore,
    stores,
    isLoading,
    isAuthenticated,
    login,
    logout,
    setCurrentStore,
    refreshStores,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}