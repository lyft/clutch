import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import {
  ButtonGroup,
  client,
  Confirmation,
  ExpandableTable,
  ExpandableRow,
  MetadataTable,
  Resolver,
  useWizardContext,
} from "@clutch-sh/core";
import { useDataLayout } from "@clutch-sh/data-layout";
import type { WizardChild } from "@clutch-sh/wizard";
import { Wizard, WizardStep } from "@clutch-sh/wizard";
import _ from "lodash";

import type { ConfirmChild, ResolverChild, WorkflowProps } from ".";

const NamespaceIdentifier: React.FC<ResolverChild> = ({ resolverType }) => {
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

const ServiceData: React.FC<WizardChild> = () => {
  const resourceData = useDataLayout("resourceData");
  const instance = resourceData.displayValue() as IClutch.k8s.v1.Service;

  return (
    <WizardStep error={resourceData.error} isLoading={resourceData.isLoading}>
       <ExpandableTable headings = {["Name", "Type", "Cluster IP", "External IP", "Port(s)", "Age"]}/>
      <ExpandableRow key={instance.name} heading={instance.name} summary={instance.name}/>
    
    </WizardStep>
    
  );
};

const ListServices: React.FC<WorkflowProps> = ({ heading, resolverType }) => {
  const dataLayout = {
    resolverInput: {},
    resourceData: {},
    deletionData: {
      deps: ["resourceData", "resolverInput"],
      hydrator: (resourceData: IClutch.k8s.v1.Service, resolverInput: { clientset: string }) => {
        const clientset = resolverInput.clientset ?? "undefined";
        return client.post("/v1/k8s/listServices", {
          clientset,
          cluster: resourceData.cluster,
          namespace: resourceData.namespace,
        //   name: resourceData.name,
        } as IClutch.k8s.v1.ListServicesRequest);
      },
    },
  };

  return (
    <Wizard dataLayout={dataLayout} heading={heading}>
      <NamespaceIdentifier name="Lookup" resolverType={resolverType} />
      <ServiceData name="Services" />
    </Wizard>
  );
};

export default ListServices;
