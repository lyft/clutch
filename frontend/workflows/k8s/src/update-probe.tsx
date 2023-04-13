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
const SECONDS_REGEX = /^([-+]?[0-9]*[sS]*)$/;
const PORT_REGEX = /^([1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$/;

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

function formatResourceData(Name: string, Value: string | number | string[], Path: string): object {
  return {
    name: Name,
    value: Value,
    textFieldLabels: {
      disabledField: "Current Value",
      updatedField: "New Value",
    },
    input: {
      type: "string",
      key: Path,
      validation: string().matches(SECONDS_REGEX),
    },
  };
}

function getHandlerObjects(
  currentDeployment: IClutch.k8s.v1.Deployment.DeploymentSpec.PodTemplateSpec.PodSpec.IContainer,
  containerBase: String
) {
  const handlerObject = [];
  if (
    currentDeployment.livenessProbe.handler.exec &&
    currentDeployment.livenessProbe.handler.exec.command
  ) {
    currentDeployment.livenessProbe.handler.exec.command.map((command, index) =>
      handlerObject.push(
        formatResourceData(
          `Command# ${index}`,
          command,
          `${containerBase}.livenessProbe.handler.exec.command`
        )
      )
    );
  }
  if (
    currentDeployment.livenessProbe.handler.grpc &&
    currentDeployment.livenessProbe.handler.grpc.service
  ) {
    handlerObject.push(
      formatResourceData(
        "GRPC Port",
        currentDeployment.livenessProbe.handler.grpc.port,
        `${containerBase}.livenessProbe.handler.grpc.port`
      )
    );
    handlerObject.push(
      formatResourceData(
        "GRPC Service",
        currentDeployment.livenessProbe.handler.grpc.service,
        `${containerBase}.livenessProbe.handler.grpc.service`
      )
    );
  }
  if (
    currentDeployment.livenessProbe.handler.httpGet &&
    currentDeployment.livenessProbe.handler.httpGet.path
  ) {
    handlerObject.push(
      formatResourceData(
        "HTTP Get Path",
        currentDeployment.livenessProbe.handler.httpGet.path,
        `${containerBase}.livenessProbe.handler.httpGet.path`
      )
    );
    handlerObject.push(
      formatResourceData(
        "HTTP Get Port",
        currentDeployment.livenessProbe.handler.httpGet.port,
        `${containerBase}.livenessProbe.handler.httpGet.port`
      )
    );
    handlerObject.push(
      formatResourceData(
        "HTTP Get Host",
        currentDeployment.livenessProbe.handler.httpGet.host,
        `${containerBase}.livenessProbe.handler.httpGet.host`
      )
    );
    handlerObject.push(
      formatResourceData(
        "HTTP Get Schema",
        currentDeployment.livenessProbe.handler.httpGet.scheme,
        `${containerBase}.livenessProbe.handler.httpGet.scheme`
      )
    );
    if (
      currentDeployment.livenessProbe.handler.httpGet.httpHeaders &&
      currentDeployment.livenessProbe.handler.httpGet.httpHeaders.length > 0
    ) {
      currentDeployment.livenessProbe.handler.httpGet.httpHeaders.map((header, index) => {
        handlerObject.push(
          formatResourceData(
            `HTTP Get Header name ${index}`,
            header.name,
            `${containerBase}.livenessProbe.handler.httpGet.httpHeaders[${index}].name`
          )
        );
        return handlerObject.push(
          formatResourceData(
            `HTTP Get Header value ${index}`,
            header.value,
            `${containerBase}.livenessProbe.handler.httpGet.httpHeaders[${index}].value`
          )
        );
      });
    }
  }
  if (
    currentDeployment.livenessProbe.handler.tcpSocket &&
    currentDeployment.livenessProbe.handler.tcpSocket.port
  ) {
    handlerObject.push(
      formatResourceData(
        "TCP Socket Host",
        currentDeployment.livenessProbe.handler.tcpSocket.host,
        `${containerBase}.livenessProbe.handler.tcpSocket.host`
      )
    );
    handlerObject.push(
      formatResourceData(
        "TCP Socker Port",
        currentDeployment.livenessProbe.handler.tcpSocket.port,
        `${containerBase}.livenessProbe.handler.tcpSocket.port`
      )
    );
  }
  return handlerObject;
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
  const handlerObject = getHandlerObjects(currentDeployment, containerBase);

  const deploymentDataObject = [
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
      name: "Initial Delay Seconds",
      value: currentDeployment.livenessProbe.initialDelaySeconds,
      textFieldLabels: {
        disabledField: "Current Value",
        updatedField: "New Value",
      },
      input: {
        type: "string",
        key: `${containerBase}.livenessProbe.initialDelaySeconds`,
        validation: string().matches(SECONDS_REGEX),
      },
    },
    {
      name: "Period Seconds",
      value: currentDeployment.livenessProbe.periodSeconds,
      textFieldLabels: {
        disabledField: "Current Value",
        updatedField: "New Value",
      },
      input: {
        type: "string",
        key: `${containerBase}.livenessProbe.periodSeconds`,
        validation: string().matches(SECONDS_REGEX),
      },
    },
    {
      name: "Timeout Seconds",
      value: currentDeployment.livenessProbe.timeoutSeconds,
      textFieldLabels: {
        disabledField: "Current Value",
        updatedField: "New Value",
      },
      input: {
        type: "string",
        key: `${containerBase}.livenessProbe.timeoutSeconds`,
        validation: string().matches(SECONDS_REGEX),
      },
    },
    {
      name: "Success Threshold",
      value: currentDeployment.livenessProbe.successThreshold,
      textFieldLabels: {
        disabledField: "Current Value",
        updatedField: "New Value",
      },
      input: {
        type: "string",
        key: `${containerBase}.livenessProbe.successThreshold`,
        validation: string().matches(SECONDS_REGEX),
      },
    },
    {
      name: "Failure Threshold",
      value: currentDeployment.livenessProbe.failureThreshold,
      textFieldLabels: {
        disabledField: "Current Value",
        updatedField: "New Value",
      },
      input: {
        type: "string",
        key: `${containerBase}.livenessProbe.failureThreshold`,
        validation: string().matches(SECONDS_REGEX),
      },
    },
    {
      name: "Termination Grace Period Seconds",
      value: currentDeployment.livenessProbe.terminationGracePeriodSeconds,
      textFieldLabels: {
        disabledField: "Current Value",
        updatedField: "New Value",
      },
      input: {
        type: "string",
        key: `${containerBase}.livenessProbe.terminationGracePeriodSeconds`,
        validation: string().matches(SECONDS_REGEX),
      },
    },
  ];

  const finalDeploymentData = [...deploymentDataObject, ...handlerObject];

  return (
    <WizardStep error={deploymentData.error} isLoading={deploymentData.isLoading}>
      <strong>Deployment Details</strong>
      <MetadataTable onUpdate={update} data={finalDeploymentData} />
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
  deployment.deploymentSpec.template.spec.containers.forEach(container => {
    Object.keys(container.livenessProbe).forEach(livenessAttribute => {
      if (livenessAttribute !== "handler") {
        const newValue = container.livenessProbe[livenessAttribute];
        const oldValue = findContainer({
          deploymentSpec: currentDeploymentData.deploymentSpec,
          containerName: container.name,
        }).livenessProbe[livenessAttribute];
        if (newValue !== oldValue) {
          if (!updatedContainer) {
            updateRows.push({ name: "Container Name", value: container.name });
            updatedContainer = true;
          }
          updateRows.push({
            name: `Old ${livenessAttribute}`,
            value: oldValue,
          });
          updateRows.push({
            name: `New ${livenessAttribute}`,
            value: newValue,
          });
        }
      }
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

const UpdateLiveness: React.FC<WorkflowProps> = ({ heading, resolverType }) => {
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
        return client.post("/v1/k8s/updateDeployment", {
          clientset,
          cluster: deploymentData.cluster,
          namespace: deploymentData.namespace,
          name: deploymentData.name,
          fields: {
            containerProbes: [
              {
                containerName: deploymentData.containerName,
                livenessProbe: {
                  timeoutSeconds: container.livenessProbe.timeoutSeconds,
                  initialDelaySeconds: container.livenessProbe.initialDelaySeconds,
                  periodSeconds: container.livenessProbe.periodSeconds,
                  successThreshold: container.livenessProbe.successThreshold,
                  failureThreshold: container.livenessProbe.failureThreshold,
                  terminationGracePeriodSeconds:
                    container.livenessProbe.terminationGracePeriodSeconds,
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

export default UpdateLiveness;
