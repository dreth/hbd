"use client";

import { useState, useEffect } from "react";
import { Input } from "@/components/ui/input";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import Link from "next/link";
import { OctagonAlert } from "lucide-react";
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from "@/components/ui/select";
import { Checkbox } from "@/components/ui/checkbox";

export default function Register() {
  const [email, setEmail] = useState("");
  const [encryptionKey, setEncryptionKey] = useState(
    "Generated Encryption Key"
  );
  const [reminderTime, setReminderTime] = useState("");
  const [timeZone, setTimeZone] = useState("");
  const [telegramApiKey, setTelegramApiKey] = useState("");
  const [telegramUser, setTelegramUser] = useState("");
  const [copySuccess, setCopySuccess] = useState("");
  const [isTimezoneDisabled, setIsTimezoneDisabled] = useState(true);

  useEffect(() => {
    // Detect and set the user's timezone
    const userTimeZone = Intl.DateTimeFormat().resolvedOptions().timeZone;
    setTimeZone(userTimeZone);
  }, []);

  const handleCopyClick = () => {
    navigator.clipboard.writeText(encryptionKey).then(
      () => {
        setCopySuccess("It's copied fam go save it somewhere safe!");
      },
      (err) => {
        setCopySuccess("Failed to copy!");
      }
    );
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // Add form submission logic here
  };

  // Get the list of supported time zones
  const timeZones = Intl.supportedValuesOf("timeZone");

  const handleTimezoneCheckboxChange = () => {
    setIsTimezoneDisabled(!isTimezoneDisabled);
  };

  return (
    <main className="flex min-h-screen flex-col items-center p-2 lg:p-8">
      <h1 className="text-lg md:text-2xl lg:text-4xl font-bold text-center my-8">
        Register for HBD Reminder App
      </h1>
      <form
        onSubmit={handleSubmit}
        className="w-full max-w-lg bg-secondary p-8 rounded-lg shadow-md space-y-6"
      >
        <div>
          <label
            htmlFor="email"
            className="block text-sm font-medium text-primary"
          >
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
        </div>
        <div>
          <label
            htmlFor="encryption-key"
            className="block text-sm font-medium text-primary"
          >
            Encryption Key
          </label>
          <div className="flex flex-col md:flex-row items-center mt-1">
            <Input
              id="encryption-key"
              type="text"
              placeholder="Generated Encryption Key"
              value={encryptionKey}
              readOnly
              className="block w-full bg-primary-foreground"
            />
            <button
              type="button"
              onClick={handleCopyClick}
              className="ml-2 px-3 py-1 mt-1 lg:mt-0 w-full lg:w-fit bg-primary text-white font-semibold rounded-md shadow-md hover:bg-blue-700 transition duration-300"
            >
              Copy
            </button>
          </div>
          {copySuccess && (
            <p className="text-sm text-green-600 mt-1">{copySuccess}</p>
          )}
        </div>
        <div>
          <label
            htmlFor="reminder-time"
            className="block text-sm font-medium text-primary"
          >
            Reminder Time
          </label>
          <Input
            id="reminder-time"
            type="time"
            placeholder="Reminder Time"
            value={reminderTime}
            onChange={(e) => setReminderTime(e.target.value)}
            className="mt-1 block w-full bg-primary-foreground"
          />
        </div>
        <div className="space-y-2">
          <label
            htmlFor="reminder-time"
            className="block text-sm font-medium text-primary"
          >
            Time Zone
          </label>
          <div className="flex items-center space-x-2">
            <Select onValueChange={setTimeZone} disabled={isTimezoneDisabled}>
              <SelectTrigger className="bg-primary-foreground">
                <SelectValue placeholder={timeZone} />
              </SelectTrigger>
              <SelectContent>
                {timeZones.map((zone) => (
                  <SelectItem key={zone} value={zone}>
                    {zone}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            {/* Checkbox to enable/disable input */}
            <div className="flex items-center space-x-2">
              <Checkbox
                id="toggleTimeZoneInput"
                checked={!isTimezoneDisabled}
                onCheckedChange={handleTimezoneCheckboxChange}
              />
              <label
                htmlFor="toggleTimeZoneInput"
                className="text-sm font-medium whitespace-nowrap leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
              >
                Want to change time zone?
              </label>
            </div>
          </div>
        </div>
        <div>
          <label
            htmlFor="telegram-api-key"
            className="block text-sm font-medium text-primary"
          >
            Telegram Bot API Key
          </label>
          <Input
            id="telegram-api-key"
            type="text"
            placeholder="Telegram Bot API Key"
            value={telegramApiKey}
            onChange={(e) => setTelegramApiKey(e.target.value)}
            className="mt-1 block w-full bg-primary-foreground"
          />
        </div>
        <div>
          <label
            htmlFor="telegram-user"
            className="block text-sm font-medium text-primary"
          >
            Telegram User
          </label>
          <Input
            id="telegram-user"
            type="text"
            placeholder="Telegram User"
            value={telegramUser}
            onChange={(e) => setTelegramUser(e.target.value)}
            className="mt-1 block w-full bg-primary-foreground"
          />
        </div>
        <div className="flex flex-col md:flex-row items-center justify-between">
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <Link href="/login">
                  <span className="text-sm text-primary">
                    Already have an account?
                  </span>
                </Link>
              </TooltipTrigger>
              <TooltipContent>
                <p>Click here to login</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
          <button
            type="submit"
            className="px-6 py-2 bg-primary w-full lg:w-fit text-white font-semibold rounded-md shadow-md hover:bg-blue-700 transition duration-300"
          >
            Register
          </button>
        </div>
      </form>
      <Alert className="max-w-lg mt-3 bg-primary-foreground">
        <OctagonAlert className="h-4 w-4" />
        <AlertTitle>Email Privacy Disclaimer: </AlertTitle>
        <AlertDescription>
          IT IS HASHED BRO WE DONT CARE ABOUT IT
        </AlertDescription>
      </Alert>
    </main>
  );
}
