import React, { useEffect, useState } from "react";
import { MapContainer, TileLayer, Marker, Popup } from "react-leaflet";
import "leaflet/dist/leaflet.css";
import { getVehiclePositions } from "./services/api"; 
import L from "leaflet";


const busIcon = L.icon({
  iconUrl: "/bus_dark.svg", 
  iconSize: [20, 20], 
  iconAnchor: [20, 40], 
  popupAnchor: [0, -40], 
});

const App: React.FC = () => {
  const [vehiclePositions, setVehiclePositions] = useState<any[]>([]);

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

  return (
    <div style={{ height: "100vh", width: "100vw" }}>
      <MapContainer
        center={[41.823989, -71.412834]} // Centered on Rhode Island
        zoom={13}
        style={{ height: "100%", width: "100%" }}
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
                <div>
                  <h3>Vehicle: {vehicle.label}</h3>
                  <p>Route ID: {vehicle.routeId}</p>
                  <p>Latitude: {vehicle.latitude}</p>
                  <p>Longitude: {vehicle.longitude}</p>
                </div>
              </Popup>
            </Marker>
          );
        })}
      </MapContainer>
    </div>
  );
};

export default App;

