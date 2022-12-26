OvenPlayer.debug(true);
const player = OvenPlayer.create("player", {
  autoStart: true,
  autoFallback: true,
  mute: true,
  sources: [
    {
      type: "webrtc",
      file: "ws://localhost:3333/live/stream",
    },
    {
      type: "hls",
      file: "http://localhost:3333/live/stream/llhls.m3u8",
    },
  ],
});
