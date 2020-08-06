import EchoWorkflow from "./echo";

const register = () => {
  return {
    developer: {
      name: "Lyft",
      contactUrl: "https://app.slack.com/client/T029A67TC/CQLRKE0ER",
    },
    path: "examples",
    group: "Examples",
    displayName: "Echo",
    routes: {
      echo: {
        path: "echo",
        component: EchoWorkflow,
        displayName: "Demo Workflow",
        description: "Workflow for demonstration purposes.",
      },
    },
  };
};

export default register;
