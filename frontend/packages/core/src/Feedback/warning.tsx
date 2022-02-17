import React from "react";

import type { ToastProps } from "./toast";
import Toast from "./toast";

export interface WarningProps extends ToastProps {
  message: any;
}

/**
 * Warning component
 * @param message the message to display in the warning
 * @deprecated use Toast component
 */
const Warning: React.FC<WarningProps> = ({ message, ...props }) => {
  if (process.env.NODE_ENV === "development") {
    // eslint-disable-next-line no-console
    console.warn("Warning component is deprecated, please use Toast component instead");
  }

  return (
    <Toast severity="warning" {...props}>
      {message}
    </Toast>
  );
};

export default Warning;
