export default class EventEmitter {
  constructor() {
    this._handlers = {};
  }

  on(t, h) {
    if (this._handlers[t] == null) {
      this._handlers[t] = [h];
    } else {
      this._handlers[t].push(h);
    }
  }

  emit(t, e) {
    if (this._handlers[t] == null) {
      return;
    }

    this._handlers[t].forEach((h) => h(e));
  }
}
