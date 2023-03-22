import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import {
  Button,
  ButtonGroup,
  client,
  Confirmation,
  MetadataTable,
  Resolver,
  Select,
  useWizardContext,
} from "@clutch-sh/core";
import { useDataLayout } from "@clutch-sh/data-layout";
import type { WizardChild } from "@clutch-sh/wizard";
import { Wizard, WizardStep } from "@clutch-sh/wizard";
import { string } from "yup";

import type { ConfirmChild, ResolverChild, WorkflowProps } from ".";

// Examples of valid quantities: 0.1, 100m, 128974848, 129e6, 129M, 128974848000m, 123Mi
const QUANTITY_REGEX = /^([+-]?[0-9.]+)([eEinumkKMGTP]*[-+]?[0-9]*)$/;

const DeploymentIdentifier: React.FC<ResolverChild> = ({ resolverType }) => {
  const { onSubmit } = useWizardContext();
  const deploymentData = useDataLayout("deploymentData");
  const inputData = useDataLayout("inputData");

  const onResolve = ({ results, input }) => {
    // Decide how to process results.
    deploymentData.assign(results[0]);
    inputData.assign(input);
    onSubmit();
  };

  return <Resolver type={resolverType} searchLimit={1} onResolve={onResolve} />;
};

function findContainer(args: {
  deploymentSpec: IClutch.k8s.v1.Deployment.IDeploymentSpec;
  containerName: string;
}): IClutch.k8s.v1.Deployment.DeploymentSpec.PodTemplateSpec.PodSpec.IContainer {
  return args.deploymentSpec.template.spec.containers.find(
    container => container.name === args.containerName
  );
}

const DeploymentDetails: React.FC<WizardChild> = () => {
  const { onSubmit, onBack } = useWizardContext();
  const deploymentData = useDataLayout("deploymentData");
  const deployment = deploymentData.displayValue() as IClutch.k8s.v1.Deployment;
  const update = (key: string, value: boolean) => {
    deploymentData.updateData(key, value);
  };

  const currentDeploymentData = useDataLayout("currentDeploymentData");

  const { containers } = deployment.deploymentSpec.template.spec;

  const [containerName, setContainerName] = React.useState(containers[0].name);

  const [containerIndex, setContainerIndex] = React.useState(0);

  React.useEffect(() => {
    // save the original values of deployment spec
    if (deployment) {
      currentDeploymentData.assign(deployment);
    }
  }, []);

  const currentDeployment = findContainer({
    deploymentSpec: deployment.deploymentSpec,
    containerName,
  });

  const containerBase = `deploymentSpec.template.spec.containers[${containerIndex}]`;

  return (
    <WizardStep error={deploymentData.error} isLoading={deploymentData.isLoading}>
      <strong>Deployment Details</strong>
      <MetadataTable
        onUpdate={update}
        data={[
          { name: "Name", value: deployment.name },
          { name: "Namespace", value: deployment.namespace },
          {
            name: "Container Name",
            value: (
              <Select
                label="Container Name"
                name="containerName"
                onChange={value => {
                  setContainerName(value);
                  setContainerIndex(containers.findIndex(container => container.name === value));
                  deploymentData.updateData("containerName", value);
                }}
                options={containers.map(container => {
                  return { label: container.name };
                })}
              />
            ),
          },
          {
            name: "CPU Limit",
            value: currentDeployment.resources.limits.cpu,
            textFieldLabels: {
              disabledField: "Current Limit",
              updatedField: "New limit",
            },
            input: {
              type: "string",
              key: `${containerBase}.resources.limits.cpu`,
              validation: string().matches(QUANTITY_REGEX),
            },
          },
          {
            name: "CPU Request",
            value: currentDeployment.resources.requests.cpu,
            textFieldLabels: {
              disabledField: "Current Request",
              updatedField: "New Request",
            },
            input: {
              type: "string",
              key: `${containerBase}.resources.requests.cpu`,
              validation: string().matches(QUANTITY_REGEX),
            },
          },
          {
            name: "Memory Limit",
            value: currentDeployment.resources.limits.memory,
            textFieldLabels: {
              disabledField: "Current Limit",
              updatedField: "New limit",
            },
            input: {
              type: "string",
              key: `${containerBase}.resources.limits.memory`,
              validation: string().matches(QUANTITY_REGEX),
            },
          },
          {
            name: "Memory Request",
            value: currentDeployment.resources.requests.memory,
            textFieldLabels: {
              disabledField: "Current Request",
              updatedField: "New Request",
            },
            input: {
              type: "string",
              key: `${containerBase}.resources.requests.memory`,
              validation: string().matches(QUANTITY_REGEX),
            },
          },
        ]}
      />
      <ButtonGroup>
        <Button text="Back" variant="neutral" onClick={() => onBack()} />
        <Button text="Update" variant="destructive" onClick={onSubmit} />
      </ButtonGroup>
    </WizardStep>
  );
};

function formatResourceString(resourceName: string, resourceRequirement: string): string {
  // Capitalize the first letter of resourceName
  const capitalizedResourceName = resourceName.charAt(0).toUpperCase() + resourceName.slice(1);

  // Capitalize and remove the s at the end of resourceRequirement
  const modifiedResourceRequirement =
    resourceRequirement.charAt(0).toUpperCase() + resourceRequirement.slice(1, -1);

  // Return the modified strings
  return `${capitalizedResourceName} ${modifiedResourceRequirement}`;
}

const Confirm: React.FC<ConfirmChild> = () => {
  const deployment = useDataLayout("deploymentData").displayValue() as IClutch.k8s.v1.Deployment;
  const updateData = useDataLayout("updateData");
  const currentDeploymentData = useDataLayout(
    "currentDeploymentData"
  ).displayValue() as IClutch.k8s.v1.Deployment;

  const updateRows: any[] = [];

  let updatedContainer = false;
  deployment.deploymentSpec.template.spec.containers.forEach(container => {
    Object.keys(container.resources).forEach(resourceRequirement => {
      Object.keys(container.resources[resourceRequirement]).forEach(resourceName => {
        const newValue = container.resources[resourceRequirement][resourceName];
        const oldValue = findContainer({
          deploymentSpec: currentDeploymentData.deploymentSpec,
          containerName: container.name,
        }).resources[resourceRequirement][resourceName];
        if (newValue !== oldValue) {
          if (!updatedContainer) {
            updateRows.push({ name: "Container Name", value: container.name });
            updatedContainer = true;
          }
          updateRows.push({
            name: `Old ${formatResourceString(resourceName, resourceRequirement)}`,
            value: oldValue,
          });
          updateRows.push({
            name: `New ${formatResourceString(resourceName, resourceRequirement)}`,
            value: newValue,
          });
        }
      });
    });
  });

  return (
    <WizardStep error={updateData.error} isLoading={updateData.isLoading}>
      <Confirmation action="Update" />
      <MetadataTable
        data={[
          { name: "Name", value: deployment.name },
          { name: "Namespace", value: deployment.namespace },
          { name: "Cluster", value: deployment.cluster },
          ...updateRows,
        ]}
      />
    </WizardStep>
  );
};

const ScaleResources: React.FC<WorkflowProps> = ({ heading, resolverType }) => {
  const dataLayout = {
    inputData: {},
    deploymentData: {},
    currentDeploymentData: {},
    updateData: {
      deps: ["deploymentData", "inputData", "currentDeploymentData"],
      hydrator: (
        deploymentData: {
          cluster: string;
          containerName: string;
          deploymentSpec: IClutch.k8s.v1.Deployment.DeploymentSpec;
          name: string;
          namespace: string;
        },
        inputData: { clientset: string },
        currentDeploymentData: IClutch.k8s.v1.Deployment
      ) => {
        const clientset = inputData.clientset ?? "undefined";
        const container = findContainer({ ...deploymentData });
        const limits: { [key: string]: string } = {
          cpu: container.resources.limits.cpu,
          memory: container.resources.limits.memory,
        };
        const requests: { [key: string]: string } = {
          cpu: container.resources.requests.cpu,
          memory: container.resources.requests.memory,
        };
        return client.post("/v1/k8s/updateDeployment", {
          clientset,
          cluster: deploymentData.cluster,
          namespace: deploymentData.namespace,
          name: deploymentData.name,
          fields: {
            containerResources: [
              {
                containerName: deploymentData.containerName,
                resources: {
                  limits,
                  requests,
                },
              },
            ],
          },
        } as IClutch.k8s.v1.UpdateDeploymentRequest);
      },
    },
  };

  return (
    <Wizard dataLayout={dataLayout} heading={heading}>
      <DeploymentIdentifier name="Lookup" resolverType={resolverType} />
      <DeploymentDetails name="Modify" />
      <Confirm name="Confirmation" />
    </Wizard>
  );
};

export default ScaleResources;
