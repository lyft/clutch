import React from "react";

import type { Workflow } from "../AppProvider/workflow";

/**
 * Union type representing different lookup keys
 */
export type HeaderItem = "NPS";

interface HeaderItemData {
  /**
   * Optional configuration data to be passed when opening a component
   */
  [key: string]: unknown;
}

export type TriggeredHeaderData = {
  /**
   * The components name, referenced from the above HeaderItems, to be used as a lookup key
   */
  [key in HeaderItem]: HeaderItemData;
};

interface ContextProps {
  workflows: Workflow[];
  /**
   * Will trigger the given header item based on the key, setting the open property as well as saving any optional data
   */
  triggerHeaderItem?: (item: HeaderItem, data?: unknown) => void;
  /**
   * Will return the triggered data, used as a lookup for listening components
   */
  triggeredHeaderData?: TriggeredHeaderData;
}

const ApplicationContext = React.createContext<ContextProps>(undefined);

const useAppContext = () => {
  return React.useContext<ContextProps>(ApplicationContext);
};

export { ApplicationContext, useAppContext };
