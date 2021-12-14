import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import {
  Button,
  ButtonGroup,
  client,
  Confirmation,
  MetadataTable,
  Resolver,
  Switch,
  useWizardContext,
} from "@clutch-sh/core";
import { useDataLayout } from "@clutch-sh/data-layout";
import type { WizardChild } from "@clutch-sh/wizard";
import { Wizard, WizardStep } from "@clutch-sh/wizard";

import type { ConfirmChild, ResolverChild, WorkflowProps } from ".";

const NodeIdentifier: React.FC<ResolverChild> = ({ resolverType }) => {
  const { onSubmit } = useWizardContext();
  const resolvedResourceData = useDataLayout("resourceData");
  const resolverInput = useDataLayout("resolverInput");

  const onResolve = ({ results, input }) => {
    // Decide how to process results.
    resolvedResourceData.assign(results[0]);
    resolverInput.assign(input);
    onSubmit();
  };

  return <Resolver type={resolverType} searchLimit={1} onResolve={onResolve} />;
};

const NodeDetails: React.FC<WizardChild> = () => {
  const { onSubmit, onBack } = useWizardContext();
  const resourceData = useDataLayout("resourceData");
  const node = resourceData.displayValue() as IClutch.k8s.v1.Node;
  const update = (_: any, value: boolean) => {
    resourceData.updateData("unschedulable", value);
  };
  return (
    <WizardStep error={resourceData.error} isLoading={resourceData.isLoading}>
      <strong>Node Details</strong>
      <MetadataTable
        onUpdate={update}
        data={[
          { name: "Name", value: node.name },
          { name: "Cluster", value: node.cluster },
          {
            name: "Unschedulable",
            value: <Switch checked={node.unschedulable} onChange={update} />,
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
  const node = useDataLayout("resourceData").displayValue() as IClutch.k8s.v1.Node;
  const updateData = useDataLayout("updateData");
  return (
    <WizardStep error={updateData.error} isLoading={updateData.isLoading}>
      <Confirmation action="Update" />
      <MetadataTable
        data={[
          { name: "Name", value: node.name },
          { name: "Cluster", value: node.cluster },
          { name: "Unschedulable", value: String(node.unschedulable) },
        ]}
      />
    </WizardStep>
  );
};

const CordonNode: React.FC<WorkflowProps> = ({ heading, resolverType, enableFeedback }) => {
  const dataLayout = {
    resolverInput: {},
    resourceData: {},
    updateData: {
      deps: ["resourceData", "resolverInput"],
      hydrator: (resourceData: IClutch.k8s.v1.Node, resolverInput: { clientset: string }) => {
        const clientset = resolverInput.clientset ?? "undefined";
        return client.post("/v1/k8s/updateNode", {
          clientset,
          cluster: resourceData.cluster,
          unschedulable: resourceData.unschedulable,
          name: resourceData.name,
        } as IClutch.k8s.v1.UpdateNodeRequest);
      },
    },
  };

  return (
    <Wizard dataLayout={dataLayout} heading={heading} enableFeedback={enableFeedback}>
      <NodeIdentifier name="Lookup" resolverType={resolverType} />
      <NodeDetails name="Verify" />
      <Confirm name="Confirmation" />
    </Wizard>
  );
};

export default CordonNode;
