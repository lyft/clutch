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
import { Wizard, WizardStep } from "@clutch-sh/wizard";
import _ from "lodash";

const PodIdentifier = ({ resolverType }) => {
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

const PodDetails = () => {
  const { onSubmit, onBack } = useWizardContext();
  const resourceData = useDataLayout("resourceData");
  const instance = resourceData.displayValue();

  return (
    <WizardStep error={resourceData.error} isLoading={resourceData.isLoading}>
      <MetadataTable
        data={[
          { name: "Name", value: instance.name },
          { name: "Cluster", value: instance.cluster },
          { name: "Namespace", value: instance.namespace },
          { name: "State", value: _.capitalize(instance.state) },
          { name: "Node IP Address", value: instance.nodeIp },
          { name: "Pod IP Address", value: instance.podIp },
        ]}
      />
      <ButtonGroup
        buttons={[
          {
            text: "Back",
            onClick: onBack,
          },
          {
            text: "Delete",
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
const Confirm = () => {
  const deletionData = useDataLayout("deletionData");

  return (
    <WizardStep error={deletionData.error} isLoading={deletionData.isLoading}>
      <Confirmation action="Deletion" />
    </WizardStep>
  );
};

const DeletePod = ({ heading, resolverType }) => {
  const dataLayout = {
    resolverInput: {},
    resourceData: {},
    deletionData: {
      deps: ["resourceData", "resolverInput"],
      hydrator: (resourceData, resolverInput) => {
        const clientset = resolverInput.clientset ?? "undefined";
        return client.post("/v1/k8s/deletePod", {
          clientset,
          cluster: resourceData.cluster,
          namespace: resourceData.namespace,
          name: resourceData.name,
        });
      },
    },
  };

  return (
    <Wizard dataLayout={dataLayout} heading={heading}>
      <PodIdentifier name="Lookup" resolverType={resolverType} />
      <PodDetails name="Modify" />
      <Confirm name="Confirmation" />
    </Wizard>
  );
};

export default DeletePod;
