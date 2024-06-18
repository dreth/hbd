"use client";

import { useState } from "react";
import { Input } from "@/components/ui/input";

export default function Dashboard() {
  // Mock user data for demonstration purposes
  const userData = {
    email: "user@example.com",
    encryptionKey:
      "GeneratedEncryptionKey1234567890123456789012345678901234567890123456789012345678901234",
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

  const handleAddBirthday = (e: React.FormEvent) => {
    e.preventDefault();
    if (name && date) {
      setBirthdays([...birthdays, { name, date }]);
      setName("");
      setDate("");
    }
  };

  return (
    <main className="flex min-h-screen flex-col items-center p-8">
      <h1 className="text-4xl font-bold text-center my-8">Dashboard</h1>
      <div className="flex flex-wrap justify-center gap-8 w-full max-w-7xl">
        <div className="w-full md:w-1/2 bg-secondary p-8 rounded-lg shadow-md space-y-6">
          <h2 className="text-2xl font-semibold">User Information</h2>
          <p>
            <strong>Email:</strong> {userData.email}
          </p>
          <p>
            <strong>Encryption Key:</strong>{" "}
            <span className="break-all">{userData.encryptionKey}</span>
          </p>
          <p>
            <strong>Reminder Time:</strong> {userData.reminderTime}
          </p>
          <p>
            <strong>Time Zone:</strong> {userData.timeZone}
          </p>
          <p>
            <strong>Telegram Bot API Key:</strong> {userData.telegramApiKey}
          </p>
          <p>
            <strong>Telegram User:</strong> {userData.telegramUser}
          </p>
        </div>

        <div className="w-full md:w-1/2 bg-secondary p-8 rounded-lg shadow-md space-y-6">
          <h2 className="text-2xl font-semibold">Add Birthday</h2>
          <form onSubmit={handleAddBirthday} className="space-y-4">
            <div className="flex space-x-4">
              <div className="w-1/2">
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
              <div className="w-1/2">
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
                  onChange={(e) => setDate(e.target.value)}
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
                className="flex justify-between bg-primary-foreground p-4 rounded-md shadow-md"
              >
                <span>{birthday.name}</span>
                <span>{birthday.date}</span>
              </li>
            ))}
          </ul>
        </div>
      </div>
    </main>
  );
}
