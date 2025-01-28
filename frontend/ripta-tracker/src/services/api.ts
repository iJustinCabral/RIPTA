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
