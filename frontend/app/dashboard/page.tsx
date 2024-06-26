"use client";

import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

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
      setBirthdays([...birthdays, { name, date }]);
      setName("");
      setDate("");
    }
  };
  const handleDeleteBirthday = (index: number) => {
    const newBirthdays = birthdays.filter((_, i) => i !== index);
    setBirthdays(newBirthdays);
  };

  const handleEditBirthday = (index: number) => {
    const birthday = birthdays[index];
    setName(birthday.name);
    setDate(birthday.date);
    handleDeleteBirthday(index);
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
      <div className="flex flex-col lg:flex-row justify-center gap-8 w-full p-2 lg:p-10">
        <div className="w-full lg:w-2/3 bg-secondary p-8 rounded-lg shadow-md space-y-6">
          <h2 className="text-2xl font-semibold">User Information</h2>
          <div>
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
              {/* Checkbox to enable/disable input */}
              <div className="flex items-center space-x-2">
                <Checkbox
                  id="toggleEmailInput"
                  checked={!isEmailDisabled}
                  onCheckedChange={handleEmailCheckboxChange}
                />
                <label
                  htmlFor="toggleEmailInput"
                  className="text-sm font-medium lg:whitespace-nowrap leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                >
                  Want to change email?
                </label>
              </div>
            </div>
          </div>
            {/* Encryption Key field */}
          <div className="flex flex-col lg:flex-row justify-between space-x-4">
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
              onClick={() => setIsEncryptionKeyVisible(!isEncryptionKeyVisible)}
              className="ml-2 px-2 py-1 bg-blue-600 text-white font-semibold rounded-md w-full md:w-fit hover:bg-blue-700 transition duration-300"
            >
              {isEncryptionKeyVisible ? "Hide" : "Show"}
            </button>
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
            {/* Checkbox to enable/disable input */}
            <div className="flex items-center space-x-2">
              <Checkbox
                id="toggleTimezoneInput"
                checked={!isTimezoneDisabled}
                onCheckedChange={() =>
                  setIsTimezoneDisabled(!isTimezoneDisabled)
                }
              />
              <label
                htmlFor="toggleTimezoneInput"
                className="text-sm font-medium lg:whitespace-nowrap leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
              >
                Want to change reminder time?
              </label>
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
            {/* Checkbox to enable/disable input */}
            <div className="flex items-center space-x-2">
              <Checkbox
                id="toggleTimeZoneInput"
                checked={!isTimezoneDisabled}
                onCheckedChange={handleTimezoneCheckboxChange}
              />
              <label
                htmlFor="toggleTimeZoneInput"
                className="text-sm font-medium lg:whitespace-nowrap leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
              >
                Want to change time zone?
              </label>
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
              {/* Checkbox to enable/disable input */}
              <div className="flex items-center space-x-2">
                <Checkbox
                  id="toggleTelegramApiKeyInput"
                  checked={!isTelegramApiKeyDisabled}
                  onCheckedChange={handleTelegramApiKeyCheckboxChange}
                />
                <label
                  htmlFor="toggleTelegramApiKeyInput"
                  className="text-sm font-medium lg:whitespace-nowrap leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                >
                  Want to change Telegram API key?
                </label>
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
              {/* Checkbox to enable/disable input */}
              <div className="flex items-center space-x-2">
                <Checkbox
                  id="toggleTelegramUserInput"
                  checked={!isTelegramUserDisabled}
                  onCheckedChange={handleTelegramUserCheckboxChange}
                />
                <label
                  htmlFor="toggleTelegramUserInput"
                  className="text-sm font-medium lg:whitespace-nowrap leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                >
                  Want to change Telegram user?
                </label>
              </div>
            </div>
          </div>
        </div>
        <div className="w-full lg:w-1/3 bg-secondary p-8 rounded-lg shadow-md space-y-6">
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
              Add Birthday
            </button>
          </form>
          <h2 className="text-2xl font-semibold mt-6">Birthdays</h2>
          <ul className="space-y-4">
            {birthdays.map((birthday, index) => (
              <li
                key={index}
                className="flex flex-col lg:flex-row justify-between bg-primary-foreground p-4 rounded-md shadow-md"
              >
                <div className="flex flex-col lg:flex-row justify-between items-center w-full mr-2">
                  <span>{birthday.name}</span>
                  <span>{birthday.date}</span>
                </div>
                <div className="flex space-x-2">
                  <button
                    onClick={() => handleEditBirthday(index)}
                    className="px-2 py-1 bg-yellow-500 text-white font-semibold rounded-md hover:bg-yellow-600 transition duration-300"
                  >
                    Edit
                  </button>
                  <button
                    onClick={() => handleDeleteBirthday(index)}
                    className="px-2 py-1 bg-red-500 text-white font-semibold rounded-md hover:bg-red-600 transition duration-300"
                  >
                    Delete
                  </button>
                </div>
              </li>
            ))}
          </ul>
        </div>{" "}
      </div>
    </main>
  );
}
