import React from "react";
import {
  ButtonGroup,
  client,
  Confirmation,
  MetadataTable,
  NotePanel,
  Resolver,
  useWizardContext,
} from "@clutch-sh/core";
import { useDataLayout } from "@clutch-sh/data-layout";
import type { WizardChild } from "@clutch-sh/wizard";
import { Wizard, WizardStep } from "@clutch-sh/wizard";

import type { ConfirmChild, ResolverChild, WorkflowProps } from ".";

const InstanceIdentifier: React.FC<ResolverChild> = ({ resolverType }) => {
  const { onSubmit } = useWizardContext();
  const resolvedResourceData = useDataLayout("resourceData");

  const onResolve = ({ results }) => {
    // Decide how to process results.
    resolvedResourceData.assign(results[0]);
    onSubmit();
  };

  return <Resolver type={resolverType} searchLimit={1} onResolve={onResolve} />;
};

const InstanceDetails: React.FC<WizardChild> = () => {
  const { onSubmit, onBack } = useWizardContext();
  const resourceData = useDataLayout("resourceData");
  const instance = resourceData.displayValue();

  const data = [
    { name: "Instance ID", value: instance.instanceId },
    { name: "Region", value: instance.region },
    { name: "State", value: instance.state },
    { name: "Instance Type", value: instance.instanceType },
    { name: "Public IP Address", value: instance.publicIpAddress },
    { name: "Private IP Address", value: instance.privateIpAddress },
    { name: "Availability Zone", value: instance.availabilityZone },
  ];

  if (instance.tags) {
    Object.keys(instance.tags).forEach(key => {
      data.push({ name: key, value: instance.tags[key] });
    });
  }

  return (
    <WizardStep error={resourceData.error} isLoading={resourceData.isLoading}>
      <MetadataTable data={data} />
      <ButtonGroup
        buttons={[
          {
            text: "Back",
            onClick: onBack,
          },
          {
            text: "Terminate",
            onClick: onSubmit,
            variant: "destructive",
          },
        ]}
      />
    </WizardStep>
  );
};

const Confirm: React.FC<ConfirmChild> = ({ notes }) => {
  const terminationData = useDataLayout("terminationData");

  return (
    <WizardStep error={terminationData.error} isLoading={terminationData.isLoading}>
      <Confirmation action="Termination" />
      <NotePanel notes={notes} />
    </WizardStep>
  );
};

const TerminateInstance: React.FC<WorkflowProps> = ({ heading, resolverType, notes = [] }) => {
  const dataLayout = {
    resourceData: {},
    terminationData: {
      deps: ["resourceData"],
      hydrator: (resourceData: { instanceId: string; region: string }) => {
        return client.post("/v1/aws/ec2/terminateInstance", {
          instance_id: resourceData.instanceId,
          region: resourceData.region,
        });
      },
    },
  };

  return (
    <Wizard dataLayout={dataLayout} heading={heading}>
      <InstanceIdentifier name="Lookup" resolverType={resolverType} />
      <InstanceDetails name="Modify" />
      <Confirm name="Confirmation" notes={notes} />
    </Wizard>
  );
};

export default TerminateInstance;
