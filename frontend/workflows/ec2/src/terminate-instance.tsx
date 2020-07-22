import React from "react";
import {
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

import type { ResolverChild, WorkflowProps } from ".";

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

  return (
    <WizardStep error={resourceData.error} isLoading={resourceData.isLoading}>
      <MetadataTable
        data={[
          { name: "Instance ID", value: instance.instanceId },
          { name: "Region", value: instance.region },
          { name: "State", value: instance.state },
          { name: "Instance Type", value: instance.instanceType },
          { name: "Public IP Address", value: instance.publicIpAddress },
          { name: "Private IP Address", value: instance.privateIpAddress },
          { name: "Availability Zone", value: instance.availabilityZone },
        ]}
      />
      <ButtonGroup
        buttons={[
          {
            text: "Back",
            onClick: onBack,
          },
          {
            text: "Terminate",
            onClick: onSubmit,
            destructive: true,
          },
        ]}
      />
    </WizardStep>
  );
};

/*
TODO: Need information boxes for
  These changes are not permanent, and will be overwritten on your next deploy. Adjust your manifest.yaml to persist changes across deploys.
and
  Note: the HPA should take just a few minutes to scale in either direction.
*/
const Confirm: React.FC<WizardChild> = () => {
  const terminationData = useDataLayout("terminationData");

  return (
    <WizardStep error={terminationData.error} isLoading={terminationData.isLoading}>
      <Confirmation action="Termination" />
    </WizardStep>
  );
};

const TerminateInstance: React.FC<WorkflowProps> = ({ heading, resolverType }) => {
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
      <Confirm name="Confirmation" />
    </Wizard>
  );
};

export default TerminateInstance;
