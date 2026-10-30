import React, { useState } from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import "./index.css";
import { injectGlobalStyles } from "./theme/globalStyles";
import { ConnectionProvider } from "./hooks/useConnection";
import { ConnectionOverlay } from "./components/ui/ConnectionOverlay";
import { Shell } from "./components/layout/Shell";
import { Dashboard } from "./pages/Dashboard";
import { Incidents } from "./pages/Incidents";
import { Settings } from "./pages/Settings";
import { MonitorDetail } from "./pages/MonitorDetail";
import { Login } from "./pages/Login";

injectGlobalStyles();

function App() {
  const [apiKey, setApiKey] = useState(
    () => localStorage.getItem("pingou_api_key") || "",
  );

  const handleLogin = (key) => setApiKey(key);

  const handleLogout = () => {
    localStorage.removeItem("pingou_api_key");
    setApiKey("");
  };

  if (!apiKey) {
    return <Login onLogin={handleLogin} />;
  }

  return (
    <ConnectionProvider>
      <ToastContainer position="bottom-right" autoClose={3000} />
      <ConnectionOverlay>
        <BrowserRouter>
          <Shell onLogout={handleLogout}>
            <Routes>
              <Route path="/" element={<Dashboard />} />
              <Route path="/monitors/:id" element={<MonitorDetail />} />
              <Route path="/incidents" element={<Incidents />} />
              <Route path="/settings" element={<Settings onLogout={handleLogout} />} />
            </Routes>
          </Shell>
        </BrowserRouter>
      </ConnectionOverlay>
    </ConnectionProvider>
  );
}

export { App };

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
);
