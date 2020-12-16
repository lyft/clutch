import React from "react";
import { useForm } from "react-hook-form";
import {
  Button,
  ButtonGroup,
  client,
  Confirmation,
  MetadataTable,
  NotePanel,
  Resolver,
  Select,
  TextField,
  useWizardContext,
} from "@clutch-sh/core";
import { useDataLayout } from "@clutch-sh/data-layout";
import type { WizardChild } from "@clutch-sh/wizard";
import { Wizard, WizardStep } from "@clutch-sh/wizard";
import styled from "@emotion/styled";
import { Grid } from "@material-ui/core";

import type { ResolverChild, WorkflowProps } from "./index";

const Form = styled.form({
  alignItems: "center",
  display: "flex",
  flexDirection: "column",
  justifyItems: "space-evenly",
  "> *": {
    padding: "8px 0",
  },
});

const StreamIdentifier: React.FC<ResolverChild> = ({ resolverType }) => {
  const { onSubmit } = useWizardContext();
  const resolvedResourceData = useDataLayout("resourceData");

  const onResolve = ({ results }) => {
    // Decide how to process results.
    resolvedResourceData.assign(results[0]);
    onSubmit();
  };

  return <Resolver type={resolverType} searchLimit={1} onResolve={onResolve} />;
};

const StreamDetails: React.FC<WizardChild> = () => {
  const { handleSubmit } = useForm({
    mode: "onChange",
  });
  const { onSubmit, onBack } = useWizardContext();
  const resourceData = useDataLayout("resourceData");
  const stream = resourceData.displayValue();
  const handleTargetShardCountChange = (value: string) => {
    resourceData.updateData("targetShardCount", value);
  };

  const values = [
    Math.ceil(stream.currentShardCount * 0.5),
    Math.ceil(stream.currentShardCount * 0.75),
    Math.ceil(stream.currentShardCount * 1),
    Math.ceil(stream.currentShardCount * 1.25),
    Math.ceil(stream.currentShardCount * 1.5),
    Math.ceil(stream.currentShardCount * 1.75),
    Math.ceil(stream.currentShardCount * 2),
  ];
  const options = values.map(value => {
    return { label: value.toString() };
  });
  return (
    <WizardStep error={resourceData.error} isLoading={resourceData.isLoading}>
      <Form onSubmit={handleSubmit(onSubmit)}>
        <TextField readOnly label="StreamName" name="streamName" value={stream.streamName} />
        <TextField readOnly label="Region" name="region" value={stream.region} />
        <Grid container alignItems="stretch" wrap="nowrap">
          <Grid item style={{ flexBasis: "50%", paddingRight: "8px" }}>
            <TextField
              readOnly
              label="Current Shard Count"
              name="currentShardCount"
              value={stream.currentShardCount}
              disabled
            />
          </Grid>
          <Grid item style={{ flexBasis: "50%", paddingRight: "8px" }}>
            <Select
              label="TargetShardCount"
              name="targetShardCount"
              onChange={handleTargetShardCountChange}
              options={options}
              defaultOption={2}
            />
          </Grid>
        </Grid>
      </Form>

      <ButtonGroup>
        <Button text="Back" onClick={onBack} />
        <Button text="Update" variant="destructive" onClick={onSubmit} />
      </ButtonGroup>
      <NotePanel
        notes={[
          {
            severity: "info",
            text:
              "These changes are not immediate. Expect some delay length correlated to the size of the stream",
          },
        ]}
      />
    </WizardStep>
  );
};

const Confirm: React.FC<WizardChild> = () => {
  const updateStreamData = useDataLayout("streamData");
  const streamData = useDataLayout("resourceData").value;
  return (
    <WizardStep error={updateStreamData.error} isLoading={updateStreamData.isLoading}>
      <Confirmation action="Update" />
      <MetadataTable
        data={[
          { name: "Stream Name", value: streamData.streamName },
          { name: "Region", value: streamData.region },
          { name: "Current Shard Count", value: streamData.currentShardCount },
          { name: "Target Shard Count", value: streamData.targetShardCount },
        ]}
      />
    </WizardStep>
  );
};

const UpdateShardCount: React.FC<WorkflowProps> = ({ heading, resolverType }) => {
  const dataLayout = {
    resourceData: {},
    streamData: {
      deps: ["resourceData"],
      hydrator: (resourceData: {
        streamName: string;
        region: string;
        targetShardCount: number;
      }) => {
        return client.post("/v1/aws/kinesis/updateShardCount", {
          stream_name: resourceData.streamName,
          region: resourceData.region,
          target_shard_count: resourceData.targetShardCount,
        });
      },
    },
  };

  return (
    <Wizard dataLayout={dataLayout} heading={heading}>
      <StreamIdentifier name="Lookup" resolverType={resolverType} />
      <StreamDetails name="Modify" />
      <Confirm name="Confirmation" />
    </Wizard>
  );
};

export default UpdateShardCount;
