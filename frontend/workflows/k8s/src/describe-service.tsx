import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
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
import _ from "lodash";

import type { ConfirmChild, ResolverChild, WorkflowProps } from ".";

const ServiceIdentifier: React.FC<ResolverChild> = ({ resolverType }) => {
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

const ServiceDetails: React.FC<WizardChild> = () => {
  const resourceData = useDataLayout("resourceData");
  const instance = resourceData.displayValue() as IClutch.k8s.v1.Service;

  return (
    <WizardStep error={resourceData.error} isLoading={resourceData.isLoading}>
      <MetadataTable
        data={[
          { name: "Name", value: instance.name },
          { name: "Cluster", value: instance.cluster },
          { name: "Namespace", value: instance.namespace },
          { name: "Available info", value: instance.data },
          { name: "Cluster IP Address", value: instance.clusterIp },
          // { name: "Labels", value: instance.labels },
          // { name: "Annotations", value: instance.annotations },
        ]}
      />
      
    </WizardStep>
  );
};

const DescribeService: React.FC<WorkflowProps> = ({ heading, resolverType }) => {
  const dataLayout = {
    resolverInput: {},
    resourceData: {},
    deletionData: {
      deps: ["resourceData", "resolverInput"],
      hydrator: (resourceData: IClutch.k8s.v1.Service, resolverInput: { clientset: string }) => {
        const clientset = resolverInput.clientset ?? "undefined";
        return client.post("/v1/k8s/describeService", {
          clientset,
          cluster: resourceData.cluster,
          namespace: resourceData.namespace,
          name: resourceData.name,
        } as IClutch.k8s.v1.DescribeServiceRequest);
      },
    },
  };

  return (
    <Wizard dataLayout={dataLayout} heading={heading}>
      <ServiceIdentifier name="Lookup" resolverType={resolverType} />
      <ServiceDetails name="Details" />
    </Wizard>
  );
};

export default DescribeService;
