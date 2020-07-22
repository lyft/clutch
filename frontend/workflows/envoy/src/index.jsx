import RemoteTriage from "./remote-triage";

const register = () => {
  return {
    developer: {
      name: "Lyft",
      contactUrl: "mailto:hello@clutch.sh",
    },
    path: "envoy",
    group: "Envoy",
    displayName: "Envoy",
    routes: {
      remoteTriage: {
        path: "triage",
        component: RemoteTriage,
        displayName: "Remote Triage",
        description: "Triage Envoy configurations.",
        trending: true,
      },
    },
  };
};

export default register;
