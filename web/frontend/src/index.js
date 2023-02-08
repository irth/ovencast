import React from "react";
import { createRoot } from "react-dom/client";

import "@fortawesome/fontawesome-free/js/all.js";
import "@fortawesome/fontawesome-free/css/all.css";
import "@fontsource/lato";

import "./index.css";

import App from "./components/App";
import { APIProvider } from "./lib/api";

const container = document.getElementById("root");
const root = createRoot(container);

root.render(
  <React.StrictMode>
    <div>
      <APIProvider url={`${location.protocol}//${location.host}/api`}>
        <App />
      </APIProvider>
    </div>
  </React.StrictMode>
);
