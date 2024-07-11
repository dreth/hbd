"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation"; // Import useRouter
import { Input } from "@/components/ui/input";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import Link from "next/link";
import { OctagonAlert, BadgeAlert } from "lucide-react";
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from "@/components/ui/select";
import { Toggle } from "@/components/ui/toggle";
import { generateEncryptionKey, registerUser } from "@/lib/api/apiService";
import { useAuth } from "@/src/context/AuthContext"; // Import useAuth

export default function Register() {
  const { setAuthInfo } = useAuth(); // Use the context to set auth info
  const [email, setEmail] = useState("");
  const [encryptionKey, setEncryptionKey] = useState("");
  const [reminderTime, setReminderTime] = useState("");
  const [timeZone, setTimeZone] = useState("");
  const [telegramApiKey, setTelegramApiKey] = useState("");
  const [telegramUser, setTelegramUser] = useState("");
  const [copySuccess, setCopySuccess] = useState("");
  const [isTimezoneDisabled, setIsTimezoneDisabled] = useState(true);
  const [registerSuccess, setRegisterSuccess] = useState<boolean | null>(null);
  const [registerError, setRegisterError] = useState<string | null>(null); // State for registration error

  const router = useRouter(); // Initialize useRouter

  useEffect(() => {
    const fetchEncryptionKey = async () => {
      try {
        const key = await generateEncryptionKey();
        setEncryptionKey(key);
      } catch (error) {
        console.error("Error fetching encryption key:", error);
      }
    };

    fetchEncryptionKey();

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

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // clear localstorage on submit
    if (typeof window !== "undefined") {
      localStorage.clear();
    }

    const userData = {
      email,
      encryption_key: encryptionKey,
      reminder_time: reminderTime,
      telegram_bot_api_key: telegramApiKey,
      telegram_user_id: telegramUser,
      timezone: timeZone,
    };

    try {
      const response = await registerUser(userData);
      if (response.success) {
        setAuthInfo(email, encryptionKey); // Set the auth info in context
        setRegisterSuccess(true);
        setRegisterError(null);
        setTimeout(() => {
          router.push("/dashboard"); // Redirect to /dashboard after successful registration
        }, 2000); // Delay for user to see the success message
      } else {
        setRegisterSuccess(false);
        setRegisterError(
          response.error || "Registration failed. Please try again."
        );
      }
    } catch (error: any) {
      setRegisterSuccess(false);
      setRegisterError(
        error.response?.data?.error || "Registration failed. Please try again."
      );
      console.error("Error registering user:", error);
    }
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
        <Alert className="max-w-lg mt-3 bg-primary-foreground">
          <OctagonAlert className="h-4 w-4" />
          <AlertTitle className="text-primary">Email Privacy Disclaimer: </AlertTitle>
          <AlertDescription>
            IT IS HASHED BRO WE DONT CARE ABOUT IT
          </AlertDescription>
        </Alert>

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
        <Alert className="max-w-lg my-3 bg-primary-foreground">
          <BadgeAlert className="h-4 w-4" />
          <AlertTitle className="text-primary">Please </AlertTitle>
          <AlertDescription>
            Remember to copy your encryption key before registering!
          </AlertDescription>
        </Alert>

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
            {/* Toggle to enable/disable input */}
            <div className="flex items-center space-x-2">
              <Toggle
                id="toggleTimeZoneInput"
                pressed={!isTimezoneDisabled}
                onPressedChange={handleTimezoneCheckboxChange}
                aria-label="Toggle Edit"
              >
                Edit
              </Toggle>
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
            Telegram User ID
          </label>
          <Input
            id="telegram-user"
            type="text"
            placeholder="Telegram User ID"
            value={telegramUser}
            onChange={(e) => setTelegramUser(e.target.value)}
            className="mt-1 block w-full bg-primary-foreground"
          />
          {telegramApiKey && (
            <Alert className="max-w-lg mt-3 bg-primary-foreground">
              <OctagonAlert className="h-4 w-4" />
              <AlertTitle className="text-primary">
                Need help finding your ID?
              </AlertTitle>
              <AlertDescription>
                Start a conversation with your bot <b>from your mobile phone</b> using <code>/start</code> send a random message to it and then follow this <Link href={`https://api.telegram.org/bot${telegramApiKey}/getUpdates`} className="font-bold underline">link</Link>. You should see a JSON response which will show a numeric ID in several places of the JSON response. That&apos;s your ID!
              </AlertDescription>
            </Alert>
          )}
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
      {registerSuccess !== null && (
        <div className="max-w-lg mt-3 bg-primary-foreground p-4 rounded-lg shadow-md">
          {registerSuccess ? (
            <p className="text-green-600">
              Registration successful! Redirecting to dashboard...
            </p>
          ) : (
            <p className="text-red-600">{registerError}</p>
          )}
        </div>
      )}
    </main>
  );
}
