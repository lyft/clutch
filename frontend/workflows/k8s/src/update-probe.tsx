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

// Examples of valid Seconds: 0.1, 100s
const SECONDS_REGEX = /^([-+]?[0-9]*[sS]*)$/;

const DeploymentIdentifier: React.FC<ResolverChild> = ({ resolverType }) => {
  const { onSubmit } = useWizardContext();
  const deploymentData = useDataLayout("deploymentData");
  const inputData = useDataLayout("inputData");

  const onResolve = ({ results, input }) => {
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

  const [containerName, setContainerName] = React.useState(
    containers && containers.length > 0 ? containers[0].name : ""
  );
  const [probeType, setProbeType] = React.useState("livenessProbe");

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

  const listedProbes = [];

  if (Object.keys(currentDeployment.livenessProbe).length !== 0) listedProbes.push("livenessProbe");
  if (Object.keys(currentDeployment.readinessProbe).length !== 0)
    listedProbes.push("readinessProbe");

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
                options={containers.map(container => ({ label: container.name }))}
              />
            ),
          },
          {
            name: "Probe Type",
            value: (
              <Select
                label="Probe Type"
                name="probeType"
                onChange={value => {
                  setProbeType(value);
                }}
                options={listedProbes.map(probe => ({ label: probe }))}
              />
            ),
          },
          {
            name: "Initial Delay Seconds",
            value: currentDeployment[probeType].initialDelaySeconds,
            textFieldLabels: {
              disabledField: "Current Value",
              updatedField: "New Value",
            },
            input: {
              type: "string",
              key: `${containerBase}.${probeType}.initialDelaySeconds`,
              validation: string().matches(SECONDS_REGEX),
            },
          },
          {
            name: "Period Seconds",
            value: currentDeployment[probeType].periodSeconds,
            textFieldLabels: {
              disabledField: "Current Value",
              updatedField: "New Value",
            },
            input: {
              type: "string",
              key: `${containerBase}.${probeType}.periodSeconds`,
              validation: string().matches(SECONDS_REGEX),
            },
          },
          {
            name: "Timeout Seconds",
            value: currentDeployment[probeType].timeoutSeconds,
            textFieldLabels: {
              disabledField: "Current Value",
              updatedField: "New Value",
            },
            input: {
              type: "string",
              key: `${containerBase}.${probeType}.timeoutSeconds`,
              validation: string().matches(SECONDS_REGEX),
            },
          },
          {
            name: "Success Threshold",
            value: currentDeployment[probeType].successThreshold,
            textFieldLabels: {
              disabledField: "Current Value",
              updatedField: "New Value",
            },
            input: {
              type: "string",
              key: `${containerBase}.${probeType}.successThreshold`,
              validation: string().matches(SECONDS_REGEX),
            },
          },
          {
            name: "Failure Threshold",
            value: currentDeployment[probeType].failureThreshold,
            textFieldLabels: {
              disabledField: "Current Value",
              updatedField: "New Value",
            },
            input: {
              type: "string",
              key: `${containerBase}.${probeType}.failureThreshold`,
              validation: string().matches(SECONDS_REGEX),
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

const Confirm: React.FC<ConfirmChild> = () => {
  const deployment = useDataLayout("deploymentData").displayValue() as IClutch.k8s.v1.Deployment;
  const updateData = useDataLayout("updateData");
  const currentDeploymentData = useDataLayout(
    "currentDeploymentData"
  ).displayValue() as IClutch.k8s.v1.Deployment;

  const updateRows: any[] = [];

  let updatedContainer = false;
  let probeType = "livenessProbe";
  deployment.deploymentSpec.template.spec.containers.forEach(container => {
    if (!updatedContainer) {
      updateRows.push({ name: "Container Name", value: container.name });
      updatedContainer = true;
    }
    ["livenessProbe", "readinessProbe"].forEach(pType => {
      Object.keys(container[pType]).forEach(probeAttribute => {
        const typeValue = typeof container[pType][probeAttribute];
        if (typeValue !== "object") {
          const newValue = container[pType][probeAttribute];
          const oldValue = findContainer({
            deploymentSpec: currentDeploymentData.deploymentSpec,
            containerName: container.name,
          }).livenessProbe[probeAttribute];
          if (newValue !== oldValue) {
            updateRows.push({
              name: `Old ${probeAttribute}`,
              value: oldValue,
            });
            updateRows.push({
              name: `New ${probeAttribute}`,
              value: newValue,
            });
            probeType = pType;
          }
        }
      });
    });
    updateRows.push({ name: "Probe type", value: probeType });
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

const UpdateLiveness: React.FC<WorkflowProps> = ({ heading, resolverType }) => {
  interface BodyRequest {
    clientset: any;
    cluster: string;
    namespace: string;
    name: string;
    fields: {
      containerProbes: {
        containerName: string;
        [key: string]: any; // Permitir otros atributos dinámicos
      }[];
    };
  }

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
        const bodyRequest: BodyRequest = {
          clientset,
          cluster: deploymentData.cluster,
          namespace: deploymentData.namespace,
          name: deploymentData.name,
          fields: {
            containerProbes: [
              {
                containerName: deploymentData.containerName,
              },
            ],
          },
        };
        ["livenessProbe", "readinessProbe"].forEach(probeType => {
          if (Object.keys(container[probeType]).length !== 0) {
            const {
              timeoutSeconds,
              initialDelaySeconds,
              periodSeconds,
              successThreshold,
              failureThreshold,
            } = container[probeType];

            bodyRequest.fields.containerProbes[0][probeType] = {
              timeoutSeconds,
              initialDelaySeconds,
              periodSeconds,
              successThreshold,
              failureThreshold,
            };
          }
        });
        return client.post(
          "/v1/k8s/updateDeployment",
          bodyRequest as IClutch.k8s.v1.UpdateDeploymentRequest
        );
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

export default UpdateLiveness;
