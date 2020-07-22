import type { Workflow } from "../AppProvider/workflow";

interface GroupedRoutes {
  [category: string]: {
    workflows: {
      displayName: string;
      path: string;
      trending: boolean;
    }[];
  };
}

const routesByGrouping = (workflows: Workflow[]): GroupedRoutes => {
  const routes = {};
  workflows.forEach(workflow => {
    const category = workflow.group;
    if (routes[category] === undefined) {
      routes[category] = { workflows: [] };
    }

    routes[category].workflows = [
      ...routes[category].workflows,
      ...workflow.routes.map(route => {
        const displayName = route.displayName
          ? `${workflow.displayName} - ${route.displayName}`
          : workflow.displayName;
        return {
          displayName,
          path: `${workflow.path}/${route.path}`,
          trending: route.trending,
        };
      }),
    ];
  });
  return routes;
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

export { routesByGrouping, searchIndexes, sortedGroupings };
