import React from "react";
import { createRoot } from "react-dom/client";

import Player from "./components/Player";

const container = document.getElementById("root");
const root = createRoot(container);

root.render(
  <React.StrictMode>
    <div>
      <Player />
      what5
    </div>
  </React.StrictMode>
);
