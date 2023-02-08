import React from "react";
import { StreamState } from "../lib/api";

import OvenPlayer from "./OvenPlayer";

export default function Player(props) {
  const host = location.host;
  const tls = location.protocol == "https:";
  const config = {
    autoStart: true,
    autoFallback: true,
    mute: true,
    sources: [
      {
        type: "webrtc (abr)",
        file: `${tls ? "wss" : "ws"}://${host}/live/stream/webrtc_abr`,
      },
      {
        type: "hls (abr)",
        file: `${tls ? "https" : "http"}://${host}/live/stream/llhls_abr.m3u8`,
      },
    ],
    hlsConfig: {
      liveSyncDuration: 2,
      liveMaxLatencyDuration: 5,
      maxLiveSyncPlaybackRate: 2,
    },
  };
  return (
    <StreamState>
      {
        ({ online }) =>
          online ? (
            <OvenPlayer {...props} config={config} />
          ) : (
            "stream offline TODO: better offline screen"
          ) // TODO: better offline screen
      }
    </StreamState>
  );
}
