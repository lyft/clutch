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
   * Will trigger the given header item based on the key, saving any optional data
   * This is useful for when a component would like to trigger an item in the header through some state change,
   * such as the NPS component opening or opening a search and prefilling it
   *
   * Example of utilizing with NPS
   * Opening - triggerHeaderItem("NPS", { defaultOption: "Example 1" })
   * Closing - triggerHeaderItem("NPS", null)
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
