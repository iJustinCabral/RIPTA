# RIPTA Real-Time Transit Tracker

This project is a **real-time transit tracker** for RIPTA (Rhode Island Public Transit Authority). It provides live data on bus locations, routes, stop details, and service alerts using RIPTA's real-time GTFS APIs.

## Features In Development
- **Live Map of Bus Locations**: Real-time vehicle positions displayed on an interactive map.
- **Route Overlays**: Visualization of bus routes with live data.
- **Stop Details**: Estimated arrival times for each stop.
- **Service Alerts**: Notifications for service disruptions or changes.

## Tech Stack
### Backend
- **Language**: Go
- **Framework**: Go’s built-in `net/http` package for dependency-free routing.
- **API Integration**: Fetches real-time data from RIPTA's GTFS APIs:
  - Trip Updates: [http://realtime.ripta.com:81/api/tripupdates?format=json](http://realtime.ripta.com:81/api/tripupdates?format=json)
  - Vehicle Positions: [http://realtime.ripta.com:81/api/vehiclepositions?format=json](http://realtime.ripta.com:81/api/vehiclepositions?format=json)
  - Service Alerts: [http://realtime.ripta.com:81/api/servicealerts?format=json](http://realtime.ripta.com:81/api/servicealerts?format=json)

### Frontend
- **Framework**: ReactJS (with TypeScript for type safety).
- **CSS Framework**: Tailwind CSS for responsive and modern styling.
- **Mapping Library**: `react-leaflet@4.2.0` for integrating interactive maps with stable React 18 support.
- **HTTP Client**: Axios for handling API requests.

### Tools
- **Package Manager**: npm for dependency management.
- **Development Environment**: Vite for fast build and development.
- **Version Control**: GitHub for repository hosting.

---

This README is meant to provide an overview of the work completed so far. Future updates will include details about fetching real-time data, dynamic map markers, and frontend/backend integration.


