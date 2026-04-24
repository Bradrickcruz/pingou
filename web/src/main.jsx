import React, { useState } from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { injectGlobalStyles } from "./theme/globalStyles";
import { Shell } from "./components/layout/Shell";
import { Dashboard } from "./pages/Dashboard";
import { Incidents } from "./pages/Incidents";
import { Settings } from "./pages/Settings";
import { Login } from "./pages/Login";

injectGlobalStyles();

function App() {
  const [apiKey, setApiKey] = useState(
    () => sessionStorage.getItem("pingou_api_key") || "",
  );

  const handleLogin = (key) => setApiKey(key);

  const handleLogout = () => {
    sessionStorage.removeItem("pingou_api_key");
    setApiKey("");
  };

  if (!apiKey) {
    return <Login onLogin={handleLogin} />;
  }

  return (
    <BrowserRouter>
      <Shell onLogout={handleLogout}>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/incidents" element={<Incidents />} />
          <Route path="/settings" element={<Settings />} />
        </Routes>
      </Shell>
    </BrowserRouter>
  );
}

export { App };

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
);
