import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import {
  Button,
  ButtonGroup,
  client,
  Confirmation,
  FeatureOff,
  FeatureOn,
  MetadataTable,
  NotePanel,
  Resolver,
  SimpleFeatureFlag,
  useWizardContext,
} from "@clutch-sh/core";
import { useDataLayout } from "@clutch-sh/data-layout";
import { Wizard, WizardStep } from "@clutch-sh/wizard";
import _ from "lodash";

import type { ConfirmChild, ResolverChild, VerifyChild, WorkflowProps } from ".";

const PodIdentifier: React.FC<ResolverChild> = ({ resolverType, notes = [] }) => {
  const { onSubmit } = useWizardContext();
  const resolvedResourceData = useDataLayout("resourceData");
  const resolverInput = useDataLayout("resolverInput");
  const onResolve = ({ results, input }) => {
    // Decide how to process results.
    resolvedResourceData.assign(results[0]);
    resolverInput.assign(input);
    onSubmit();
  };

  return <Resolver type={resolverType} searchLimit={1} onResolve={onResolve} notes={notes} />;
};

const PodDetails: React.FC<VerifyChild> = ({ notes = [] }) => {
  const { onSubmit, onBack } = useWizardContext();
  const resourceData = useDataLayout("resourceData");
  const instance = resourceData.displayValue() as IClutch.k8s.v1.Pod;
  const locationNotes = notes.filter(note => note.location === "verify");

  return (
    <WizardStep error={resourceData.error} isLoading={resourceData.isLoading}>
      <strong>Pod Details</strong>
      <NotePanel notes={locationNotes} />
      <MetadataTable
        data={[
          { name: "Name", value: instance.name },
          { name: "Cluster", value: instance.cluster },
          { name: "Namespace", value: instance.namespace },
          { name: "State", value: _.capitalize(instance.state.toString()) },
          { name: "Node IP Address", value: instance.nodeIp },
          { name: "Pod IP Address", value: instance.podIp },
        ]}
      />
      <ButtonGroup>
        <SimpleFeatureFlag feature="k8sDashOrigin">
          <FeatureOn>
            <Button text="Back" variant="neutral" onClick={() => onBack({ toOrigin: true })} />
          </FeatureOn>
          <FeatureOff>
            <Button text="Back" variant="neutral" onClick={() => onBack()} />
          </FeatureOff>
        </SimpleFeatureFlag>
        <Button text="Delete" variant="destructive" onClick={onSubmit} />
      </ButtonGroup>
    </WizardStep>
  );
};

const Confirm: React.FC<ConfirmChild> = () => {
  const deletionData = useDataLayout("deletionData");
  const podData = useDataLayout("resourceData");
  const { name, cluster, namespace } = podData.displayValue();
  return (
    <WizardStep error={deletionData.error} isLoading={deletionData.isLoading}>
      <Confirmation action="Deletion" />
      <MetadataTable
        data={[
          { name: "Name", value: name },
          { name: "Cluster", value: cluster },
          { name: "Namespace", value: namespace },
        ]}
      />
    </WizardStep>
  );
};

const DeletePod: React.FC<WorkflowProps> = ({ heading, resolverType, notes = [] }) => {
  const dataLayout = {
    resolverInput: {},
    resourceData: {},
    deletionData: {
      deps: ["resourceData", "resolverInput"],
      hydrator: (resourceData: IClutch.k8s.v1.Pod, resolverInput: { clientset: string }) => {
        const clientset = resolverInput.clientset ?? "undefined";
        return client.post("/v1/k8s/deletePod", {
          clientset,
          cluster: resourceData.cluster,
          namespace: resourceData.namespace,
          name: resourceData.name,
        } as IClutch.k8s.v1.DeletePodRequest);
      },
    },
  };

  return (
    <Wizard dataLayout={dataLayout} heading={heading}>
      <PodIdentifier name="Lookup" resolverType={resolverType} notes={notes} />
      <PodDetails name="Verify" notes={notes} />
      <Confirm name="Result" />
    </Wizard>
  );
};

export default DeletePod;
