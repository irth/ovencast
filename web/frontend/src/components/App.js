import React from "react";
import "./App.css";
import Chat from "./Chat";
import Nav from "./Nav";
import Player from "./Player";
import StreamInfo from "./StreamInfo";

export default function App() {
  return (
    <div className="appLayout">
      <Nav className="nav" brand="OvenCast" />
      <div className="streamColumn">
        <Player className="player" />
        <StreamInfo
          className="streamInfo"
          title="Test streama"
          subtitle="This is just a test stream"
        />
        <div className="streamDescription">
          <div className="streamDescriptionInner">
            <h1>What is OvenCast?</h1>
            <p>
              OvenCast is a distribution of OvenMediaEngine preconfigured for
              low-latency self-hosted livestreaming.
            </p>
            <h2>Features</h2>
            <ul>
              <li>Low latency streaming via WebRTC</li>
              <li>Chat</li>
              <li>stuff</li>
            </ul>
          </div>
        </div>
      </div>
      <Chat className="chat" />
    </div>
  );
}
