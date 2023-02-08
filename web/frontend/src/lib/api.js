import React, { createContext, useEffect, useState, useContext } from "react";
import EventEmitter from "./events";

class APIError extends Error {
  constructor(message) {
    super(message);
    this.name = "APIError";
  }
}

const StateContext = createContext(null);

export function APIProvider({ url, children }) {
  const [streamState, setStreamState] = useState({
    offline: true,
  });

  useEffect(() => {
    const api = new API(url);
    api.connect();
    api.state().then((state) => {
      setStreamState(state);
      console.log(state);
    });

    api.on("state", (e) => {
      setStreamState(e.message);
    });

    return () => {
      api.disconnect();
    };
  }, []);

  return (
    <StateContext.Provider value={streamState}>
      {children}
    </StateContext.Provider>
  );
}

export function StreamState({ children }) {
  const streamState = useContext(StateContext);
  return children(streamState);
}

export class API extends EventEmitter {
  constructor(url) {
    super();
    this.ws = null;
    this.url = url.replace(/\/+$/, "");

    const parsedUrl = new URL(this.url);
    const wsProto = parsedUrl.protocol == "https:" ? "wss" : "ws";
    this.wsURL = `${wsProto}://${parsedUrl.host}${parsedUrl.pathname}/websocket`;
  }

  connect() {
    // TODO: handle reconnections
    if (this.ws != null) {
      this.ws.close();
      this.ws = null;
    }

    console.log("connecting to", this.wsURL);
    this.ws = new WebSocket(this.wsURL);
    this.ws.addEventListener("message", this._wshandler.bind(this));
    // TODO: handle onerror
  }

  disconnect() {
    if (this.ws != null) {
      console.log("disconnecting from", this.wsURL);
      this.ws.close();
    }
  }

  _wshandler(e) {
    const command = JSON.parse(e.data);
    const type = command.type;
    const message = command.message;

    console.log(type, message);

    this.emit(type, { type, message });
  }

  async _req(endpoint) {
    const v = await fetch(`${this.url}${endpoint}`);
    const res = await v.json();
    if (!res.ok) throw new APIError(res.error);
    return res.response;
  }

  async state() {
    return await this._req("/state");
  }
}
