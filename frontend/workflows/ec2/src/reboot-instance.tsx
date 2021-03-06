import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import {
  Button,
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
  const instance = resourceData.displayValue() as IClutch.aws.ec2.v1.Instance;

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
      <ButtonGroup>
        <Button text="Back" variant="neutral" onClick={onBack} />
        <Button text="Reboot" variant="destructive" onClick={onSubmit} />
      </ButtonGroup>
    </WizardStep>
  );
};

const Confirm: React.FC<ConfirmChild> = ({ notes }) => {
  const rebootData = useDataLayout("rebootData");

  return (
    <WizardStep error={rebootData.error} isLoading={rebootData.isLoading}>
      <Confirmation action="Reboot" />
      <NotePanel notes={notes} />
    </WizardStep>
  );
};

const RebootInstance: React.FC<WorkflowProps> = ({ heading, resolverType, notes = [] }) => {
  const dataLayout = {
    resourceData: {},
    rebootData: {
      deps: ["resourceData"],
      hydrator: (resourceData: { instanceId: string; region: string }) => {
        return client.post("/v1/aws/ec2/rebootInstance", {
          instance_id: resourceData.instanceId,
          region: resourceData.region,
        });
      },
    },
  };

  return (
    <Wizard dataLayout={dataLayout} heading={heading}>
      <InstanceIdentifier name="Lookup" resolverType={resolverType} />
      <InstanceDetails name="Verify" />
      <Confirm name="Confirmation" notes={notes} />
    </Wizard>
  );
};

export default RebootInstance;
