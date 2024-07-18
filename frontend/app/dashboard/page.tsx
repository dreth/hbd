"use client";

import { useState, useEffect } from "react";
import { Input } from "@/components/ui/input";
import { Alert, AlertTitle, AlertDescription } from "@/components/ui/alert";
import { AlertCircle } from "lucide-react";
import { OctagonAlert } from "lucide-react";
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
  modifyUser,
  deleteUser,
  loginUser,
  addBirthday,
  modifyBirthday,
  deleteBirthday,
} from "@/lib/api/apiService";
import { useAuth } from "@/src/context/AuthContext";
import { useRouter } from "next/navigation";

// Add CSS class for ring effect
const ringClass = "ring-2 ring-blue-500";

export default function Dashboard() {
  const { email, encryptionKey } = useAuth();
  const router = useRouter();

  // State to hold user data
  const [userData, setUserData] = useState({
    email: email || "",
    encryptionKey: encryptionKey || "",
    reminderTime: "",
    timeZone: "",
    telegramApiKey: "",
    telegramUser: "",
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
  const [isEncryptionKeyVisible, setIsEncryptionKeyVisible] = useState(false);
  const [isTelegramApiKeyVisible, setIsTelegramApiKeyVisible] = useState(false);
  const [isEmailDisabled, setIsEmailDisabled] = useState(true);
  const [isTelegramApiKeyDisabled, setIsTelegramApiKeyDisabled] =
    useState(true);
  const [isTelegramUserDisabled, setIsTelegramUserDisabled] = useState(true);
  const [isTimezoneDisabled, setIsTimezoneDisabled] = useState(true);
  const [timeZone, setTimeZone] = useState(userData.timeZone);
  const [editIndex, setEditIndex] = useState<number | null>(null); // State to track the index of the birthday being edited
  const [deleteIndex, setDeleteIndex] = useState<number | null>(null); // State to track the index of the birthday to be deleted
  const [confirmDeleteUser, setConfirmDeleteUser] = useState<boolean>(false); // State for delete user confirmation

  // State to hold error messages
  const [nameError, setNameError] = useState<string | null>(null);
  const [dateError, setDateError] = useState<string | null>(null);

  useEffect(() => {
    // Only access localStorage on the client side
    if (typeof window !== "undefined") {
      setUserData((prevData) => ({
        ...prevData,
        email: email || localStorage.getItem("email") || "",
        encryptionKey:
          encryptionKey || localStorage.getItem("encryptionKey") || "",
        reminderTime: localStorage.getItem("reminderTime") || "",
        timeZone: localStorage.getItem("timeZone") || "",
        telegramApiKey: localStorage.getItem("telegramApiKey") || "",
        telegramUser: localStorage.getItem("telegramUser") || "",
      }));
    }
  }, [email, encryptionKey]);

  useEffect(() => {
    // Fetch user data from login payload
    const fetchUserData = async () => {
      try {
        const response = await loginUser({
          email: userData.email,
          encryption_key: userData.encryptionKey,
        });
        setUserData((prevData) => ({
          ...prevData,
          reminderTime: response.reminder_time,
          timeZone: response.timezone,
          telegramApiKey: response.telegram_bot_api_key,
          telegramUser: response.telegram_user_id,
        }));
        setTimeZone(response.timezone);
        localStorage.setItem("reminderTime", response.reminder_time);
        localStorage.setItem("timeZone", response.timezone);
        localStorage.setItem("telegramApiKey", response.telegram_bot_api_key);
        localStorage.setItem("telegramUser", response.telegram_user_id);
        localStorage.setItem("encryptionKey", userData.encryptionKey);

        if (response.birthdays) {
          setBirthdays(response.birthdays);
        }
      } catch (error) {
        console.error("Error fetching user data", error);
      }
    };
    if (userData.email && userData.encryptionKey) {
      fetchUserData();
    }
  }, [userData.email, userData.encryptionKey]);

  // Handlers to toggle the disabled state
  const handleEmailCheckboxChange = () => {
    setIsEmailDisabled(!isEmailDisabled);
  };

  const handleTelegramApiKeyCheckboxChange = () => {
    setIsTelegramApiKeyDisabled(!isTelegramApiKeyDisabled);
  };

  const handleTelegramUserCheckboxChange = () => {
    setIsTelegramUserDisabled(!isTelegramUserDisabled);
  };

  const handleUpdateUser = async () => {
    try {
      await modifyUser({
        auth: {
          email: userData.email,
          encryption_key: userData.encryptionKey,
        },
        new_email: userData.email,
        new_reminder_time: userData.reminderTime,
        new_telegram_bot_api_key: userData.telegramApiKey,
        new_telegram_user_id: userData.telegramUser,
        new_timezone: timeZone,
      });
      console.log("User updated successfully");
      localStorage.setItem("email", userData.email);
      localStorage.setItem("reminderTime", userData.reminderTime);
      localStorage.setItem("telegramApiKey", userData.telegramApiKey);
      localStorage.setItem("telegramUser", userData.telegramUser);
      localStorage.setItem("timeZone", timeZone);
      localStorage.setItem("encryptionKey", userData.encryptionKey);
    } catch (error) {
      console.error("Error updating user", error);
    }
  };

  const handleDeleteUser = async () => {
    if (!confirmDeleteUser) {
      setConfirmDeleteUser(true);
    } else {
      try {
        await deleteUser({
          email: userData.email,
          encryption_key: userData.encryptionKey,
        });
        console.log("User deleted successfully");
        localStorage.clear();
        router.push("/"); // Redirect to home page after deletion
      } catch (error) {
        console.error("Error deleting user", error);
      }
    }
  };

  const handleAddBirthday = async (e: { preventDefault: () => void }) => {
    e.preventDefault();
    // Reset error messages
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
        const response = await addBirthday({
          auth: {
            email: userData.email,
            encryption_key: userData.encryptionKey,
          },
          name,
          date,
        });
        if (response.success) {
          setBirthdays([...birthdays, { id: response.id, name, date }]);
          setName("");
          setDate("");
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
    setEditIndex(index); // Set the edit index to the current birthday
  };

  const handleUpdateBirthday = async (e: { preventDefault: () => void }) => {
    e.preventDefault();
    if (editIndex !== null && name && date) {
      try {
        const response = await modifyBirthday({
          auth: {
            email: userData.email,
            encryption_key: userData.encryptionKey,
          },
          id: birthdays[editIndex].id, // Assuming each birthday has a unique id
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
          setEditIndex(null); // Reset edit index after updating
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
        const response = await deleteBirthday({
          auth: {
            email: userData.email,
            encryption_key: userData.encryptionKey,
          },
          id: birthdayToDelete.id,
          date: birthdayToDelete.date,
          name: birthdayToDelete.name,
        });
        if (response.success) {
          const newBirthdays = birthdays.filter((_, i) => i !== deleteIndex);
          setBirthdays(newBirthdays);
          setDeleteIndex(null); // Reset delete index after deleting
        }
      } catch (error) {
        console.error("Error deleting birthday", error);
      }
    }
  };

  const handleDateChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const inputDate = e.target.value;
    const [year, month, day] = inputDate.split("-");
    if (year === "0000" || year === "0000") {
      // Keep the date as is if the year is set to 0000
      setDate(`0000-${month}-${day}`);
    } else {
      // Otherwise, format it to yyyy-mm-dd
      setDate(inputDate);
    }
  };

  // Get the list of supported time zones
  const timeZones = Intl.supportedValuesOf("timeZone");

  const handleTimezoneCheckboxChange = () => {
    setIsTimezoneDisabled(!isTimezoneDisabled);
  };

  const handleNameChangeBirthday = (name: string) => {
    setName(name);
  };

  return (
    <main className="flex min-h-screen flex-col items-center p-8">
      <h1 className="text-4xl font-bold text-center my-8">Dashboard</h1>
      <div className="flex flex-col justify-center gap-8 w-full max-w-screen-2xl p-2 lg:p-10">
        <Accordion type="single" collapsible>
          <AccordionItem value="item-1" className="bg-secondary rounded-lg px-6">
            <AccordionTrigger className="text-2xl font-semibold">
              User Information
            </AccordionTrigger>
            <AccordionContent>
              <div className="w-full bg-secondary p-8 rounded-lg shadow-md space-y-6">
                {/* Encryption Key field */}
                <div className="flex flex-col lg:flex-row justify-between space-x-4 mb-4">
                  <div>
                    <strong>Encryption Key:</strong>{" "}
                    {isEncryptionKeyVisible ? (
                      <span className="break-all">
                        {userData.encryptionKey}
                      </span>
                    ) : (
                      <span className="break-all">
                        ****************************************************************
                      </span>
                    )}
                  </div>
                  <button
                    onClick={() =>
                      setIsEncryptionKeyVisible(!isEncryptionKeyVisible)
                    }
                    className="ml-2 px-2 py-1 bg-blue-600 text-white font-semibold rounded-md w-full md:w-fit hover:bg-blue-700 transition duration-300"
                  >
                    {isEncryptionKeyVisible ? "Hide" : "Show"}
                  </button>
                </div>
                {/* Email input field */}
                <div className="flex flex-col lg:flex-row justify-between items-center gap-3">
                  <strong>Email:</strong>
                  <Input
                    type="email"
                    placeholder="new email?"
                    value={userData.email}
                    className="bg-primary-foreground"
                    disabled={isEmailDisabled}
                    onChange={(e) =>
                      setUserData({ ...userData, email: e.target.value })
                    }
                  />
                  {/* Toggle to enable/disable input */}
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
                {/* Reminder Time input field */}
                <div className="flex flex-col lg:flex-row justify-between items-center gap-3">
                  <strong className="lg:whitespace-nowrap">
                    Reminder Time:
                  </strong>
                  <Input
                    type="time"
                    placeholder="new reminder time?"
                    value={userData.reminderTime}
                    className="bg-primary-foreground"
                    onChange={(e) =>
                      setUserData({ ...userData, reminderTime: e.target.value })
                    }
                  />
                </div>
                {/* Time zone input field */}
                <div className="flex flex-col lg:flex-row justify_between items-center gap-3">
                  <strong className="lg:whitespace-nowrap">Time Zone:</strong>
                  <Select
                    value={timeZone}
                    onValueChange={setTimeZone}
                    disabled={isTimezoneDisabled}
                  >
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
                {/* Telegram API Key field */}
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
                    className="bg-primary-foreground"
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
                  {/* Toggle to enable/disable input */}
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
                {/* Telegram User ID field */}
                <div className="flex flex-col lg:flex-row justify-between items-center gap-3">
                  <strong className="lg:whitespace-nowrap">
                    Telegram User ID:
                  </strong>
                  <Input
                    type="text"
                    placeholder="new telegram user?"
                    value={userData.telegramUser}
                    className="bg-primary-foreground"
                    disabled={isTelegramUserDisabled}
                    onChange={(e) =>
                      setUserData({ ...userData, telegramUser: e.target.value })
                    }
                  />
                  {/* Toggle to enable/disable input */}
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
                <div className="flex justify-center">
                  <button
                    onClick={handleUpdateUser}
                    className="w-full px-6 py-2 bg-green-600 text-white font-semibold rounded-md shadow-md hover:bg-green-700 transition duration-300"
                  >
                    Update User Data
                  </button>
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
                  className="mt-1 block w-full bg-primary-foreground"
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
                  placeholder="Date"
                  value={date}
                  onChange={handleDateChange}
                  className="mt-1 block w-full bg-primary-foreground"
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
                  } flex flex-col lg:flex-row justify-between bg-primary-foreground p-4 rounded-md shadow-md`}
                >
                  <div className="flex flex-col lg:flex-row justify-between lg:justify-normal items-center w-full mr-2">
                    <span className="font-medium mr-2">{birthday.name}</span>
                    <span>{birthday.date}</span>
                  </div>
                  <div className="flex space-x-2 justify-center">
                    <button
                      onClick={() => handleEditBirthday(index)}
                      className="px-2 py-1 bg-yellow-500 text-white font-semibold rounded-md hover:bg-yellow-600 transition duration-300"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => setDeleteIndex(index)}
                      className="px-2 py-1 bg-red-500 text-white font-semibold rounded-md hover:bg-red-600 transition duration-300"
                    >
                      Delete
                    </button>
                  </div>
                </li>
                {deleteIndex === index && (
                  <Alert
                    variant="destructive"
                    className="mt-2 bg-primary-foreground"
                  >
                    <AlertCircle className="h-4 w-4" />
                    <AlertTitle>Are you sure?</AlertTitle>
                    <AlertDescription>
                      Are you sure you want to delete this birthday?
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
      </div>
      {/* Delete user field */}
      <div className="flex justify-center">
        {!confirmDeleteUser ? (
          <button
            onClick={() => setConfirmDeleteUser(true)}
            className="w-full px-6 py-2 bg-red-600 text-white font-semibold rounded-md shadow-md hover:bg-red-700 transition duration-300"
          >
            Delete User
          </button>
        ) : (
          <div className="flex flex-col justify-center items-center space-y-4">
            <Alert className="max-w-lg mt-3 bg-destructive-foreground">
              <OctagonAlert className="h-4 w-4" />
              <AlertTitle className="text-destructive">
                There is no way back from this{" "}
              </AlertTitle>
              <AlertDescription>
                Please ve careful with this action
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
    </main>
  );
}
