import _ from "lodash";

import type { UserConfiguration } from ".";
import type { ConfiguredRoute, Workflow, WorkflowConfiguration } from "./workflow";

const workflowRoutes = (
  workflowId: string,
  workflow: WorkflowConfiguration,
  configuration: UserConfiguration
): ConfiguredRoute[] => {
  const workflowConfig = configuration?.[workflowId] || {};
  const allRoutes = Object.keys(workflowConfig).map(key => {
    // if workflow contains an icon, return an empty object
    if (key === "icon") {
      return {} as ConfiguredRoute;
      // if workflow does not contain route with user-specified key return an empty object
    }
    if (workflow.routes[key] === undefined) {
      /* eslint-disable-next-line no-console */
      console.warn(
        `[${workflowId}][${key}] Not registered: Invalid config - route does not exist. Valid routes: ${Object.keys(
          workflow.routes
        )}`
      );
      return {} as ConfiguredRoute;
    }
    return {
      ...workflow.routes[key],
      ...workflowConfig[key],
    };
  });
  // filter out routes that are empty
  _.remove(allRoutes, r => !!_.isEmpty(r));

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
  // eslint-disable-next-line
  validRoutes.map(r => (r.path = r.path.replace(/^\/+/, "").replace(/\/+$/, "")));

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
  return new Promise(async resolve => {
    let validWorkflows = Object.keys(workflows || [])
      .map((workflowId: string) => {
        const workflow = workflows[workflowId]();
        const icon = configuration?.[workflowId]?.icon || { path: "" };
        try {
          return {
            ...workflow,
            icon,
            routes: workflowRoutes(workflowId, workflow, configuration),
          };
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
    try {
      await Promise.all(filters.map(f => f(validWorkflows).then(w => (validWorkflows = w))));
    } catch (e) {
      /* eslint-disable-next-line no-console */
      console.warn("Error applying filters to workflows", e);
    }
    resolve(validWorkflows);
  });
};

export { registeredWorkflows, workflowRoutes };
