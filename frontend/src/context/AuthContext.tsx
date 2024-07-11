"use client";


import { createContext, useContext, useState, ReactNode } from 'react';

interface AuthContextType {
  email: string;
  encryptionKey: string;
  setAuthInfo: (email: string, encryptionKey: string) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [email, setEmail] = useState('');
  const [encryptionKey, setEncryptionKey] = useState('');

  const setAuthInfo = (email: string, encryptionKey: string) => {
    setEmail(email);
    setEncryptionKey(encryptionKey);
  };

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
