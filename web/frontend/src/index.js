import React from "react";
import { createRoot } from "react-dom/client";

import "@fortawesome/fontawesome-free/js/all.js";
import "@fortawesome/fontawesome-free/css/all.css";

import App from "./components/App";

const container = document.getElementById("root");
const root = createRoot(container);

root.render(
  <React.StrictMode>
    <div>
      <App />
    </div>
  </React.StrictMode>
);
