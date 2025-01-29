import axios from "axios";

// This will change to a server URL later when deployed
const API_BASE_URL = "http://localhost:8080/api";

export const getTripUpdates = async () => {
  const response = await axios.get(`${API_BASE_URL}/tripupdates`);
  return response.data;
};

export const getVehiclePositions = async () => {
  const response = await axios.get(`${API_BASE_URL}/vehiclepositions`);
  return response.data;
};

export const getServiceAlerts = async () => {
  const response = await axios.get(`${API_BASE_URL}/servicealerts`);
  return response.data;
};

export const getRoutePolyline = async (routeId: string) => {
  try {
    const response = await axios.get(`${API_BASE_URL}/route/${routeId}`);
    return response.data
      .filter((point: { lat: number; lon: number }) => point.lat && point.lon)
      .map((point: { lat: number; lon: number }) => [point.lat, point.lon]);
  } catch (error) {
    console.error("Error fetching route polyline:", error);
    return []
  }
};
