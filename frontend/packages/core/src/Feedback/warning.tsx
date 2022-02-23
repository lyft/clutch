import React from "react";

import Toast from "./toast";

export interface WarningProps {
  message: any;
  onClose?: () => void;
}

/**
 * Warning component
 * @param message the message to display in the warning
 * @deprecated use Toast component
 */
const Warning: React.FC<WarningProps> = ({ message, onClose }) => {
  if (process.env.NODE_ENV === "development") {
    // eslint-disable-next-line no-console
    console.warn("Warning component is deprecated, please use Toast component instead");
  }

  return (
    <Toast severity="warning" onClose={onClose}>
      {message}
    </Toast>
  );
};

export default Warning;
