import EventEmitter from "./events";

export class Socket extends EventEmitter {
  constructor(url, path) {
    super();
    this.ws = null;
    const parsedUrl = new URL(url);
    const pathname = parsedUrl.pathname.replace(/\/+$/, "");
    const wsProto = parsedUrl.protocol == "https:" ? "wss" : "ws";
    this.wsURL = `${wsProto}://${parsedUrl.host}${pathname}/${path}`;
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
}
