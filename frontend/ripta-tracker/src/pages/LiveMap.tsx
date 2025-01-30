import React, { useEffect, useState } from "react";
import { MapContainer, TileLayer, Marker, Popup, Polyline } from "react-leaflet";
import "leaflet/dist/leaflet.css";
import { getVehiclePositions, getRoutePolyline } from "../services/api"; 
import L from "leaflet";


const busIcon = L.icon({
  iconUrl: "/bus_dark.svg", 
  iconSize: [20, 20], 
  iconAnchor: [20, 40], 
  popupAnchor: [0, -40], 
});

const LiveMap: React.FC = () => {
  const [vehiclePositions, setVehiclePositions] = useState<any[]>([]);
  const [routePolyline, setRoutePolyline] = useState<[number, number][]>([]);

  // Fetch vehicle positions from Go backend
  useEffect(() => {
    const fetchVehiclePositions = async () => {
      try {
        const vehicles = await getVehiclePositions();
        console.log("Vehicle Positions Response:", vehicles);

        // Extract vehicle data from the response
        const positions = vehicles.entity.map((entity: any) => ({
          id: entity.id,
          label: entity.vehicle.vehicle.label,
          latitude: entity.vehicle.position.latitude,
          longitude: entity.vehicle.position.longitude,
          routeId: entity.vehicle.trip.route_id,
        }));

        console.log("Parsed vehicle positions:", positions)
        setVehiclePositions(positions);
      } catch (error) {
        console.error("Error fetching vehicle positions:", error);
      }
    };

    fetchVehiclePositions();

    // Refresh data every 30 seconds
    const interval = setInterval(fetchVehiclePositions, 30000);
    return () => clearInterval(interval); // Cleanup interval on component unmount
  }, []);

  // Fetch and show route polyline
  const handleRouteClick = async (routeId: string) => {
    try {
      const polyline = await getRoutePolyline(routeId);
      setRoutePolyline(polyline);
    } catch (error) {
      console.error("Could net fetch route polyline")
    }
  }

  return (
  <div className="h-[calc(100vh-4rem)] w-screen mt-16">
    <MapContainer
      center={[41.823989, -71.412834]} // Centered on Rhode Island
      zoom={13}
      className="h-full w-full"
    >
      <TileLayer
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
      />

      {vehiclePositions.map((vehicle, index) => {
        if (!vehicle.latitude || !vehicle.longitude) {
          console.error(`Invalid coordinates for vehicle ID ${vehicle.id}`);
          return null;
        }

        return (
          <Marker
            key={vehicle.id || index}
            position={[vehicle.latitude, vehicle.longitude]}
            icon={busIcon}
          >
            <Popup>
              <div className="text-center">
                <h3 className="text-lg font-bold">{vehicle.label}</h3>
                <p className="text-gray-700">Route ID: {vehicle.routeId}</p>
                <button
                  onClick={() => handleRouteClick(vehicle.routeId)} // Fetch and show route
                  className="mt-2 bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md cursor-pointer transition duration-200"
                >
                  Show Route
                </button>
                <p className="text-gray-500 text-sm">Lat: {vehicle.latitude}</p>
                <p className="text-gray-500 text-sm">Lon: {vehicle.longitude}</p>
              </div>
            </Popup>
          </Marker>
        );
      })}

      {/* Route Polyline */}
      {routePolyline.length > 0 && (
        <Polyline
          positions={routePolyline.filter(([lat, lon]) => lat && lon)} // Filter invalid points
          color="blue"
          weight={5}
        />
      )}
    </MapContainer>
  </div>
  );
};

export default LiveMap;

