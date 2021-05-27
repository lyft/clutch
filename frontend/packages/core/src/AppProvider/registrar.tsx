import type { UserConfiguration } from ".";
import type { ConfiguredRoute, Workflow, WorkflowConfiguration } from "./workflow";

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

/**
 * Determine all user registered workflows on the application and apply filters, if any.
 * @param workflows a map of workflow keys to functions that return their configuration.
 * @param configuration the user configuration, usually read in from the clutch.config.js file.
 * @param filters a list of filters to apply to the user registered workflows.
 * @returns
 */
const registeredWorkflows = async (
  workflows: { [key: string]: () => WorkflowConfiguration },
  configuration: UserConfiguration,
  filters: { (workflows: Workflow[]): Promise<Workflow[]> }[] = []
): Promise<Workflow[]> => {
  let validWorkflows = Object.keys(workflows || [])
    .map((workflowId: string) => {
      const workflow = workflows[workflowId]();
      try {
        return { ...workflow, routes: workflowRoutes(workflowId, workflow, configuration) };
      } catch {
        // n.b. if the routes aren't configured properly we drop the workflow
        /* eslint-disable-next-line no-console */
        console.warn(
          `Skipping registration of ${workflowId || "unknown"} workflow due to invalid config`
        );
        return null;
      }
    })
    .filter(workflow => workflow !== null);
  filters.forEach(f => {
    f(validWorkflows).then(w => {
      validWorkflows = w;
    });
  });
  return validWorkflows;
};

export default registeredWorkflows;
