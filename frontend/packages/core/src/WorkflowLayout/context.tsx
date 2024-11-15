import React from "react";
import { useLocation } from "react-router-dom";

export interface WorkflowLayoutContextProps {
  title?: string;
  subtitle?: string;
  headerContent?: React.ReactNode;
  setTitle: (title: string) => void;
  setSubtitle: (subtitle: string) => void;
  setHeaderContent: (headerContent: React.ReactNode) => void;
}

const INITIAL_STATE = {
  title: null,
  subtitle: null,
  headerContent: null,
  setTitle: () => {},
  setSubtitle: () => {},
  setHeaderContent: () => {},
};

const WorkflowLayoutContext = React.createContext<WorkflowLayoutContextProps>(INITIAL_STATE);

const workflowLayoutContextReducer = (state, action) => {
  switch (action.type) {
    case "set_title":
      return { ...state, title: action.payload };
    case "set_subtitle":
      return { ...state, subtitle: action.payload };
    case "set_content":
      return { ...state, headerContent: action.payload };
    default:
      throw new Error("Unhandled action type");
  }
};

const WorkflowLayoutContextProvider = ({ children }: { children: React.ReactNode }) => {
  const [state, dispatch] = React.useReducer(workflowLayoutContextReducer, INITIAL_STATE);

  const providerValue = React.useMemo(
    () => ({
      ...state,
      setTitle: (title: string) => {
        dispatch({ type: "set_title", payload: title });
      },
      setSubtitle: (subtitle: string) => {
        dispatch({ type: "set_subtitle", payload: subtitle });
      },
      setHeaderContent: (headerContent: string) => {
        dispatch({ type: "set_content", payload: headerContent });
      },
    }),
    [state]
  );

  return (
    <WorkflowLayoutContext.Provider value={providerValue}>
      {children}
    </WorkflowLayoutContext.Provider>
  );
};

const useWorkflowLayoutContext = () => {
  const location = useLocation();
  const context = React.useContext<WorkflowLayoutContextProps>(WorkflowLayoutContext);

  if (!context) {
    throw new Error("useWorkflowLayoutContext was invoked outside of a valid context");
  }

  // Reset state on route change
  React.useEffect(() => {
    context.setTitle(null);
    context.setSubtitle(null);
    context.setHeaderContent(null);
  }, [location.pathname]);

  return context;
};

export { WorkflowLayoutContextProvider, useWorkflowLayoutContext };
