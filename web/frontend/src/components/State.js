import React, { useContext, useEffect, useState } from "react";
import { Socket } from "../lib/socket";

const StateContext = React.createContext(null);

export const useStreamState = () => useContext(StateContext);

export const StreamStateProvider = ({ children }) => {
  const [state, setState] = useState({
    online: true,
  });

  useEffect(() => {
    const stateWS = new Socket(window.location, "api/state?ws=1");
    stateWS.on("state", (state) => setState(state.message));
    stateWS.connect();

    return () => stateWS.disconnect();
  }, []);

  return (
    <StateContext.Provider value={state}>{children}</StateContext.Provider>
  );
};
