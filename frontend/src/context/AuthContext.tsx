"use client";

import { createContext, useContext, useState, ReactNode, useEffect } from 'react';

interface AuthContextType {
  email: string;
  encryptionKey: string;
  setAuthInfo: (email: string, encryptionKey: string) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [email, setEmail] = useState<string | null>(null);
  const [encryptionKey, setEncryptionKey] = useState<string | null>(null);

  const setAuthInfo = (email: string, encryptionKey: string) => {
    setEmail(email);
    setEncryptionKey(encryptionKey);
    localStorage.setItem('email', email);
    localStorage.setItem('encryptionKey', encryptionKey);
  };

  const logout = () => {
    setEmail(null);
    setEncryptionKey(null);
    localStorage.removeItem('email');
    localStorage.removeItem('encryptionKey');
  };

  useEffect(() => {
    const storedEmail = localStorage.getItem('email');
    const storedKey = localStorage.getItem('encryptionKey');
    if (storedEmail && storedKey) {
      setEmail(storedEmail);
      setEncryptionKey(storedKey);
    }
  }, []);

  return (
    <AuthContext.Provider value={{ email: email ?? '', encryptionKey: encryptionKey ?? '', setAuthInfo, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
