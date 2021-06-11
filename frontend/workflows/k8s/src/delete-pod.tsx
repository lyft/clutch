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
  Resolver,
  SimpleFeatureFlag,
  useWizardContext,
} from "@clutch-sh/core";
import { useDataLayout } from "@clutch-sh/data-layout";
import type { WizardChild } from "@clutch-sh/wizard";
import { Wizard, WizardStep } from "@clutch-sh/wizard";
import _ from "lodash";

import type { ConfirmChild, ResolverChild, WorkflowProps } from ".";

const PodIdentifier: React.FC<ResolverChild> = ({ resolverType }) => {
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

const PodDetails: React.FC<WizardChild> = () => {
  const { onSubmit, onBack } = useWizardContext();
  const resourceData = useDataLayout("resourceData");
  const instance = resourceData.displayValue() as IClutch.k8s.v1.Pod;

  return (
    <WizardStep error={resourceData.error} isLoading={resourceData.isLoading}>
      <strong>Pod Details</strong>
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

/*
TODO: Need information boxes for
  These changes are not permanent, and will be overwritten on your next deploy. Adjust your manifest.yaml to persist changes across deploys.
and
  Note: the HPA should take just a few minutes to scale in either direction.
*/
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

const DeletePod: React.FC<WorkflowProps> = ({ heading, resolverType }) => {
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
      <PodIdentifier name="Lookup" resolverType={resolverType} />
      <PodDetails name="Verify" />
      <Confirm name="Confirmation" />
    </Wizard>
  );
};

export default DeletePod;
