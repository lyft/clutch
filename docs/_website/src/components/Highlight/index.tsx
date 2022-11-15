import React from "react";

interface HighlightProps {
  children: React.ReactChildren;
}

const Highlight = ({ children }: HighlightProps): JSX.Element => (
  <span
    style={{
      backgroundColor: "var(--ifm-color-primary)",
      borderRadius: "2px",
      color: "#fff",
      padding: "0.2rem",
    }}
  >
    {children}
  </span>
);

export default Highlight;
