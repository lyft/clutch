import type React from "react";

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
  name: string;
  children: React.ReactNode;
}

const FeatureFlag = ({ name, children }: FeatureFlagProps) => {
  featureFlags().then(flags => {
    const flag = flags.find((f: { [name: string]: any }) => f.name === name);

    if (flag && flag.active) {
      return children;
    }

    return null;
  });
};

export { featureFlags, FeatureFlag };
