"use client";

import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Alert, AlertTitle, AlertDescription } from "@/components/ui/alert";
import { AlertCircle } from "lucide-react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Toggle } from "@/components/ui/toggle";

// Add CSS class for ring effect
const ringClass = "ring-2 ring-blue-500";

export default function Dashboard() {
  // Mock user data for demonstration purposes
  const userData = {
    email: "user@example.com",
    encryptionKey:
      "1234567890123456789012345678901234567890123456789012345678901234",
    reminderTime: "08:00",
    timeZone: "UTC",
    telegramApiKey: "your-telegram-api-key",
    telegramUser: "your-telegram-user",
  };

  // State to handle new birthday inputs
  const [name, setName] = useState("");
  const [date, setDate] = useState("");
  const [birthdays, setBirthdays] = useState([
    { name: "John Doe", date: "2024-06-20" },
    { name: "Jane Smith", date: "2024-07-15" },
  ]);
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

  const handleAddBirthday = (e: { preventDefault: () => void }) => {
    e.preventDefault();
    if (name && date) {
      if (editIndex !== null) {
        // If editing, update the existing birthday
        const updatedBirthdays = [...birthdays];
        updatedBirthdays[editIndex] = { name, date };
        setBirthdays(updatedBirthdays);
        setEditIndex(null); // Reset edit index after updating
      } else {
        // If adding a new birthday
        setBirthdays([...birthdays, { name, date }]);
      }
      setName("");
      setDate("");
    }
  };

  const handleDeleteBirthday = () => {
    if (deleteIndex !== null) {
      const newBirthdays = birthdays.filter((_, i) => i !== deleteIndex);
      setBirthdays(newBirthdays);
      setDeleteIndex(null); // Reset delete index after deleting
    }
  };

  const handleEditBirthday = (index: number) => {
    const birthday = birthdays[index];
    setName(birthday.name);
    setDate(birthday.date);
    setEditIndex(index); // Set the edit index to the current birthday
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

  return (
    <main className="flex min-h-screen flex-col items-center p-8">
      <h1 className="text-4xl font-bold text-center my-8">Dashboard</h1>
      <div className="flex flex-col justify-center gap-8 w-full max-w-screen-2xl p-2 lg:p-10">
        <div className="w-full bg-secondary p-8 rounded-lg shadow-md space-y-6">
          <h2 className="text-2xl font-semibold">Add Birthday</h2>
          <form onSubmit={handleAddBirthday} className="space-y-4">
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
                  onChange={(e) => setName(e.target.value)}
                  className="mt-1 block w-full bg-primary-foreground"
                />
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
                  <Alert variant="destructive" className="mt-2 bg-primary-foreground">
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
        <div className="w-full bg-secondary p-8 rounded-lg shadow-md space-y-6">
          <h2 className="text-2xl font-semibold">User Information</h2>
          <div>
            {/* Encryption Key field */}
            <div className="flex flex-col lg:flex-row justify-between space-x-4 mb-4">
              <div>
                <strong>Encryption Key:</strong>{" "}
                {isEncryptionKeyVisible ? (
                  <span className="break-all">{userData.encryptionKey}</span>
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
          </div>
          <div className="flex flex-col lg:flex-row justify-between items-center gap-3">
            <strong className="lg:whitespace-nowrap">Reminder Time:</strong>
            <Input
              type="time"
              placeholder="new reminder time?"
              value={userData.reminderTime}
              className="bg-primary-foreground"
              disabled={isTimezoneDisabled}
            />
            {/* Toggle to enable/disable input */}
            <div className="flex items-center space-x-2">
              <Toggle
                id="toggleTimezoneInput"
                pressed={!isTimezoneDisabled}
                onPressedChange={() =>
                  setIsTimezoneDisabled(!isTimezoneDisabled)
                }
                aria-label="Toggle Edit"
              >
                Edit
              </Toggle>
            </div>
          </div>
          <div className="flex flex-col lg:flex-row justify-between items-center gap-3">
            <strong className="lg:whitespace-nowrap">Time Zone:</strong>
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
          </div>{" "}
          <div>
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
          </div>
          <div>
            {/* Telegram User field */}
            <div className="flex flex-col lg:flex-row justify-between items-center gap-3">
              <strong className="lg:whitespace-nowrap">Telegram User:</strong>
              <Input
                type="text"
                placeholder="new telegram user?"
                value={userData.telegramUser}
                className="bg-primary-foreground"
                disabled={isTelegramUserDisabled}
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
          </div>
        </div>
      </div>
    </main>
  );
}
