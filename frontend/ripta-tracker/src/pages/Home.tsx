import React from "react"

const Home: React.FC = () => {
  return (
    <div className="flex flex-col items-center justify-center h-screen bg-gray-100 text-center px-4 pt-20">
      {/* Title */}
      <h1 className="text-4xl font-bold text-gray-900 mb-4">
        Welcome to RIPTA Tracker
      </h1>

      {/* Subtitle */}
      <p className="text-lg text-gray-700 mb-6 max-w-2xl">
        Get real-time updates on RIPTA buses, routes, and schedules.
      </p>

      {/* AI Chatbox Placeholder */}
      <div className="w-full max-w-lg bg-white shadow-lg rounded-lg p-6">
        <p className="text-gray-500">Ask "Where’s my bus?"</p>
        <input
          type="text"
          placeholder="Coming soon..."
          className="w-full border p-2 mt-2 rounded-md"
          disabled
        />
      </div>
    </div>
  );
};

export default Home;

