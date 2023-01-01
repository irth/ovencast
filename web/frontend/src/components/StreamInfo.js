import React from "react";

import "./StreamInfo.css";

export default function StreamInfo({ title, subtitle, ...rest }) {
  return (
    <div {...rest} className={`${rest["className"] || ""} streamInfoComponent`}>
      <div class="title">{title}</div>
      <div class="subtitle">{subtitle}</div>
    </div>
  );
}
