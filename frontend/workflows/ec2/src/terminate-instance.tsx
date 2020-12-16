import React from "react";
import {
  Accordion,
  AccordionDetails,
  Button,
  ButtonGroup,
  client,
  Confirmation,
  MetadataTable,
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

  const metadata = [];
  if (instance.tags) {
    Object.keys(instance.tags).forEach(key => {
      metadata.push({ name: key, value: instance.tags[key] });
    });
  }

  return (
    <WizardStep error={resourceData.error} isLoading={resourceData.isLoading}>
      <MetadataTable data={data} />
      <Accordion title="Metadata">
        <AccordionDetails>
          <MetadataTable data={metadata} />
        </AccordionDetails>
      </Accordion>
      <ButtonGroup justify="flex-end">
        <Button text="Back" variant="neutral" onClick={onBack} />
        <Button text="Terminate" variant="destructive" onClick={onSubmit} />
      </ButtonGroup>
    </WizardStep>
  );
};

const Confirm: React.FC<ConfirmChild> = ({ notes }) => {
  const terminationData = useDataLayout("terminationData");
  const instance = useDataLayout("resourceData").displayValue();
  const formattedNotes = notes?.map(note => <div>{note.text}</div>);

  const data = [
    { name: "Instance ID", value: instance.instanceId },
    { name: "Region", value: instance.region },
  ];

  return (
    <WizardStep error={terminationData.error} isLoading={terminationData.isLoading}>
      <Confirmation action="Termination">{notes && formattedNotes}</Confirmation>
      <MetadataTable data={data} />
    </WizardStep>
  );
};

const TerminateInstance: React.FC<WorkflowProps> = ({ heading, resolverType, notes }) => {
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
