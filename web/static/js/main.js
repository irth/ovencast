OvenPlayer.debug(true);

const host = location.host;
const tls = location.protocol == "https";

const player = OvenPlayer.create("player", {
  autoStart: true,
  autoFallback: true,
  mute: true,
  sources: [
    {
      type: "webrtc",
      file: `${tls ? "wss" : "ws"}://${host}/live/stream`,
    },
    {
      type: "hls",
      file: `${tls ? "https" : "http"}://${host}/live/stream/llhls.m3u8`,
    },
  ],
});
