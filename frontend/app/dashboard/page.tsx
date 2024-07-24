"use client";

import Link from "next/link";
import { useState, useEffect } from "react";
import { Input } from "@/components/ui/input";
import { Alert, AlertTitle, AlertDescription } from "@/components/ui/alert";
import {
  AlertCircle,
  OctagonAlert,
  Cake,
  Gift,
  GitPullRequestArrow,
  BookOpen,
  Coffee,
} from "lucide-react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { Toggle } from "@/components/ui/toggle";
import {
  getUserData,
  modifyUser,
  deleteUser,
  addBirthday,
  modifyBirthday,
  deleteBirthday,
} from "@/lib/api/apiService";
import { useAuth } from "@/src/context/AuthContext";
import { useRouter } from "next/navigation";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";

const ringClass = "ring-2 ring-blue-500";

export default function Dashboard() {
  const { email, token, logout } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!email || !token) {
      router.push("/");
    }
  }, [email, token, router]);

  const [userData, setUserData] = useState({
    email: email || "",
    reminderTime: "",
    timeZone: "",
    telegramApiKey: "",
    telegramUser: "",
    newPassword: "", // Added new password field
  });

  const [name, setName] = useState("");
  const [date, setDate] = useState("");
  const [birthdays, setBirthdays] = useState<
    {
      id: any;
      name: string;
      date: string;
    }[]
  >([]);
  const [isTelegramApiKeyVisible, setIsTelegramApiKeyVisible] = useState(false);
  const [isEmailDisabled, setIsEmailDisabled] = useState(true);
  const [isPasswordDisabled, setIsPasswordDisabled] = useState(true);
  const [isTelegramApiKeyDisabled, setIsTelegramApiKeyDisabled] =
    useState(true);
  const [isTelegramUserDisabled, setIsTelegramUserDisabled] = useState(true);
  const [isTimezoneDisabled, setIsTimezoneDisabled] = useState(true);
  const [isReminderTimeDisabled, setIsReminderTimeDisabled] = useState(true);
  const [timeZone, setTimeZone] = useState(userData.timeZone);
  const [editIndex, setEditIndex] = useState<number | null>(null);
  const [deleteIndex, setDeleteIndex] = useState<number | null>(null);
  const [confirmDeleteUser, setConfirmDeleteUser] = useState<boolean>(false);

  const [nameError, setNameError] = useState<string | null>(null);
  const [dateError, setDateError] = useState<string | null>(null);

  const getIcon = (index: number) => {
    return index % 2 === 0 ? (
      <Cake className="h-6 w-6 text-accent" />
    ) : (
      <Gift className="h-6 w-6 text-primary" />
    );
  };

  useEffect(() => {
    const fetchUserData = async () => {
      try {
        const response = await getUserData(token);
        setUserData((prevData) => ({
          ...prevData,
          reminderTime: response.reminder_time,
          timeZone: response.timezone,
          telegramApiKey: response.telegram_bot_api_key,
          telegramUser: response.telegram_user_id,
        }));
        setTimeZone(response.timezone);
        if (response.birthdays) {
          setBirthdays(response.birthdays);
        }
      } catch (error) {
        console.error("Error fetching user data", error);
      }
    };
    if (token) {
      fetchUserData();
    }
  }, [token]);

  const handleEmailCheckboxChange = () => {
    setIsEmailDisabled(!isEmailDisabled);
  };

  const handlePasswordCheckboxChange = () => {
    setIsPasswordDisabled(!isPasswordDisabled);
  };

  const handleTelegramApiKeyCheckboxChange = () => {
    setIsTelegramApiKeyDisabled(!isTelegramApiKeyDisabled);
  };

  const handleTelegramUserCheckboxChange = () => {
    setIsTelegramUserDisabled(!isTelegramUserDisabled);
  };

  const handleUpdateUser = async () => {
    try {
      const response = await modifyUser(token, {
        new_email: userData.email,
        new_password: userData.newPassword, // Include new password
        new_reminder_time: userData.reminderTime,
        new_telegram_bot_api_key: userData.telegramApiKey,
        new_telegram_user_id: userData.telegramUser,
        new_timezone: timeZone,
      });
      console.log("User updated successfully");
      setUserData((prevData) => ({
        ...prevData,
        token: response.token,
      }));
      localStorage.setItem("token", response.token); // Update token in local storage
    } catch (error) {
      console.error("Error updating user", error);
    }
  };

  const handleDeleteUser = async () => {
    if (!confirmDeleteUser) {
      setConfirmDeleteUser(true);
    } else {
      try {
        await deleteUser(token);
        console.log("User deleted successfully");
        logout();
        router.push("/");
      } catch (error) {
        console.error("Error deleting user", error);
      }
    }
  };

  const handleAddBirthday = async (e: { preventDefault: () => void }) => {
    e.preventDefault();
    setNameError(null);
    setDateError(null);

    if (!name) {
      setNameError("Name is required.");
    }
    if (!date) {
      setDateError("Date is required.");
    }
    if (name && date) {
      try {
        const response = await addBirthday(token, {
          name,
          date,
        });
        console.log("Add birthday response:", response);
        if (response.id && response.name && response.date) {
          setBirthdays([
            ...birthdays,
            { id: response.id, name: response.name, date: response.date },
          ]);
          setName("");
          setDate("");
        } else {
          console.error("Failed to add birthday: Invalid response format");
        }
      } catch (error) {
        console.error("Error adding birthday", error);
      }
    }
  };

  const handleEditBirthday = (index: number) => {
    const birthday = birthdays[index];
    setName(birthday.name);
    setDate(birthday.date);
    setEditIndex(index);
  };

  const handleUpdateBirthday = async (e: { preventDefault: () => void }) => {
    e.preventDefault();
    if (editIndex !== null && name && date) {
      try {
        const response = await modifyBirthday(token, {
          id: birthdays[editIndex].id,
          name,
          date,
        });
        if (response.success) {
          const updatedBirthdays = [...birthdays];
          updatedBirthdays[editIndex] = {
            id: birthdays[editIndex].id,
            name,
            date,
          };
          setBirthdays(updatedBirthdays);
          setEditIndex(null);
          setName("");
          setDate("");
        }
      } catch (error) {
        console.error("Error modifying birthday", error);
      }
    }
  };

  const handleDeleteBirthday = async () => {
    if (deleteIndex !== null) {
      try {
        const birthdayToDelete = birthdays[deleteIndex];
        const response = await deleteBirthday(token, {
          id: birthdayToDelete.id,
          date: birthdayToDelete.date,
          name: birthdayToDelete.name,
        });
        if (response.success) {
          setBirthdays((prevBirthdays) =>
            prevBirthdays.filter((_, i) => i !== deleteIndex)
          );
          setDeleteIndex(null);
        }
      } catch (error) {
        console.error("Error deleting birthday", error);
      }
    }
  };

  const handleDateChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setDate(e.target.value);
  };

  const timeZones = Intl.supportedValuesOf("timeZone");

  const handleTimezoneCheckboxChange = () => {
    setIsTimezoneDisabled(!isTimezoneDisabled);
  };

  const handleReminderTimeCheckboxChange = () => {
    setIsReminderTimeDisabled(!isReminderTimeDisabled);
  };

  const handleNameChangeBirthday = (name: string) => {
    setName(name);
  };

  return (
    <main className="flex min-h-screen flex-col items-center p-8">
      <h1 className="text-4xl font-bold text-center my-8">Dashboard</h1>
      <div className="flex flex-col justify-center gap-8 w-full max-w-screen-2xl p-2 lg:p-10">
        {/* User Info Section */}
        <Accordion type="single" collapsible>
          <AccordionItem
            value="item-1"
            className="bg-secondary rounded-lg px-6"
          >
            <AccordionTrigger className="text-2xl font-semibold">
              User Information
            </AccordionTrigger>
            <AccordionContent>
              <div className="w-full bg-secondary lg:p-8 rounded-lg space-y-6">
                <div className="flex flex-col lg:flex-row justify-between items-center gap-3">
                  <strong>Email:</strong>
                  <Input
                    type="email"
                    placeholder="new email?"
                    value={userData.email}
                    className="bg-primary-foreground dark:bg-background"
                    disabled={isEmailDisabled}
                    onChange={(e) =>
                      setUserData({ ...userData, email: e.target.value })
                    }
                  />
                  <div className="flex items-center space-x-2">
                    <Toggle
                      id="toggleEmailInput"
                      pressed={!isEmailDisabled}
                      onPressedChange={handleEmailCheckboxChange}
                      aria-label="Toggle Edit"
                    >
                      Edit
                    </Toggle>
                  </div>
                </div>
                <div className="flex flex-col lg:flex-row justify-between items-center gap-3">
                  <strong className="whitespace-nowrap">New Password:</strong>
                  <Input
                    type="password"
                    placeholder="Just leave it blank if you don't want to change it"
                    className="bg-primary-foreground dark:bg-background"
                    disabled={isPasswordDisabled}
                    onChange={(e) =>
                      setUserData({ ...userData, newPassword: e.target.value })
                    }
                  />
                  <div className="flex items-center space-x-2">
                    <Toggle
                      id="togglePasswordInput"
                      pressed={!isPasswordDisabled}
                      onPressedChange={handlePasswordCheckboxChange}
                      aria-label="Toggle Edit"
                    >
                      Edit
                    </Toggle>
                  </div>
                </div>
                <div className="flex flex-col lg:flex-row justify-between items-center gap-3">
                  <strong className="lg:whitespace-nowrap">
                    Reminder Time:
                  </strong>
                  <Input
                    type="time"
                    placeholder="new reminder time?"
                    value={userData.reminderTime}
                    className="bg-primary-foreground dark:bg-background"
                    disabled={isReminderTimeDisabled}
                    onChange={(e) =>
                      setUserData({ ...userData, reminderTime: e.target.value })
                    }
                  />
                  <div className="flex items-center space-x-2">
                    <Toggle
                      id="toggleTimeZoneInput"
                      pressed={!isReminderTimeDisabled}
                      onPressedChange={handleReminderTimeCheckboxChange}
                      aria-label="Toggle Edit"
                    >
                      Edit
                    </Toggle>
                  </div>
                </div>
                <div className="flex flex-col lg:flex-row justify_between items-center gap-3">
                  <strong className="lg:whitespace-nowrap">Time Zone:</strong>
                  <Select
                    value={timeZone}
                    onValueChange={setTimeZone}
                    disabled={isTimezoneDisabled}
                  >
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
                <div className="flex flex-col lg:flex-row justify-between items-center gap-3">
                  <strong className="lg:whitespace-nowrap">
                    Telegram Bot API Key:
                  </strong>
                  <Input
                    type="text"
                    placeholder="new telegram API key?"
                    value={
                      isTelegramApiKeyVisible
                        ? userData.telegramApiKey
                        : "************"
                    }
                    className="bg-primary-foreground dark:bg-background"
                    disabled={isTelegramApiKeyDisabled}
                    onChange={(e) =>
                      setUserData({
                        ...userData,
                        telegramApiKey: e.target.value,
                      })
                    }
                  />
                  <button
                    onClick={() =>
                      setIsTelegramApiKeyVisible(!isTelegramApiKeyVisible)
                    }
                    className="ml-2 px-2 py-1 bg-blue-600 text-white font-semibold rounded-md hover:bg-blue-700 transition duration-300"
                  >
                    {isTelegramApiKeyVisible ? "Hide" : "Show"}
                  </button>
                  <div className="flex items-center space-x-2">
                    <Toggle
                      id="toggleTelegramApiKeyInput"
                      pressed={!isTelegramApiKeyDisabled}
                      onPressedChange={handleTelegramApiKeyCheckboxChange}
                      aria-label="Toggle Edit"
                    >
                      Edit
                    </Toggle>
                  </div>
                </div>
                <div className="flex flex-col lg:flex-row justify-between items-center gap-3">
                  <strong className="lg:whitespace-nowrap">
                    Telegram User ID:
                  </strong>
                  <Input
                    type="text"
                    placeholder="new telegram user?"
                    value={userData.telegramUser}
                    className="bg-primary-foreground dark:bg-background"
                    disabled={isTelegramUserDisabled}
                    onChange={(e) =>
                      setUserData({ ...userData, telegramUser: e.target.value })
                    }
                  />
                  <div className="flex items-center space-x-2">
                    <Toggle
                      id="toggleTelegramUserInput"
                      pressed={!isTelegramUserDisabled}
                      onPressedChange={handleTelegramUserCheckboxChange}
                      aria-label="Toggle Edit"
                    >
                      Edit
                    </Toggle>
                  </div>
                </div>
                <div className="flex justify-center gap-2">
                  <button
                    onClick={handleUpdateUser}
                    className="w-full max-h-9 px-6 py-2 bg-primary text-white font-semibold rounded-md shadow-md hover:bg-sky-700 transition duration-300"
                  >
                    Update User Data
                  </button>
                  {!confirmDeleteUser ? (
                    <button
                      onClick={() => setConfirmDeleteUser(true)}
                      className="w-full max-w-40 px-6 py-2 mx-auto bg-destructive text-white font-semibold rounded-md shadow-md hover:bg-red-700 transition duration-300"
                    >
                      Delete User
                    </button>
                  ) : (
                    <div className="flex flex-col justify-center items-center space-y-4">
                      <Alert className="max-w-lg mt-3">
                        <OctagonAlert className="h-4 w-4" />
                        <AlertTitle className="text-destructive">
                          There is no way back from this{" "}
                        </AlertTitle>
                        <AlertDescription>
                          Please be careful with this action, all the
                          information will be deleted and you wont be able to
                          recover it.
                        </AlertDescription>
                      </Alert>
                      <div className="flex flex-col lg:flex-row justify-between items-center gap-6">
                        <button
                          onClick={handleDeleteUser}
                          className="w-full lg:w-64 px-3 py-2 bg-red-700 text-white font-semibold rounded-md shadow-md hover:bg-red-900 transition duration-300 whitespace-nowrap"
                        >
                          Confirm Delete
                        </button>
                        <button
                          onClick={() => setConfirmDeleteUser(false)}
                          className="w-full lg:w-64 px-3 py-2 bg-gray-600 text-white font-semibold rounded-md shadow-md hover:bg-gray-700 transition duration-300"
                        >
                          Cancel
                        </button>
                      </div>
                    </div>
                  )}
                </div>
              </div>
            </AccordionContent>
          </AccordionItem>
        </Accordion>

        {/* Birthdays Section */}
        <div className="w-full bg-secondary p-8 rounded-lg shadow-md space-y-6">
          <h2 className="text-2xl font-semibold">Add Birthday</h2>
          <form
            onSubmit={
              editIndex !== null ? handleUpdateBirthday : handleAddBirthday
            }
            className="space-y-4"
          >
            <div className="flex flex-col lg:flex-row space-x-0 lg:space-x-4 space-y-2 lg:space-y-0">
              <div className="w-full lg:w-1/2">
                <label
                  htmlFor="name"
                  className="block text-sm font-medium text-primary"
                >
                  Name
                </label>
                <Input
                  id="name"
                  type="text"
                  placeholder="Name"
                  value={name}
                  onChange={(e) => handleNameChangeBirthday(e.target.value)}
                  className="mt-1 block w-full bg-primary-foreground dark:bg-background"
                />
                {nameError && (
                  <p className="text-red-600 text-sm mt-1">{nameError}</p>
                )}
              </div>
              <div className="w-full lg:w-1/2">
                <label
                  htmlFor="date"
                  className="block text-sm font-medium text-primary"
                >
                  Date
                </label>
                <Input
                  id="date"
                  type="date"
                  placeholder="yyyy-mm-dd"
                  value={date}
                  onChange={handleDateChange}
                  className="mt-1 block w-full bg-primary-foreground dark:bg-background"
                />
                {dateError && (
                  <p className="text-red-600 text-sm mt-1">{dateError}</p>
                )}
              </div>
            </div>
            <button
              type="submit"
              className="px-6 py-2 bg-blue-600 text-white font-semibold rounded-md shadow-md w-full hover:bg-blue-700 transition duration-300"
            >
              {editIndex !== null ? "Update Birthday" : "Add Birthday"}
            </button>
          </form>
          <h2 className="text-2xl font-semibold mt-6">Birthdays</h2>
          <ul className="gap-4 grid grid-cols-2">
            {birthdays.map((birthday, index) => (
              <div key={index} className="col-span-2 md:col-span-1">
                <li
                  className={`${
                    editIndex === index ? ringClass : ""
                  } flex flex-col lg:flex-row justify-between bg-background p-4 rounded-md shadow-md`}
                >
                  <div className="flex flex-col lg:flex-row justify-between  items-center w-full mr-2">
                    <TooltipProvider>
                      <Tooltip>
                        <TooltipTrigger
                          asChild
                          className="truncate font-medium w-40"
                        >
                          <p>{birthday.name}</p>
                        </TooltipTrigger>
                        <TooltipContent>
                          <p>{birthday.name}</p>
                        </TooltipContent>
                      </Tooltip>
                    </TooltipProvider>
                    {getIcon(index)}
                    <span>{birthday.date.split("-").reverse().join("/")}</span>
                  </div>
                  <div className="flex space-x-2 justify-center">
                    <button
                      onClick={() => handleEditBirthday(index)}
                      className="px-2 py-1 bg-amber-600 text-white font-semibold rounded-md hover:bg-amber-800 transition duration-300"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => setDeleteIndex(index)}
                      className="px-2 py-1 bg-rose-500 text-white font-semibold rounded-md hover:bg-rose-600 transition duration-300"
                    >
                      Delete
                    </button>
                  </div>
                </li>
                {deleteIndex === index && (
                  <Alert variant="destructive" className="mt-2 bg-background">
                    <AlertCircle className="h-4 w-4" />
                    <AlertTitle>
                      Are you sure you want to delete this birthday?
                    </AlertTitle>
                    <AlertDescription>
                      There is no way back from this action.
                      <div className="mt-4 flex justify-end space-x-2">
                        <button
                          onClick={() => setDeleteIndex(null)}
                          className="px-4 py-2 bg-gray-500 text-white font-semibold rounded-md hover:bg-gray-600 transition duration-300"
                        >
                          Cancel
                        </button>
                        <button
                          onClick={handleDeleteBirthday}
                          className="px-4 py-2 bg-red-600 text-white font-semibold rounded-md hover:bg-red-700 transition duration-300"
                        >
                          Delete
                        </button>
                      </div>
                    </AlertDescription>
                  </Alert>
                )}
              </div>
            ))}
          </ul>
        </div>

        <hr className="w-full border-primary" />
        <div className="flex justify-evenly mt-4">
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <Link
                  href="https://swagger.io/"
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
                  href="https://github.com/"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  <GitPullRequestArrow className="w-5 h-5 hover:text-accent" />
                </Link>
              </TooltipTrigger>
              <TooltipContent>
                <p>Github</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <Link
                  href="https://buymeacoffee.com/"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  <Coffee className="w-5 h-5 hover:text-primary animate-pulse" />
                </Link>
              </TooltipTrigger>
              <TooltipContent>
                <p>Buy me a coffee!</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        </div>
      </div>
    </main>
  );
}
