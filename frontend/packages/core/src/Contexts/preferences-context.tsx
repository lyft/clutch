import React from "react";
import _ from "lodash";

export type ActionType = "SetPref" | "RemovePref" | "SetLocalPref" | "RemoveLocalPref";
const STORAGE_KEY = "userPreferences";
type State = { key: string; value?: unknown };
type Action = { type: ActionType; payload: State };
type Dispatch = (action: Action) => void;
type UserPreferencesProviderProps = { children: React.ReactNode };
const DEFAULT_PREFERENCES: State = {
  timeFormat: "UTC",
  theme: "",
} as any;
interface ContextProps {
  preferences: State;
  dispatch: Dispatch;
}

type ContextType = ContextProps | undefined;

export interface PreferencesContextProps {
  preferences: {
    [key: string]: unknown;
  };
  getPref: (key: string) => unknown;
  setPref: (key: string, value: unknown) => void;
}

const preferencesReducer = (preferences: State, action: Action): State => {
  switch (action.type) {
    case "SetPref": {
      const updatedPref = { ...preferences, [action.payload.key]: action.payload.value };
      try {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(updatedPref));
      } catch {}
      return updatedPref;
    }
    case "RemovePref": {
      const updatedPref = { ...preferences };
      delete updatedPref[action.payload.key];
      try {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(updatedPref));
      } catch {}
      return updatedPref;
    }
    case "SetLocalPref": {
      return { ...preferences, [action.payload.key]: action.payload.value };
    }
    case "RemoveLocalPref": {
      const updatedPref = { ...preferences };
      delete updatedPref[action.payload.key];
      return updatedPref;
    }
    default: {
      throw new Error(`Unhandled action type: ${action.type}`);
    }
  }
};

const UserPreferencesContext = React.createContext<ContextType>(undefined);

const UserPreferencesProvider = ({ children }: UserPreferencesProviderProps) => {
  // Load preferences as default value and if none then default to value
  let pref = DEFAULT_PREFERENCES;
  try {
    pref = JSON.parse(localStorage.getItem(STORAGE_KEY) || "");
    // If there are any missing default preferences, add them
    Object.keys(DEFAULT_PREFERENCES).forEach(key => {
      if (_.isEmpty(pref[key]) || !pref[key]) {
        pref[key] = DEFAULT_PREFERENCES[key];
      }
    });
  } catch {
    localStorage.removeItem(STORAGE_KEY);
  }

  localStorage.setItem(STORAGE_KEY, JSON.stringify(pref));
  const [state, dispatch] = React.useReducer(preferencesReducer, pref);

  const value = React.useMemo(() => ({ preferences: state, dispatch }), [state, dispatch]);

  return (
    <UserPreferencesContext.Provider value={value}>{children}</UserPreferencesContext.Provider>
  );
};

const useUserPreferences = () => {
  const context = React.useContext(UserPreferencesContext);
  if (!context) {
    throw new Error("useUserPreferences was invoked outside of a valid context");
  }
  return context;
};

export { useUserPreferences, UserPreferencesProvider };
