import React from "react";

import "./Nav.css";

export default function Nav({ brand, ...rest }) {
  return (
    <nav {...rest} className={`${rest["className"] || ""} navComponent`}>
      <div class="brand">{brand}</div>
      <div class="nav">
        <a href="/login">log in</a>
      </div>
    </nav>
  );
}
