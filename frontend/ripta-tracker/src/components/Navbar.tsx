import React from "react"
import { Link } from "react-router-dom";

const Navbar: React.FC = () => {
  return (
    <nav className="fixed top-0 left-0 w-full bg-white shadow-md z-50">
      <div className="max-w-6xl mx-auto px-4">
        <div className="flex justify-between items-center py-3">
          {/* Logo */}
          <Link to="/" className="text-2xl font-bold text-gray-800">
            RIPTA Tracker
          </Link>

          {/* Navigation Links */}
          <div className="flex space-x-6">
            <Link to="/" className="text-gray-600 hover:text-blue-600">
              Home
            </Link>
            <Link to="/map" className="text-gray-600 hover:text-blue-600">
              Live Map
            </Link>
            <Link to="/schedules" className="text-gray-600 hover:text-blue-600">
              Schedules
            </Link>
          </div>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
