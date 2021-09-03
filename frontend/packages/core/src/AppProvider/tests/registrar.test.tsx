import React from "react";

import { registeredWorkflows, workflowRoutes } from "../registrar";
import type { DefaultWorkflowConfig, GatewayConfig, GatewayRoute, RouteConfigs } from "../types";

const WORKFLOW_ID = "@clutch-sh/ec2";
const WORKFLOW_CONFIG = {
  developer: {
    name: "Lyft",
    contactUrl: "mailto:hello@example.com",
  },
  path: "ec2",
  group: "AWS",
  displayName: "EC2",
  routes: {
    terminateInstance: {
      path: "instance/terminate",
      displayName: "Terminate Instance",
      description: "Terminate an EC2 instance.",
      requiredConfigProps: ["resolverType"],
      component: () => <></>,
    },
    rebootInstance: {
      path: "instance/reboot",
      displayName: "Reboot Instance",
      description: "Reboot an EC2 Instance",
      requiredConfigProps: ["resolverType"],
      component: () => <></>,
    },
    resizeAutoscalingGroup: {
      path: "asg/resize",
      displayName: "Resize Autoscaling Group",
      description: "Resize an autoscaling group.",
      requiredConfigProps: ["resolverType"],
      component: () => <></>,
    },
  },
} as DefaultWorkflowConfig;
const ROUTES = {
  terminateInstance: {
    component: () => <></>,
    path: "instance/term",
    description: "Custom description.",
    trending: true,
    componentProps: {
      resolverType: "clutch.aws.ec2.v1.Instance",
      notes: [
        {
          severity: "info",
          text: "Note: the instance may take several minutes to shut down.",
        },
      ],
    },
  },
} as RouteConfigs;

const GATEWAY_CONFIG = {
  [WORKFLOW_ID]: ROUTES,
} as GatewayConfig;

describe("workflowRoutes", () => {
  let warn;
  let gatewayRoutes;
  beforeEach(() => {
    gatewayRoutes = JSON.parse(JSON.stringify(ROUTES));
    warn = jest.spyOn(console, "warn").mockImplementation(() => {});
  });

  afterEach(() => {
    warn.mockReset();
  });

  afterAll(() => {
    warn.mockRestore();
  });

  it("handles empty gateway route configuration", () => {
    const routes = workflowRoutes(WORKFLOW_ID, WORKFLOW_CONFIG, {});
    expect(routes).toEqual([]);
  });

  it("warns on invalid gateway routes", () => {
    // eslint-disable-next-line
    gatewayRoutes["fakeRoute"] = {} as GatewayRoute;
    workflowRoutes(WORKFLOW_ID, WORKFLOW_CONFIG, gatewayRoutes);
    expect(warn).toHaveBeenCalledWith(
      "[@clutch-sh/ec2][fakeRoute] Invalid gateway config: route with specified name does not exist"
    );
  });

  it("removes invalid gateway routes", () => {
    // eslint-disable-next-line
    gatewayRoutes["fakeRoute"] = {} as GatewayRoute;
    const routes = workflowRoutes(WORKFLOW_ID, WORKFLOW_CONFIG, gatewayRoutes);
    expect(routes).toHaveLength(1);
  });

  it("warns on invalid gateway route configurations", () => {
    delete gatewayRoutes.terminateInstance.componentProps;
    workflowRoutes(WORKFLOW_ID, WORKFLOW_CONFIG, gatewayRoutes);
    expect(warn).toHaveBeenCalledWith(
      "[@clutch-sh/ec2][terminateInstance] Invalid gateway config: route is missing required props"
    );
  });

  it("removes routes with invalid gateway configurations", () => {
    delete gatewayRoutes.terminateInstance.componentProps;
    const routes = workflowRoutes(WORKFLOW_ID, WORKFLOW_CONFIG, gatewayRoutes);
    expect(routes).toHaveLength(0);
  });

  it("removes empty routes", () => {
    // @ts-ignore
    const routes = workflowRoutes(WORKFLOW_ID, WORKFLOW_CONFIG, { terminateInstance: {} });
    expect(routes).toEqual([]);
  });

  it("returns valid routes", () => {
    const routes = workflowRoutes(WORKFLOW_ID, WORKFLOW_CONFIG, gatewayRoutes);
    expect(routes).toHaveLength(1);
  });
});

describe("registeredWorkflows", () => {
  let warn;
  let gatewayConfig;
  const availableWorkflows = { [WORKFLOW_ID]: () => WORKFLOW_CONFIG };
  beforeEach(() => {
    gatewayConfig = JSON.parse(JSON.stringify(GATEWAY_CONFIG));
    warn = jest.spyOn(console, "warn").mockImplementation(() => {});
  });

  afterEach(() => {
    warn.mockClear();
  });

  afterAll(() => {
    warn.mockRestore();
  });

  it("handles empty workflows", () => {
    return registeredWorkflows().then(workflows => {
      expect(workflows).toHaveLength(0);
    });
  });

  it("handles empty gateway configuration", () => {
    return registeredWorkflows(undefined).then(workflows => {
      expect(workflows).toHaveLength(0);
    });
  });

  it("allows workflow overrides from gateway", () => {
    gatewayConfig[WORKFLOW_ID].overrides = { displayName: "Amazon Web Services" };
    return registeredWorkflows(availableWorkflows, gatewayConfig).then(workflows => {
      expect(workflows[0].displayName).toBe("Amazon Web Services");
    });
  });

  it("removes workflow overrides from gateway config", () => {
    gatewayConfig[WORKFLOW_ID].overrides = { displayName: "Amazon Web Services" };
    return registeredWorkflows(availableWorkflows, gatewayConfig).then(() => {
      expect(warn).not.toHaveBeenCalled();
    });
  });

  it("warns on any error parsing gateway route configs", () => {
    return registeredWorkflows(availableWorkflows, { [WORKFLOW_ID]: null }).then(() => {
      expect(warn).toHaveBeenCalledWith("[@clutch-sh/ec2] Not registered: invalid config");
    });
  });

  it("removes gateway route configs on any error", () => {
    return registeredWorkflows(availableWorkflows, { [WORKFLOW_ID]: null }).then(workflows => {
      expect(workflows).toHaveLength(0);
    });
  });

  it("warns if gateway workflow has zero valid routes", () => {
    delete gatewayConfig[WORKFLOW_ID].terminateInstance;
    return registeredWorkflows(availableWorkflows, gatewayConfig).then(() => {
      expect(warn).toHaveBeenCalledWith("[@clutch-sh/ec2] Not registered: zero routes found");
    });
  });

  it("removes gateway workflows with zero valid routes", () => {
    delete gatewayConfig[WORKFLOW_ID].terminateInstance;
    return registeredWorkflows(availableWorkflows, gatewayConfig).then(workflows => {
      expect(workflows).toHaveLength(0);
    });
  });

  it("returns gateway configured workflow routes", () => {
    return registeredWorkflows(availableWorkflows, gatewayConfig).then(workflows => {
      expect(workflows[0].routes[0].path).toBe("instance/term");
    });
  });

  it("applies filters to each workflow", () => {
    const filter = jest.fn().mockImplementation(() => new Promise(resolve => resolve([])));
    filter.mockReturnValue(new Promise(resolve => resolve([])));
    return registeredWorkflows(availableWorkflows, gatewayConfig, [filter]).then(workflows => {
      expect(filter).toHaveBeenLastCalledWith([
        {
          developer: { contactUrl: "mailto:hello@example.com", name: "Lyft" },
          displayName: "EC2",
          group: "AWS",
          path: "ec2",
          routes: [
            expect.objectContaining({
              componentProps: {
                notes: [
                  {
                    severity: "info",
                    text: "Note: the instance may take several minutes to shut down.",
                  },
                ],
                resolverType: "clutch.aws.ec2.v1.Instance",
              },
              description: "Custom description.",
              displayName: "Terminate Instance",
              path: "instance/term",
              requiredConfigProps: ["resolverType"],
              trending: true,
            }),
          ],
        },
      ]);
      expect(workflows).toHaveLength(0);
    });
  });
});
