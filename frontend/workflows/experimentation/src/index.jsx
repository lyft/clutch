import ListExperiments from "./list-experiments";

const register = function register() {
  return {
    developer: {
      name: "Lyft",
      contactUrl: "mailto:hello@clutch.sh",
    },
    path: "experimentation",
    group: "Experimentation",
    displayName: "Experimentation",
    routes: {
      listExperiments: {
        path: "list",
        displayName: "List Experiments",
        description: "List Experiments.",
        component: ListExperiments,
      },
    },
  };
};

export default register;
