import * as React from "react";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { CircleHelp } from "lucide-react";
import { Input } from "@/components/ui/input";

function TelegramApiKeyInput() {
  const [telegramApiKey, setTelegramApiKey] = React.useState("");
  const [useDefaultBot, setUseDefaultBot] = React.useState(false);

  // Retrieve the default bot API key from the environment variable
  const defaultBotApiKey = process.env.HBD_DEFAULT_BOT_API_KEY || "";

  // Handle checkbox change
  const handleCheckboxChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setUseDefaultBot(e.target.checked);
    if (e.target.checked) {
      setTelegramApiKey(defaultBotApiKey);
    } else {
      setTelegramApiKey("");
    }
  };

  return (
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
            <p className="text-primary text-lg">Need help finding your API key?</p>
            <p>Follow these steps to get your Telegram API key:</p>
            <ol className="list-decimal ml-4">
              <li>Open Telegram.</li>
              <li>Search for <b>BotFather</b>.</li>
              <li>Start a chat with <b>BotFather</b>.</li>
              <li>
                Use the{" "}
                <code className="bg-blue-100 dark:bg-primary p-0.5 rounded-md">
                  /newbot
                </code>{" "}
                command to create a new bot.
              </li>
              <li>Follow the instructions to create a new bot.</li>
              <li>Copy the API key.</li>
            </ol>
            <p>This API key allows the application to send messages to you through the bot.</p>
            <br />
            <p>Alternatively, you can check the &quot;Use default bot&quot; checkbox and the default bot will be used (if configured by the host)</p>
          </PopoverContent>
        </Popover>
      </label>
      
      {/* Checkbox to toggle between using default bot and custom API key */}
      <div className="mt-2 flex items-center">
        <input
          id="use-default-bot"
          type="checkbox"
          checked={useDefaultBot}
          onChange={handleCheckboxChange}
          className="mr-2"
        />
        <label htmlFor="use-default-bot" className="text-sm text-primary">
          Use the default bot
        </label>
      </div>

      {/* Conditionally render the input field based on the checkbox state */}
      {!useDefaultBot && (
        <Input
          id="telegram-api-key"
          type="text"
          placeholder="Telegram Bot API Key"
          value={telegramApiKey}
          onChange={(e) => setTelegramApiKey(e.target.value)}
          className="mt-1 block w-full bg-primary-foreground dark:bg-background"
        />
      )}
    </div>
  );
}

export default TelegramApiKeyInput;
