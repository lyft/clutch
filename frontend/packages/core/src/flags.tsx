import * as React from "react";

import { client } from "./Network/index";

export interface FeatureFlags {
  [name: string]: any;
}

const featureFlags = (): FeatureFlags => {
  return client
    .post("/v1/featureflag/getFlags")
    .then(response => {
      return response.data.flags;
    })
    .catch(error => {
      // TODO: add instrumentation here
      console.warn("failed to fetch flags: ", error); // eslint-disable-line
      return {};
    });
};

export interface FeatureFlagProps {
  flagName: string;
  children: React.ReactNode;
}

const FeatureFlag = ({ flagName, children }: FeatureFlagProps) => {
  const [flags, setFlags] = React.useState({});
  featureFlags().then(f => setFlags(f));
  const flag = flags?.[flagName];
  if (flag !== undefined && flag.booleanValue === true) {
    return <>children</>;
  }

  return null;
};

export { featureFlags, FeatureFlag };
