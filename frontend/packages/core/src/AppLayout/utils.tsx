import * as _ from "lodash";

import type { WorkflowIcon } from "../AppProvider";
import type { Route, Workflow } from "../AppProvider/workflow";

interface GroupedRoutes {
  [category: string]: {
    icon?: WorkflowIcon;
    workflows: {
      displayName: string;
      path: string;
      trending: boolean;
    }[];
  };
}

interface TrendingWorkflow {
  displayName: string;
  group: string;
  description: string;
  path: string;
  icon: string;
}

const getDisplayName = (workflow: Workflow, route: Route, delimiter: string = ":"): string => {
  let { displayName } = workflow;
  if (route.displayName) {
    displayName = `${
      displayName.toLowerCase() !== workflow.group.toLowerCase()
        ? `${displayName}${delimiter} `
        : ""
    }${route.displayName}`;
  }

  return displayName;
};

const workflowsByTrending = (workflows: Workflow[]): TrendingWorkflow[] => {
  const trending = [];
  const trendingIcons = {};

  workflows.forEach(workflow => {
    if (workflow?.icon?.path && !trendingIcons[workflow.group]) {
      trendingIcons[workflow.group] = workflow.icon.path;
    }
  });

  workflows.forEach(workflow => {
    workflow.routes.forEach(route => {
      if (route.trending) {
        trending.push({
          displayName: getDisplayName(workflow, route),
          group: workflow.group,
          description: route.description,
          path: `${workflow.path}/${route.path}`,
          icon: trendingIcons[workflow.group] ?? "",
        });
      }
    });
  });

  return trending;
};

const routesByGrouping = (workflows: Workflow[]): GroupedRoutes => {
  const routes = {};
  workflows.forEach(workflow => {
    const category = workflow.group;

    routes[category] ??= {
      workflows: [],
      icon: workflow.icon,
    };

    routes[category].icon.path = routes[category].icon.path || workflow.icon?.path;

    routes[category].workflows = [
      ...routes[category].workflows,
      ...workflow.routes.map(route => {
        return {
          displayName: getDisplayName(workflow, route, " -"),
          path: `${workflow.path}/${route.path}`,
          trending: route.trending || false,
        };
      }),
    ];
  });
  return routes;
};

/**
 * Will break a path down and iterate through given workflows to see if there is a matching path.
 */
const workflowByRoute = (workflows: Workflow[], route: string): Workflow => {
  const [baseRoute, ...subRoutes] = route.split("/").filter(Boolean);
  const subRoute = subRoutes.join("/");
  let returnFlow = null;

  const filtered = workflows.filter((workflow: Workflow) => workflow.path === baseRoute);

  filtered.some((workflow: Workflow) => {
    return workflow.routes.some((wroute: any) => {
      if (wroute.path === subRoute) {
        returnFlow = workflow;
      }
      return returnFlow !== null;
    });
  });

  return returnFlow;
};

const sortedGroupings = (workflows: Workflow[]): string[] => {
  return Object.keys(routesByGrouping(workflows)).sort();
};

export interface SearchIndex {
  category: string;
  label: string;
  path: string;
}

const searchIndexes = (workflows: Workflow[]): SearchIndex[] => {
  let indexOptions = [];
  workflows.forEach(workflow => {
    const category = workflow.group;
    indexOptions = [
      ...indexOptions,
      ...workflow.routes.map(route => {
        const label = route.displayName
          ? `${workflow.displayName} - ${route.displayName}`
          : workflow.displayName;
        return {
          category,
          label,
          path: `${workflow.path}/${route.path}`,
        };
      }),
    ];
  });
  return indexOptions;
};

/** Filter out all of the workflows that are configured to be `hideNav: true`.
 * This prevents the workflows from being discoverable by the user from the UI where used.
 * Some example usages are in the search and drawer navigation components.
 *
 * The routes for all configured workflows will still be reachable
 * by manually providing the full path in the URI.
 */
const filterHiddenRoutes = (workflows: Workflow[]): Workflow[] =>
  _.cloneDeep(workflows).filter(workflow => {
    const publicRoutes = workflow.routes.filter(route => {
      return !(route?.hideNav !== undefined ? route.hideNav : false);
    });
    workflow.routes = publicRoutes; /* eslint-disable-line no-param-reassign */
    return publicRoutes.length !== 0;
  });

export {
  filterHiddenRoutes,
  routesByGrouping,
  searchIndexes,
  sortedGroupings,
  workflowByRoute,
  workflowsByTrending,
};
