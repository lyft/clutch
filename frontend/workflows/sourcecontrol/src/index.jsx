import CreateRepository from "./create-repository";

const register = () => {
  return {
    developer: {
      name: "Lyft",
      contactUrl: "mailto:hello@clutch.sh",
    },
    path: "SCM",
    group: "Source Control",
    displayName: "Source Control",
    routes: {
      createRepository: {
        path: "createRepository",
        component: CreateRepository,
        displayName: "Create Repository",
        description: "Create a new repository.",
      },
    },
  };
};

export default register;
