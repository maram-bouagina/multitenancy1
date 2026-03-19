'use client';

import { createContext, useContext, useState, useEffect, ReactNode } from 'react';
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

  useEffect(() => {
    // Check for existing auth data on mount
    const storedToken = localStorage.getItem('auth_token');
    const storedUser = localStorage.getItem('auth_user');

    if (storedToken && storedUser) {
      try {
        const parsedUser = JSON.parse(storedUser);
        setToken(storedToken);
        setUser(parsedUser);
        // Load stores for this tenant
        loadStores(storedToken, parsedUser.id);
      } catch (error) {
        console.error('Failed to parse stored auth data:', error);
        logout();
      }
    } else {
      setIsLoading(false);
    }
  }, []);

  const loadStores = async (authToken: string, tenantId: string) => {
    try {
      // Set token temporarily for this request
      const originalToken = apiClient.getToken();
      apiClient.setToken(authToken);

      const userStores = await apiClient.getStores();

      setStores(userStores);

      // Set first store as current if no current store is set
      if (userStores.length > 0 && !currentStore) {
        setCurrentStore(userStores[0]);
      }

      // Restore original token
      if (originalToken) {
        apiClient.setToken(originalToken);
      }
    } catch (error) {
      console.error('Failed to load stores:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const login = (response: LoginResponse) => {
    setUser(response.tenant);
    setToken(response.token);

    // Store in localStorage
    localStorage.setItem('auth_token', response.token);
    localStorage.setItem('auth_user', JSON.stringify(response.tenant));

    // Load stores for this tenant
    loadStores(response.token, response.tenant.id);
  };

  const logout = () => {
    setUser(null);
    setToken(null);
    setCurrentStore(null);
    setStores([]);

    // Clear localStorage
    localStorage.removeItem('auth_token');
    localStorage.removeItem('auth_user');

    setIsLoading(false);
  };

  const refreshStores = async () => {
    if (user && token) {
      await loadStores(token, user.id);
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