import * as React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";

import type { ClutchError } from "./Network/errors";
import { client } from "./Network/index";

const FEATURE_FLAG_POLL_RATE = +process.env.REACT_APP_FF_POLL || 300000;
const FF_CACHE_TTL = +process.env.REACT_APP_FF_CACHE_TTL || 60000;

export interface FeatureFlags {
  [name: string]: IClutch.featureflag.v1.IFlag;
}

const featureFlags = (): Promise<FeatureFlags> => {
  const cachedFlags = JSON.parse(sessionStorage.getItem("featureFlags"));
  if (cachedFlags) {
    const elapsedTime = new Date(new Date().getTime() - cachedFlags.timestamp).getTime();
    if (elapsedTime < FF_CACHE_TTL) {
      return new Promise(resolve => resolve(cachedFlags.flags));
    }
  }
  return client
    .post("/v1/featureflag/getFlags", {}, { timeout: 1000 })
    .then(response => {
      const { flags } = response.data as IClutch.featureflag.v1.GetFlagsResponse;
      const cache = {
        timestamp: new Date().getTime(),
        flags,
      };
      sessionStorage.setItem("featureFlags", JSON.stringify(cache));
      return flags;
    })
    .catch((error: ClutchError) => {
      // TODO: add instrumentation here
      console.warn("failed to fetch flags: ", error); // eslint-disable-line
      return cachedFlags?.flags || {};
    });
};

interface SimpleFeatureFlagProps {
  /** The name of the feature flag to lookup */
  feature: string;
  /** A simple feature flag component */
  children:
    | React.ReactElement<SimpleFeatureFlagProps>
    | React.ReactElement<SimpleFeatureFlagProps>[];
}

interface GenericSimpleFeatureFlagStateProps {
  /** Children that will be rendered if the feature is in the expected state. */
  children: React.ReactNode;
  /** If the feature is enabled. */
  enabled: boolean;
  /** The expected feature state to render the children. */
  expectedState: boolean;
}

/**
 * A generic simple feature flag component. This exists to mask the feature enabled
 * state as a prop on the exposed On/Off components.
 */
const GenericSimpleFeatureFlagState = ({
  children,
  enabled,
  expectedState,
}: GenericSimpleFeatureFlagStateProps) => <>{enabled === expectedState && children}</>;

interface SimpleFeatureFlagStateProps {
  /** Children that will be rendered if the feature is on/off */
  children: React.ReactNode;
}

const FeatureOn = ({ children }: SimpleFeatureFlagStateProps) => <>{children}</>;

const FeatureOff = ({ children }: SimpleFeatureFlagStateProps) => <>{children}</>;

/**
 * A feature flag wrapper that evaluates a binary value of a specified flag to determine
 * if it's children should be shown.
 */
const SimpleFeatureFlag = ({ feature, children }: SimpleFeatureFlagProps) => {
  const cachedFlags = JSON.parse(sessionStorage.getItem("featureFlags"));
  const [flags, setFlags] = React.useState(cachedFlags?.flags || {});
  const [featureEnabled, setFeatureEnabled] = React.useState(false);

  const loadFlags = () => {
    featureFlags().then(f => setFlags(f));
  };

  React.useEffect(() => {
    loadFlags();
    const interval = setInterval(loadFlags, FEATURE_FLAG_POLL_RATE);
    return () => clearInterval(interval);
  }, []);

  React.useEffect(() => {
    const flag = flags?.[feature];
    if (flag !== undefined) {
      setFeatureEnabled(flag.booleanValue);
    }
  }, [flags]);

  const statefulChildren = React.Children.map(children, child => (
    <GenericSimpleFeatureFlagState
      enabled={featureEnabled}
      expectedState={child.type === FeatureOn}
    >
      {child.props.children}
    </GenericSimpleFeatureFlagState>
  ));

  return <>{statefulChildren}</>;
};

export { FEATURE_FLAG_POLL_RATE, featureFlags, FeatureOff, FeatureOn, SimpleFeatureFlag };
