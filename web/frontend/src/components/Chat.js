import React from "react";
import FontAwesome from "@fortawesome/fontawesome-free";

import "./Chat.css";

export default function Chat({ ...rest }) {
  console.log(FontAwesome);
  return (
    <div {...rest} className={`${rest["className"] || ""} chatComponent`}>
      <div className="chatHistory">
        <div className="spacer"></div>
        <div className="messageAuthor">author1</div>
        <div className="messageContents">
          this is a chat message example that is a bit longer
        </div>
        <div className="messageAuthor">author2</div>
        <div className="messageContents">contents2</div>
        <div className="messageAuthor">author3</div>
        <div className="messageContents">contents3</div>
        <div className="messageAuthor">author331234</div>
        <div className="messageContents">contents3</div>
      </div>
      <div className="chatInput">
        <input placeholder="Enter chat message..." type="text"></input>
        <button>
          <i className="fa-solid fa-paper-plane"></i>
        </button>
      </div>
    </div>
  );
}
