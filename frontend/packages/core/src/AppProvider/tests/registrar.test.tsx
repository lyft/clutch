import React from "react";

import { workflowRoutes } from "../registrar";

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
};
const USER_CONFIGURATION = {
  "@clutch-sh/ec2": {
    terminateInstance: {
      component: () => <></>,
      path: "instance/terminate",
      description: "Terminate an EC2 instance.",
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
  },
};

describe("workflowRoutes", () => {
  let warn;
  beforeAll(() => {
    warn = jest.spyOn(console, "warn").mockImplementation();
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  afterAll(() => {
    warn.mockRestore();
  });

  it("handles empty user configuration", () => {
    const routes = workflowRoutes(WORKFLOW_ID, WORKFLOW_CONFIG, {});
    expect(routes).toEqual([]);
  });

  it("handles non-existant workflow IDs", () => {
    const routes = workflowRoutes("test-workflow-id", WORKFLOW_CONFIG, USER_CONFIGURATION);
    expect(routes).toEqual([]);
  });

  it("warns when user-specified workflow route is missing", () => {
    const gatewayCfg = { ...USER_CONFIGURATION };
    // eslint-disable-next-line
    gatewayCfg["@clutch-sh/ec2"]["fakeRoute"] = {};
    workflowRoutes(WORKFLOW_ID, WORKFLOW_CONFIG, gatewayCfg);
    expect(warn).toHaveBeenCalledWith(
      "[@clutch-sh/ec2][fakeRoute] Not registered: Invalid config - route does not exist. Valid routes: terminateInstance,rebootInstance,resizeAutoscalingGroup"
    );
  });

  it("filters out empty routes", () => {
    const gatewayCfg = { ...USER_CONFIGURATION };
    // eslint-disable-next-line
    gatewayCfg["@clutch-sh/ec2"]["fakeRoute"] = {};
    const routes = workflowRoutes(WORKFLOW_ID, WORKFLOW_CONFIG, gatewayCfg);
    expect(routes).toHaveLength(1);
  });

  it("warns on missing required route props", () => {
    const gatewayCfg = { ...USER_CONFIGURATION };
    delete gatewayCfg["@clutch-sh/ec2"].terminateInstance.componentProps;
    workflowRoutes(WORKFLOW_ID, WORKFLOW_CONFIG, gatewayCfg);
    expect(warn).toHaveBeenCalledWith(
      "[@clutch-sh/ec2][instance/terminate] Not registered: Invalid config - missing required component props resolverType"
    );
  });
  it("filters out routes missing required props", () => {
    const gatewayCfg = { ...USER_CONFIGURATION };
    delete gatewayCfg["@clutch-sh/ec2"].terminateInstance.componentProps;
    const routes = workflowRoutes(WORKFLOW_ID, WORKFLOW_CONFIG, gatewayCfg);
    expect(routes).toHaveLength(0);
  });

  it("returns valid routes", () => {
    const routes = workflowRoutes(WORKFLOW_ID, WORKFLOW_CONFIG, USER_CONFIGURATION);
    expect(routes).toEqual([]);
  });
});
