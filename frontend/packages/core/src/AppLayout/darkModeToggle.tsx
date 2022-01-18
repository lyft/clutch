import * as React from "react";
import { Switch } from "@clutch-sh/core";

const DarkModeToggle = ({ isDarkMode, toggleDarkMode }) => {
  const handleChange = event => {
    toggleDarkMode(event.target.checked);
  };

  return (
    <>
      <span style={{ paddingLeft: "15px", color: "white" }}>Dark Mode</span>
      <Switch checked={isDarkMode} disabled={false} onChange={handleChange} />
    </>
  );
};

export default DarkModeToggle;
