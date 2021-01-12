import type { ConfiguredRoute, Workflow, WorkflowConfiguration } from "./workflow";
import type { UserConfiguration } from ".";

const workflowRoutes = (
  workflowId: string,
  workflow: WorkflowConfiguration,
  configuration: UserConfiguration
): ConfiguredRoute[] => {
  const workflowConfig = configuration?.[workflowId] || {};
  const allRoutes = Object.keys(workflowConfig).map(key => {
    return {
      ...workflow.routes[key],
      ...workflowConfig[key],
    };
  });

  const validRoutes = allRoutes.filter(route => {
    const requiredRouteProps = route?.requiredConfigProps || [];
    const missingProps = requiredRouteProps.filter((prop: string) => {
      return route.componentProps?.[prop] === undefined;
    });

    const isValidRoute = missingProps.length === 0;
    if (!isValidRoute) {
      /* eslint-disable-next-line no-console */
      console.warn(
        `[${workflowId}][${route.path}] Not registered: Invalid config - missing required component props ${missingProps}`
      );
    }
    return isValidRoute;
  });

  return validRoutes;
};

const registeredWorkflows = (
  workflows: { [key: string]: () => WorkflowConfiguration },
  configuration: UserConfiguration
): Workflow[] => {
  return Object.keys(workflows || [])
    .map((workflowId: string) => {
      const workflow = workflows[workflowId]();
      try {
        return { ...workflow, routes: workflowRoutes(workflowId, workflow, configuration) };
      } catch {
        // n.b. if the routes aren't configured properly we
        // drop the workflow
        /* eslint-disable-next-line no-console */
        console.warn(
          `Skipping registration of ${workflowId || "unknown"} workflow due to invalid config`
        );
        return null;
      }
    })
    .filter(workflow => workflow !== null);
};

export default registeredWorkflows;
