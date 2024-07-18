"use client";

import { createContext, useContext, useState, ReactNode, useEffect } from "react";

interface AuthContextType {
  email: string;
  encryptionKey: string;
  setAuthInfo: (email: string, encryptionKey: string) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [email, setEmail] = useState('');
  const [encryptionKey, setEncryptionKey] = useState('');

  useEffect(() => {
    const storedEmail = localStorage.getItem('email');
    const storedEncryptionKey = localStorage.getItem('encryptionKey');
    if (storedEmail && storedEncryptionKey) {
      setEmail(storedEmail);
      setEncryptionKey(storedEncryptionKey);
    }
  }, []);

  const setAuthInfo = (email: string, encryptionKey: string) => {
    setEmail(email);
    setEncryptionKey(encryptionKey);
    localStorage.setItem('email', email);
    localStorage.setItem('encryptionKey', encryptionKey);
  };

  const logout = () => {
    setEmail('');
    setEncryptionKey('');
    localStorage.removeItem('email');
    localStorage.removeItem('encryptionKey');
  };

  return (
    <AuthContext.Provider value={{ email, encryptionKey, setAuthInfo, logout }}>
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
