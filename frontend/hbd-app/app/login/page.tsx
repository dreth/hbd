"use client"

import { useState } from "react";
import { Input } from "@/components/ui/input";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import Link from "next/link";

export default function Login() {
  const [email, setEmail] = useState("");
  const [encryptionKey, setEncryptionKey] = useState("");
  const [emailError, setEmailError] = useState("");
  const [keyError, setKeyError] = useState("");

  const validateEmail = (email: string) => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  };

  const validateEncryptionKey = (key: string) => {
    return key.length === 64;
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    let isValid = true;

    if (!validateEmail(email)) {
      setEmailError("Please enter a valid email address fam.");
      isValid = false;
    } else {
      setEmailError("");
    }

    if (!validateEncryptionKey(encryptionKey)) {
      setKeyError("Encryption key must be 64 characters long fam.");
      isValid = false;
    } else {
      setKeyError("");
    }

    if (isValid) {
      // Proceed with form submission
      console.log("Form is valid. Submitting...");
    }
  };

  return (
    <main className="flex min-h-screen flex-col items-center p-8">
      <h1 className="text-lg md:text-2xl lg:text-4xl font-bold text-center my-8">
        Login to HBD Reminder App
      </h1>
      <form
        onSubmit={handleSubmit}
        className="w-full max-w-md bg-secondary p-8 rounded-lg shadow-md space-y-6"
      >
        <div>
          <label htmlFor="email" className="block text-sm font-medium text-primary">
            Email
          </label>
          <Input
            id="email"
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="mt-1 block w-full bg-primary-foreground"
          />
          {emailError && <p className="text-red-600 text-sm mt-1">{emailError}</p>}
        </div>
        <div>
          <label htmlFor="encryption-key" className="block text-sm font-medium text-primary">
            Encryption Key
          </label>
          <Input
            id="encryption-key"
            type="password"
            placeholder="Encryption Key"
            value={encryptionKey}
            onChange={(e) => setEncryptionKey(e.target.value)}
            className="mt-1 block w-full bg-primary-foreground"
          />
          {keyError && <p className="text-red-600 text-sm mt-1">{keyError}</p>}
        </div>
        <div className="flex flex-col lg:flex-row items-center justify-between">
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <Link href="/register">
                  <span className="text-sm text-primary cursor-help">Forgot your encryption key?</span>
                </Link>
              </TooltipTrigger>
              <TooltipContent className="bg-destructive">
                <p>gg fam go start over</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
          <button
            type="submit"
            className="px-6 py-2 bg-primary text-white font-semibold rounded-md shadow-md hover:bg-blue-700 transition duration-300"
          >
            Login
          </button>
        </div>
      </form>
      <p className="text-sm text-gray-600 mt-4">
        Dont have an account?{" "}
        <Link href="/register">
          <span className="text-blue-600 hover:underline">Register</span>
        </Link>
      </p>
    </main>
  );
}
