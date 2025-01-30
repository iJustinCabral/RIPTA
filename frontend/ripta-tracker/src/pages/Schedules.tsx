import React, { useEffect, useState } from "react";
import { getSchedule, getRoutes } from "../services/api"; 

const Schedule: React.FC = () => {
  const [routes, setRoutes] = useState<any[]>([]);
  const [selectedRoute, setSelectedRoute] = useState<string>("");
  const [schedule, setSchedule] = useState<any[]>([]);
  const [visibleStops, setVisibleStops] = useState(10);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchRoutes = async () => {
      try {
        const routeData = await getRoutes();
        setRoutes(routeData);
        if (routeData.length > 0) {
          setSelectedRoute(routeData[0].route_id);
        }
      } catch (err) {
        console.error("Error fetching routes:", err);
      }
    };

    fetchRoutes();
  }, []);

  useEffect(() => {
    if (!selectedRoute) return;

    const fetchSchedule = async () => {
      setLoading(true);
      setError("");

      try {
        const response = await getSchedule(selectedRoute);

        // Get current time in HH:MM:SS format
        const currentTime = new Date();
        const currentTimeString = currentTime.toTimeString().split(" ")[0];

        // Filter stops to show only upcoming stops
        const upcomingStops = response.filter((stop: any) => stop.ArrivalTime > currentTimeString);

        setSchedule(upcomingStops);
      } catch (err) {
        setError("Failed to load schedule.");
        console.error("Error fetching schedule:", err);
      } finally {
        setLoading(false);
      }
    };

    fetchSchedule();
  }, [selectedRoute]); // ✅ Fetches new schedule when route changes

  const loadMoreStops = () => {
    setVisibleStops((prev) => prev + 10);
  };

  return (
    <div className="max-w-3xl mx-auto p-6 bg-white shadow-lg rounded-md">
      <h2 className="text-2xl font-semibold text-center text-gray-800 mb-4">
        Bus Schedule
      </h2>

      {/* Route Selector */}
      <div className="mb-4">
        <label className="block text-gray-700 font-semibold mb-2">Select Route:</label>
        <select
          className="w-full border-gray-300 rounded-md p-2 text-gray-700 bg-white shadow-sm"
          value={selectedRoute}
          onChange={(e) => setSelectedRoute(e.target.value)}
        >
          {routes.map((route) => (
            <option key={route.route_id} value={route.route_id}>
              {route.route_id} - {route.route_long_name}
            </option>
          ))}
        </select>
      </div>

      {/* Loading & Error Messages */}
      {loading && <p className="text-center text-gray-500">Loading schedule...</p>}
      {error && <p className="text-center text-red-500">{error}</p>}

      {/* Schedule List */}
      {!loading && !error && schedule.length === 0 && (
        <p className="text-center text-gray-500">No more upcoming stops for today.</p>
      )}

      {!loading && !error && schedule.length > 0 && (
        <>
          <div className="space-y-4">
            {schedule.slice(0, visibleStops).map((stop, index) => (
              <div key={index} className="border p-3 rounded-md bg-gray-50">
                <h3 className="font-semibold text-lg text-gray-700">{stop.stop_name}</h3>
                <p className="text-gray-600">🕒 Arrival: {stop.ArrivalTime}</p>
                <p className="text-gray-600">🕒 Departure: {stop.DepartureTime}</p>
              </div>
            ))}
          </div>

          {visibleStops < schedule.length && (
            <button
              onClick={loadMoreStops}
              className="mt-4 w-full bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md transition duration-200"
            >
              Load More Stops
            </button>
          )}
        </>
      )}
    </div>
  );
};

export default Schedule;

