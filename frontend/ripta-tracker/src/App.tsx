import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import Navbar from "./components/Navbar";
import Home from "./pages/Home";
import LiveMap from "./pages/LiveMap";
import Schedules from "./pages/Schedules";

const App: React.FC = () => {
  return (
    <Router>
      <Navbar />
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/map" element={<LiveMap />} />
        <Route path="/schedules" element={<Schedules />} />
      </Routes>
    </Router>
  );
};

export default App;

