import React from "react";
import { useForm } from "react-hook-form";
import {
  Button,
  client,
  Confirmation,
  NotePanel,
  Resolver,
  TextField,
  useWizardContext,
} from "@clutch-sh/core";
import { useDataLayout } from "@clutch-sh/data-layout";
import type { WizardChild } from "@clutch-sh/wizard";
import { Wizard, WizardStep } from "@clutch-sh/wizard";
import { Grid, Select } from "@material-ui/core";

import type { ResolverChild, WorkflowProps } from "./index";

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
  const [targetShardCount, setTargetShardCount] = React.useState("");
  const handleTargetShardCountChange = e => {
    setTargetShardCount(e.target.value);
    resourceData.updateData("targetShardCount", e.target.value);
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

  return (
    <WizardStep error={resourceData.error} isLoading={resourceData.isLoading}>
      <form onSubmit={handleSubmit(onSubmit)}>
        <TextField label="StreamName" name="streamName" value={stream.streamName} />
        <TextField label="Region" name="region" value={stream.region} />
        <TextField
          label="Current Shard Count"
          name="currentShardCount"
          value={stream.currentShardCount}
        />
        Select Target Shard Count:
        <Select
          labelId="TargetShardCount"
          id="targetShardCount"
          value={targetShardCount}
          onChange={handleTargetShardCountChange}
          defaultValue={values[2]}
          displayEmpty
        >
          <option value={values[0]}>{values[0]}</option>
          <option value={values[1]}>{values[1]}</option>
          <option value={values[2]}>{values[2]}</option>
          <option value={values[3]}>{values[3]}</option>
          <option value={values[4]}>{values[4]}</option>
          <option value={values[5]}>{values[5]}</option>
          <option value={values[6]}>{values[6]}</option>
        </Select>
      </form>
      <Grid container justify="center">
        <Button text="Back" onClick={onBack} />
        <Button text="Update" variant="destructive" onClick={onSubmit} />
      </Grid>
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
  const streamData = useDataLayout("streamData");

  return (
    <WizardStep error={streamData.error} isLoading={streamData.isLoading}>
      <Confirmation action="Update" />
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
