"use client";

import { useState, useEffect } from "react";
import { Input } from "@/components/ui/input";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import Link from "next/link";
import Image from "next/image";
import { useRouter } from "next/navigation";
import {
  loginUser,
  generateEncryptionKey,
  registerUser,
} from "@/lib/api/apiService";
import { useAuth } from "@/src/context/AuthContext";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import { OctagonAlert, BadgeAlert, CircleHelp } from "lucide-react";
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from "@/components/ui/select";
import { Toggle } from "@/components/ui/toggle";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";

export default function Home() {
  const [email, setEmail] = useState("");
  const [encryptionKey, setEncryptionKey] = useState("");
  const [emailError, setEmailError] = useState("");
  const [keyError, setKeyError] = useState("");
  const [loginError, setLoginError] = useState("");
  const [loginSuccess, setLoginSuccess] = useState<boolean | null>(null);

  const [reminderTime, setReminderTime] = useState("");
  const [timeZone, setTimeZone] = useState("");
  const [telegramApiKey, setTelegramApiKey] = useState("");
  const [telegramUser, setTelegramUser] = useState("");
  const [copySuccess, setCopySuccess] = useState("");
  const [isTimezoneDisabled, setIsTimezoneDisabled] = useState(true);
  const [registerSuccess, setRegisterSuccess] = useState<boolean | null>(null);
  const [registerError, setRegisterError] = useState<string | null>(null);

  const router = useRouter();
  const { setAuthInfo } = useAuth();

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

    const userTimeZone = Intl.DateTimeFormat().resolvedOptions().timeZone;
    setTimeZone(userTimeZone);
  }, []);

  const validateEmail = (email: string) => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  };

  const validateEncryptionKey = (key: string) => {
    return key.length === 64;
  };

  const handleLoginSubmit = async (e: React.FormEvent) => {
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
      try {
        const response = await loginUser({
          email,
          encryption_key: encryptionKey,
        });
        setAuthInfo(email, encryptionKey);
        setLoginSuccess(true);
        console.log("Login successful:", response);
        router.push("/dashboard");
      } catch (error) {
        setLoginError("Invalid email or encryption key.");
        setLoginSuccess(false);
        console.error("Login error:", error);
      }
    }
  };

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

  const handleRegisterSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

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
        setAuthInfo(email, encryptionKey);
        setRegisterSuccess(true);
        setRegisterError(null);
        setTimeout(() => {
          router.push("/dashboard");
        }, 2000);
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

  const timeZones = Intl.supportedValuesOf("timeZone");

  const handleTimezoneCheckboxChange = () => {
    setIsTimezoneDisabled(!isTimezoneDisabled);
  };

  return (
    <main className="flex min-h-screen flex-col items-center p-5 lg:p-10">
      <div className="col-span-1 min-h-full">
        <h1 className="text-lg md:text-2xl lg:text-4xl font-bold text-center mb-2">
          HBD
        </h1>
        <Tabs defaultValue="login" className="">
          <TabsList className="flex justify-center bg-background">
              <TabsTrigger value="login">Login</TabsTrigger>
              <TabsTrigger value="signup">Sign up</TabsTrigger>
          </TabsList>
          <TabsContent value="login" className="w-[600px]">
            <form
              onSubmit={handleLoginSubmit}
              className="w-full max-w-md bg-secondary p-8 rounded-lg shadow-md space-y-6"
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
                  className="mt-1 block w-full bg-primary-foreground dark:bg-background"
                />
                {emailError && (
                  <p className="text-red-600 text-sm mt-1">{emailError}</p>
                )}
              </div>
              <div>
                <label
                  htmlFor="encryption-key"
                  className="block text-sm font-medium text-primary"
                >
                  Encryption Key
                </label>
                <Input
                  id="encryption-key"
                  type="password"
                  placeholder="Encryption Key"
                  value={encryptionKey}
                  onChange={(e) => setEncryptionKey(e.target.value)}
                  className="mt-1 block w-full bg-primary-foreground dark:bg-background"
                />
                {keyError && (
                  <p className="text-red-600 text-sm mt-1">{keyError}</p>
                )}
              </div>
              {loginError && (
                <p className="text-red-600 text-sm mt-1">{loginError}</p>
              )}
              <div className="flex flex-col lg:flex-row items-center justify-between">
                <TooltipProvider>
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Link href="/register">
                        <span className="text-sm text-primary cursor-help">
                          Forgot your encryption key?
                        </span>
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
              Don&apos;t have an account?{" "}
              <Link href="/register">
                <span className="text-blue-600 hover:underline">Register</span>
              </Link>
            </p>
          </TabsContent>
          <TabsContent value="signup" className="w-[600px]">
            <form
              onSubmit={handleRegisterSubmit}
              className="w-full max-w-md bg-secondary p-8 rounded-lg shadow-md space-y-6"
            >
              <div className="space-y-4">
                <h3 className="font-medium text-primary">
                  Copy your encryption key before registering!
                </h3>
                <div className="flex flex-col md:flex-row items-center mt-1">
                  <Input
                    id="encryption-key"
                    type="text"
                    placeholder="Generated Encryption Key"
                    value={encryptionKey}
                    readOnly
                    className="block w-full bg-primary-foreground dark:bg-background"
                  />
                  <button
                    type="button"
                    onClick={handleCopyClick}
                    className="ml-2 px-3 py-1 mt-1 lg:mt-0 w-full lg:w-auto bg-blue-600 text-white font-semibold rounded-md shadow-md hover:bg-blue-700 transition duration-300"
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
                  className="mt-1 block w-full bg-primary-foreground dark:bg-background"
                />
              </div>
              <div className="flex items-center space-x-3">
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
                    className="mt-1 block w-full bg-primary-foreground dark:bg-background"
                  />
                </div>
                <div className="space-y-1 w-full">
                  <label
                    htmlFor="reminder-time"
                    className="block text-sm font-medium text-primary"
                  >
                    Time Zone
                  </label>
                  <Select onValueChange={setTimeZone}>
                    <SelectTrigger className="bg-primary-foreground dark:bg-background">
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
                </div>
              </div>
              <div>
                <label
                  htmlFor="telegram-api-key"
                  className="text-sm font-medium text-primary flex items-center"
                >
                  Telegram Bot API Key
                  <Popover>
                    <PopoverTrigger>
                      <CircleHelp className="ml-2 text-secondary-foreground w-4 h-4" />
                    </PopoverTrigger>
                    <PopoverContent>
                      <p className="text-primary text-lg">
                        Need help finding your ID?
                      </p>
                      <p>
                        Start a conversation with your bot{" "}
                        <b>from your mobile phone</b> using{" "}
                        <code className="bg-blue-100 dark:bg-primary p-0.5 rounded-md">
                          /start
                        </code>{" "}
                        send a random message to it and then follow this{" "}
                        <Link
                          href={`https://api.telegram.org/bot${telegramApiKey}/getUpdates`}
                          className="font-bold underline"
                        >
                          link
                        </Link>
                        . You should see a JSON response which will show a
                        numeric ID in several places of the JSON response.
                        That&apos;s your ID!
                      </p>
                    </PopoverContent>
                  </Popover>
                </label>
                <Input
                  id="telegram-api-key"
                  type="text"
                  placeholder="Telegram Bot API Key"
                  value={telegramApiKey}
                  onChange={(e) => setTelegramApiKey(e.target.value)}
                  className="mt-1 block w-full bg-primary-foreground dark:bg-background"
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
                  className="mt-1 block w-full bg-primary-foreground dark:bg-background"
                />
                {telegramApiKey && (
                  <Alert className="max-w-lg mt-3 bg-primary-foreground dark:bg-background">
                    <OctagonAlert className="h-4 w-4" />
                    <AlertTitle className="text-primary">
                      Need help finding your ID?
                    </AlertTitle>
                    <AlertDescription>
                      Start a conversation with your bot{" "}
                      <b>from your mobile phone</b> using{" "}
                      <code className="bg-blue-100 dark:bg-primary p-0.5 rounded-md">
                        /start
                      </code>{" "}
                      send a random message to it and then follow this{" "}
                      <Link
                        href={`https://api.telegram.org/bot${telegramApiKey}/getUpdates`}
                        className="font-bold underline"
                      >
                        link
                      </Link>
                      . You should see a JSON response which will show a numeric
                      ID in several places of the JSON response. That&apos;s
                      your ID!
                    </AlertDescription>
                  </Alert>
                )}
              </div>
              <Alert className="max-w-lg mt-3 bg-primary-foreground dark:bg-background">
                <OctagonAlert className="h-4 w-4" />
                <AlertTitle className="text-primary">
                  Email Privacy Disclaimer:{" "}
                </AlertTitle>
                <AlertDescription>
                  IT IS HASHED BRO WE DON&apos;T CARE ABOUT IT
                </AlertDescription>
              </Alert>
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
              <div className="p-4 ">
                {registerSuccess ? (
                  <p className="text-green-600">
                    Registration successful! Redirecting to dashboard...
                  </p>
                ) : (
                  <p className="text-red-600">{registerError}</p>
                )}
              </div>
            )}
          </TabsContent>
        </Tabs>
      </div>
    </main>
  );
}
