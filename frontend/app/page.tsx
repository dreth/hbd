import Image from "next/image";
import Link from "next/link";

export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-between p-5 lg:p-24">
      <h1 className="text-lg md:text-2xl lg:text-4xl font-bold text-center">
        Welcome to HBD birthday reminder App
      </h1>
      <Image
        src="/replace.webp"
        alt="Birthday cake"
        width={600}
        height={600}
        className="rounded-lg"
      />
      <p className="text-lg text-center">
        Never forget a birthday again with this simple app
      </p>
      <div className="flex space-x-4">
        <Link href="/login">
          <span className="px-6 py-3 bg-primary text-white font-semibold rounded-md shadow-md hover:bg-blue-700 transition duration-300">
            Login
          </span>
        </Link>
        <Link href="/register">
          <span className="px-6 py-3 bg-secondary text-white font-semibold rounded-md shadow-md hover:bg-green-700 transition duration-300">
            Register
          </span>
        </Link>
      </div>

    </main>
  );
}
