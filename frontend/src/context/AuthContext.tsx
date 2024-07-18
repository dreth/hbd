"use client";

import { createContext, useContext, useState, ReactNode, useEffect } from 'react';

interface AuthContextType {
  email: string;
  encryptionKey: string;
  setAuthInfo: (email: string, encryptionKey: string) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [email, setEmail] = useState(localStorage.getItem('email') || '');
  const [encryptionKey, setEncryptionKey] = useState(localStorage.getItem('encryptionKey') || '');

  const setAuthInfo = (email: string, encryptionKey: string) => {
    setEmail(email);
    setEncryptionKey(encryptionKey);
    localStorage.setItem('email', email);
    localStorage.setItem('encryptionKey', encryptionKey);
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
    <AuthContext.Provider value={{ email, encryptionKey, setAuthInfo }}>
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
