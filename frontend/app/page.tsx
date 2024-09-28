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
import { useRouter } from "next/navigation";
import {
  loginUser,
  generatePassword,
  registerUser,
} from "@/lib/api/apiService";
import { useAuth } from "@/src/context/AuthContext";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { OctagonAlert, CircleHelp, BookOpen, Coffee } from "lucide-react";
import { GitHubLogoIcon } from "@radix-ui/react-icons";
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from "@/components/ui/select";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import TelegramApiKeyInput from "@/components/ui/telegram-api-key-input";

export default function Home() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [emailError, setEmailError] = useState("");
  const [passwordError, setPasswordError] = useState("");
  const [loginError, setLoginError] = useState("");
  const [loginSuccess, setLoginSuccess] = useState<boolean | null>(null);
  const [reminderTime, setReminderTime] = useState("");
  const [timeZone, setTimeZone] = useState("");
  const [searchTerm, setSearchTerm] = useState("");
  const [telegramApiKey, setTelegramApiKey] = useState("");
  const [telegramUser, setTelegramUser] = useState("");
  const [copySuccess, setCopySuccess] = useState("");
  const [registerSuccess, setRegisterSuccess] = useState<boolean | null>(null);
  const [registerError, setRegisterError] = useState<string | null>(null);

  const router = useRouter();
  const { setAuthInfo } = useAuth();

  useEffect(() => {
    const userTimeZone = Intl.DateTimeFormat().resolvedOptions().timeZone;
    setTimeZone(userTimeZone);
  }, []);

  const validateEmail = (email: string) => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  };

  const validatePassword = (password: string) => {
    return password.length >= 1;
  };

  const handleLoginSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    let isValid = true;

    if (!validateEmail(email)) {
      setEmailError("Please enter a valid email address.");
      isValid = false;
    } else {
      setEmailError("");
    }

    if (!validatePassword(password)) {
      setPasswordError("Password must be at least 8 characters long.");
      isValid = false;
    } else {
      setPasswordError("");
    }

    if (isValid) {
      try {
        const response = await loginUser({ email, password });
        setAuthInfo(email, response.token);
        setLoginSuccess(true);
        router.push("/dashboard");
      } catch (error) {
        setLoginError("Invalid email or password.");
        setLoginSuccess(false);
        console.error("Login error:", error);
      }
    }
  };

  const handleGeneratePasswordClick = async () => {
    try {
      const generatedPassword = await generatePassword();
      setPassword(generatedPassword);

      navigator.clipboard.writeText(generatedPassword).then(
        () => {
          setCopySuccess("Password generated and copied to clipboard!");
        },
        (err) => {
          setCopySuccess("Failed to copy!");
        }
      );
    } catch (error) {
      console.error("Error generating password", error);
    }
  };

  const handleCopyClick = () => {
    navigator.clipboard.writeText(password).then(
      () => {
        setCopySuccess("Password copied!");
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
      password,
      reminder_time: reminderTime,
      telegram_bot_api_key: telegramApiKey,
      telegram_user_id: telegramUser,
      timezone: timeZone,
    };

    try {
      const response = await registerUser(userData);
      if (response.token) {
        localStorage.setItem("token", response.token);
        localStorage.setItem("email", email);

        setAuthInfo(email, response.token);

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

  const filteredTimeZones = Intl.supportedValuesOf("timeZone").filter((zone) =>
    zone.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <main className="flex min-h-screen flex-col items-center p-5 lg:p-10">
      <div className="col-span-1 min-h-full">
        <h1 className="text-lg md:text-2xl lg:text-4xl font-bold text-center mb-2">
          HBD
        </h1>
        <Tabs defaultValue="login" className="flex flex-col justify-start">
          <TabsList className="flex justify-center">
            <TabsTrigger value="login">Login</TabsTrigger>
            <TabsTrigger value="signup">Sign up</TabsTrigger>
          </TabsList>
          <TabsContent value="login">
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
                  htmlFor="password"
                  className="block text-sm font-medium text-primary"
                >
                  Password
                </label>
                <Input
                  id="password"
                  type="password"
                  placeholder="Password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="mt-1 block w-full bg-primary-foreground dark:bg-background"
                />
                {passwordError && (
                  <p className="text-red-600 text-sm mt-1">{passwordError}</p>
                )}
              </div>
              {loginError && (
                <p className="text-red-600 text-sm mt-1">{loginError}</p>
              )}
              <div className="flex flex-col lg:flex-row items-center justify-between gap-3">
                <TooltipProvider>
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <p className="text-sm text-primary cursor-help">
                        Forgot your password?
                      </p>
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
          </TabsContent>
          <TabsContent value="signup">
            <form
              onSubmit={handleRegisterSubmit}
              className="w-full max-w-md bg-secondary p-8 rounded-lg shadow-md space-y-6"
            >
              <div className="space-y-4">
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
                <div>
                  <label
                    htmlFor="password"
                    className="block text-sm font-medium text-primary"
                  >
                    Password
                    <Popover>
                      <PopoverTrigger>
                        <CircleHelp className="ml-2 text-secondary-foreground w-4 h-4" />
                      </PopoverTrigger>
                      <PopoverContent>
                        <p className="text-primary">
                          You can generate a password or enter your own.
                        </p>
                      </PopoverContent>
                    </Popover>
                  </label>
                  <div className="flex items-center mt-1">
                    <Input
                      id="password"
                      type="text"
                      placeholder="Generated or Entered Password"
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
                      className="block w-full bg-primary-foreground dark:bg-background"
                    />
                  </div>
                  <div className="flex justify-between items-center gap-2">
                    <button
                      type="button"
                      onClick={handleGeneratePasswordClick}
                      className="mt-2 px-3 py-1 w-full bg-blue-600 text-white font-semibold rounded-md shadow-md hover:bg-blue-700 transition duration-300"
                    >
                      Generate
                    </button>
                    <button
                      type="button"
                      onClick={handleCopyClick}
                      className="mt-2 px-3 py-1 w-auto bg-sky-400 text-white font-semibold rounded-md shadow-md hover:bg-sky-700 transition duration-300"
                    >
                      Copy
                    </button>
                  </div>
                  {copySuccess && (
                    <p className="text-sm text-green-600 mt-1">{copySuccess}</p>
                  )}
                </div>
                <div className="flex flex-col md:flex-row md:items-center gap-3">
                  <div>
                    <label
                      htmlFor="reminder-time"
                      className="block text-sm font-medium text-primary whitespace-nowrap"
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
                        <SelectValue
                          placeholder={timeZone || "Select a time zone"}
                        />
                      </SelectTrigger>
                      <SelectContent>
                        <div className="p-2">
                          <input
                            type="text"
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                            placeholder="Search time zones"
                            className="w-full p-2 border rounded"
                          />
                        </div>
                        {filteredTimeZones.map((zone) => (
                          <SelectItem key={zone} value={zone}>
                            {zone}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                </div>{" "}
                <TelegramApiKeyInput></TelegramApiKeyInput>
                <div>
                  <label
                    htmlFor="telegram-user"
                    className="block text-sm font-medium text-primary"
                  >
                    Telegram User ID
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
                        <br />
                        <p>If you use the default bot, message the @RawDataBot on Telegram and it will return you a JSON response with your chat ID under &quot;id&quot; in several document fields, it should be a string of numbers like 637278744</p>
                      </PopoverContent>
                    </Popover>
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
                        . You should see a JSON response which will show a
                        numeric ID in several places of the JSON response.
                        That&apos;s your ID!
                      </AlertDescription>
                    </Alert>
                  )}
                </div>
                <Alert className="max-w-lg mt-3 bg-primary-foreground dark:bg-background">
                  <OctagonAlert className="h-4 w-4" />
                  <AlertTitle className="text-primary">
                    Email and Password Privacy Disclaimer:{" "}
                  </AlertTitle>
                  <AlertDescription>
                    IT IS HASHED BRO WE DON&apos;T CARE ABOUT IT
                  </AlertDescription>
                </Alert>
                <div className="flex justify-center">
                  <button
                    type="submit"
                    className="px-6 py-2 bg-primary w-full text-white font-semibold rounded-md shadow-md hover:bg-blue-700 transition duration-300"
                  >
                    Register
                  </button>
                </div>
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
        <br />
        <hr className="border-primary" />
        <div className="flex justify-evenly mt-4">
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <Link
                  href="/api/swagger/index.html"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  <BookOpen className="w-5 h-5 hover:text-accent" />
                </Link>
              </TooltipTrigger>
              <TooltipContent>
                <p>Swagger</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <Link
                  href="https://github.com/dreth/hbd"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  <GitHubLogoIcon className="w-5 h-5 hover:text-accent" />
                </Link>
              </TooltipTrigger>
              <TooltipContent>
                <p>Github</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        </div>
        <br />
      </div>
    </main>
  );
}
