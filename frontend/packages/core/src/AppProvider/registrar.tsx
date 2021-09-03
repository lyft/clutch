import _ from "lodash";

import type {
  DefaultWorkflowConfig,
  GatewayConfig,
  GatewayRoute,
  RouteConfigs,
  Workflow,
  WorkflowFilters,
  Workflows,
} from "./types";

/** Warn end users a workflow or route is not being registered with a reason. */
const warnUnregistered = (
  workflowId: string,
  reason: string,
  message: string,
  routeName?: string
) => {
  /* eslint-disable-next-line no-console */
  console.warn(`[${workflowId}]${routeName ? `[${routeName}]` : ""} ${reason}: ${message}`);
};

/**
 * Determine if a route is valid by checking:
 *   * does the gateway route config have the required route props?
 */
const isValidRoute = (route: GatewayRoute): boolean => {
  const requiredRouteProps = route?.requiredConfigProps || [];
  const missingProps = requiredRouteProps.filter((prop: string) => {
    return route?.componentProps?.[prop] === undefined;
  });

  return missingProps.length === 0;
};

/**
 * Determine all valid routes registered on the gateway.
 *
 * @param workflowId
 * @param workflow
 * @param configuration
 */
const workflowRoutes = (
  workflowId: string,
  workflow: DefaultWorkflowConfig,
  routes: RouteConfigs = {}
): GatewayRoute[] => {
  const validRoutes = Object.keys(routes).map(routeName => {
    // if workflow does not contain route with gateway config key ignore gateway route
    if (workflow.routes?.[routeName] === undefined) {
      warnUnregistered(
        workflowId,
        "Invalid gateway config",
        `route with specified name does not exist`,
        routeName
      );
      return null;
    }
    const routeConfig = { ...workflow.routes[routeName], ...routes[routeName] };
    if (!isValidRoute(routeConfig)) {
      warnUnregistered(
        workflowId,
        "Invalid gateway config",
        `route is missing required props`,
        routeName
      );
      return null;
    }
    return routeConfig;
  });
  // filter out routes that are empty
  _.remove(validRoutes, r => _.isEmpty(r));

  return validRoutes;
};

/**
 * Determine the workflows registered on the gateway and apply filters, if any.
 *
 * @param workflows A map of unique workflow keys (package names) to functions that return their configuration.
 * @param gatwayConfig Configuration specific to the Clutch gateway running.
 * @param filters A list of filters to apply to the gateway's registered workflows.
 */
const registeredWorkflows = async (
  workflows: Workflows = {},
  gatewayConfig: GatewayConfig = {},
  filters: WorkflowFilters = []
): Promise<Workflow[]> => {
  let wkflws = Object.keys(workflows)
    .map((workflowId: string) => {
      let routes = [];
      const defaultWorkflowConfig = workflows[workflowId]();
      const gatewayWorkflowConfig = gatewayConfig[workflowId];
      const gatewayWorkflowOverrides = gatewayWorkflowConfig?.overrides || {};
      // remove gateway workflow overrides before gathering routes as they are at the same level
      delete gatewayWorkflowConfig?.overrides;
      try {
        routes = workflowRoutes(
          workflowId,
          defaultWorkflowConfig,
          gatewayWorkflowConfig as RouteConfigs
        );
      } catch {
        // if the routes aren't configured properly we drop the workflow
        warnUnregistered(workflowId, "Not registered", "invalid config");
        return null;
      }

      // if workflow has no routes drop the workflow
      if (_.isEmpty(routes)) {
        warnUnregistered(workflowId, "Not registered", "zero routes found");
        return null;
      }

      return { ...defaultWorkflowConfig, ...gatewayWorkflowOverrides, routes };
    })
    .filter(workflow => workflow !== null);
  await filters.forEach(async f => {
    wkflws = await f(wkflws);
  });
  return wkflws;
};

export { isValidRoute, registeredWorkflows, warnUnregistered, workflowRoutes };
