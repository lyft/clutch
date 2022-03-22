import type { Environments } from "./types";

const EnvironmentIcon = (environment: Environments) => {
  switch (environment.toLowerCase()) {
    case "setup":
      return "🔨";
    case "staging":
      return "🥚";
    case "canary":
      return "🐣";
    case "production":
      return "🐥";
    default:
      return "❓";
  }
};

export default EnvironmentIcon;
