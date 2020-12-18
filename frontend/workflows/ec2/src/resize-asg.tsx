import React from "react";
import {
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
import { number } from "yup";

import type { ConfirmChild, ResolverChild, WorkflowProps } from ".";

const GroupIdentifier: React.FC<ResolverChild> = ({ resolverType }) => {
  const { onSubmit } = useWizardContext();
  const groupData = useDataLayout("groupData");

  const onResolve = ({ results }) => {
    // Decide how to process results.
    groupData.assign(results[0]);
    onSubmit();
  };

  return <Resolver type={resolverType} searchLimit={1} onResolve={onResolve} />;
};

const GroupDetails: React.FC<WizardChild> = () => {
  const { onSubmit, onBack } = useWizardContext();
  const groupData = useDataLayout("groupData");
  const group = groupData.displayValue();
  const update = (key: string, value: string) => {
    groupData.updateData(key, value);
  };

  return (
    <WizardStep error={groupData.error} isLoading={groupData.isLoading}>
      <MetadataTable
        onUpdate={update}
        data={[
          { name: "Name", value: group.name },
          { name: "Region", value: group.region },
          { name: "Termination Policy", value: group.terminationPolicies },
          {
            name: "Min Size",
            value: group.size.min,
            input: {
              type: "number",
              key: "size.min",
              validation: number().integer().moreThan(0),
            },
          },
          {
            name: "Max Size",
            value: group.size.max,
            input: { type: "number", key: "size.max" },
          },
          {
            name: "Desired Size",
            value: group.size.desired,
            input: { type: "number", key: "size.desired" },
          },
          { name: "Availability Zone", value: group.zones },
        ]}
      />
      <ButtonGroup justify="flex-end">
        <Button text="Back" variant="neutral" onClick={onBack} />
        <Button text="Resize" variant="destructive" onClick={onSubmit} />
      </ButtonGroup>
    </WizardStep>
  );
};

// TODO (sperry): possibly show the previous size values
const Confirm: React.FC<ConfirmChild> = ({ notes }) => {
  const group = useDataLayout("groupData").displayValue();
  const resizeData = useDataLayout("resizeData");
  const formattedNotes = notes?.map(note => <div>{note.text}</div>);

  return (
    <WizardStep error={resizeData.error} isLoading={resizeData.isLoading}>
      <Confirmation action="Resize">{notes && formattedNotes}</Confirmation>
      <MetadataTable
        data={[
          { name: "Name", value: group.name },
          { name: "New Min Size", value: group.size.min },
          { name: "New Max Size", value: group.size.max },
          { name: "New Desired Size", value: group.size.desired },
        ]}
      />
    </WizardStep>
  );
};

const ResizeAutoscalingGroup: React.FC<WorkflowProps> = ({ heading, resolverType, notes }) => {
  const dataLayout = {
    groupData: {},
    resizeData: {
      deps: ["groupData"],
      hydrator: (groupData: { name: string; region: string; size: string }) => {
        return client.post("/v1/aws/ec2/resizeAutoscalingGroup", {
          name: groupData.name,
          region: groupData.region,
          size: groupData.size,
        });
      },
    },
  };

  return (
    <Wizard dataLayout={dataLayout} heading={heading}>
      <GroupIdentifier name="Lookup" resolverType={resolverType} />
      <GroupDetails name="Modify" />
      <Confirm name="Confirmation" notes={notes} />
    </Wizard>
  );
};

export default ResizeAutoscalingGroup;
