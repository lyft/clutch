/** A route belonging to a specific workflow. */
interface Route {
  /** The component to render for this route. */
  component: React.FC<any>;
  /** A description of the route */
  description: string;
  /** Name of the route used for display purposes to end users. */
  displayName?: string;
  /** Path of the route. */
  path: string;
  /** Properties required by the Component that are set only via the config. */
  requiredConfigProps?: string[];
  /** Is the workflow discoverable via search and drawer navigation. This defaults to false. */
  hideNav?: boolean;
  /**
   * The feature flag name used to determine if the route should be registered.
   *
   * If this is not set the route will always be registered.
   */
  featureFlag?: string;
}

/** A worfklow route configuration for a gateway. */
export interface GatewayRoute extends Route {
  componentProps?: object;
  trending?: boolean;
}

/** A map of unique route keys to gateway configurations. */
export interface RouteConfigs {
  [routeName: string]: GatewayRoute;
}

/** Overrides to apply to a workflow within a gateway. */
export interface WorkflowOverrides {
  overrides?: Partial<WorkflowBase>;
}

/** Configuration specific to the gateway running. */
export interface GatewayConfig {
  /** Unique package name to collection of routes */
  [packageName: string]: RouteConfigs | WorkflowOverrides;
}

export interface BaseWorkflowProps {
  heading: string;
}

interface Developer {
  contactUrl: string;
  name: string;
}

interface WorkflowBase {
  developer: Developer;
  displayName: string;
  group: string;
  path: string;
  routes: unknown;
}

/** A workflow running within a gateway. */
export interface Workflow extends WorkflowBase {
  routes: GatewayRoute[];
}

/** Default configuration for a workflow. */
export interface DefaultWorkflowConfig extends WorkflowBase {
  routes: {
    [key: string]: Route;
  };
}

/**
 * A map of unique workflow keys (package names) to functions that return their configuration.
 */
export interface Workflows {
  [packageName: string]: () => DefaultWorkflowConfig;
}

/** A collection of workflow filters. */
export type WorkflowFilters = { (workflows: Workflow[]): Promise<Workflow[]> }[];
