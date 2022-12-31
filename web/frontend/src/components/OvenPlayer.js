import React from "react";
import { useEffect, useRef } from "react";

import Hls from "hls.js";
window.Hls = Hls; // for OvenPlayer

import oven from "ovenplayer";

export default function OvenPlayer({ config, ...rest }) {
  const wrapperEl = useRef(null);
  useEffect(() => {
    console.log("(re)creating player", config);
    const playerEl = document.createElement("div");
    wrapperEl.current.append(playerEl);
    let playerInstance = oven.create(playerEl, config);

    return () => {
      playerInstance.remove();
      playerInstance = null;
    };
  }, [config]);

  return <div {...rest} ref={wrapperEl}></div>;
}
