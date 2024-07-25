"use client";

import React, { createContext, useContext, useState, ReactNode } from 'react';

interface AuthContextType {
  email: string | null;
  token: string | null;  
  setAuthInfo: (email: string, token: string) => void; 
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [email, setEmail] = useState<string | null>(null);
  const [token, setToken] = useState<string | null>(null);  

  const setAuthInfo = (email: string, token: string) => {
    setEmail(email);
    setToken(token);  
  };

  const logout = () => {
    setEmail(null);
    setToken(null);  
    localStorage.removeItem('token');
  };

  return (
    <AuthContext.Provider value={{ email, token, setAuthInfo, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
