import React from "react";
import {
  BrowserRouter as Router,
  Outlet,
  Route,
  Routes,
  useLocation,
  useNavigate,
} from "react-router-dom";
import _ from "lodash";

import AppLayout from "../AppLayout";
import { ApplicationContext } from "../Contexts/app-context";
import type { HydrateData, HydratedData } from "../Contexts/storage-context";
import { StorageContext } from "../Contexts/storage-context";
import { FEATURE_FLAG_POLL_RATE, featureFlags } from "../flags";
import Landing from "../landing";
import NotFound from "../not-found";

import { registeredWorkflows } from "./registrar";
import { Theme } from "./themes";
import type { ConfiguredRoute, Workflow, WorkflowConfiguration } from "./workflow";
import ErrorBoundary from "./workflow";

export interface UserConfiguration {
  [packageName: string]: {
    [key: string]: ConfiguredRoute;
  };
}

interface ClutchAppProps {
  availableWorkflows: {
    [key: string]: () => WorkflowConfiguration;
  };
  configuration: UserConfiguration;
}

interface ShortLinkHydratorProps {
  hydrate: (HydrateData) => void;
}

/**
 * Filter workflow routes using available feature flags.
 * @param workflows a list of valid Workflow objects.
 */
const featureFlagFilter = (workflows: Workflow[]): Promise<Workflow[]> => {
  return featureFlags().then(flags =>
    workflows.filter(workflow => {
      /* eslint-disable-next-line no-param-reassign */
      workflow.routes = workflow.routes.filter(route => {
        const show =
          route.featureFlag === undefined ||
          (flags?.[route.featureFlag] !== undefined &&
            flags[route.featureFlag].booleanValue === true);
        return show;
      });
      return workflow.routes.length !== 0;
    })
  );
};

const ShortLinkHydrator = ({ hydrate }: ShortLinkHydratorProps) => {
  const { pathname } = useLocation();
  const navigate = useNavigate();

  const returnedData: { [key: string]: HydrateData } = {
    "1234": {
      route: "/dash",
      data: {
        ProjectSelector: {
          dashState: {
            "0": {
              clutch: { checked: true },
              clutchdata: { checked: true },
              lyftkube: { checked: true },
            },
            "1": {
              acl: { checked: false },
              tcseditor: { checked: false },
              cdcconnector: { checked: false },
              deployapi: { checked: false },
              incidentassistant: { checked: false },
              infraevents: { checked: false },
              dynamodbquerier: { checked: false },
              omnibot: { checked: false },
              servicecatalog: { checked: false },
            },
            "2": {
              deployapi: { checked: false },
              sloportal: { checked: false },
              loadtestbot: { checked: false },
              omnibot: { checked: false },
            },
          },
        },
        TimelineView: {
          dashTimelineEventFilters: ["alerts", "deploys"],
          dashSplitEvents: false,
        },
      },
    },
    "2456": {
      route: "/dash",
      data: {
        ProjectSelector: {
          dashState: {
            "0": {
              lyftkube: { checked: true },
              ridesapi: { checked: true, custom: true },
            },
            "1": {
              acl: { checked: false },
              tcseditor: { checked: false },
              cdcconnector: { checked: true },
              deployapi: { checked: false },
              incidentassistant: { checked: false },
              infraevents: { checked: false },
              dynamodbquerier: { checked: false },
              omnibot: { checked: false },
              servicecatalog: { checked: false },
              jobscheduler: { checked: false },
              ridestate: { checked: false },
              rateandpay: { checked: false },
              ridescheduler: { checked: false },
              badcompanions: { checked: false },
              dsfeatures: { checked: false },
              riderperks: { checked: false },
              locations: { checked: false },
              ridehistory: { checked: false },
              driverearnings: { checked: false },
              enterprisedispatch: { checked: false },
              trips: { checked: false },
              userpreferences: { checked: false },
              cancels: { checked: false },
              drivermode: { checked: false },
              reflection: { checked: false },
              enterprise: { checked: false },
              messaging: { checked: false },
              green: { checked: false },
              viewport: { checked: false },
              karma: { checked: false },
              wallet: { checked: false },
              taskqueue: { checked: false },
              pusher: { checked: false },
              hourscompliance: { checked: false },
              rideprograms: { checked: false },
              drivermonitor: { checked: false },
              driverqueues: { checked: false },
              routestate: { checked: false },
              businessprograms: { checked: false },
              fare: { checked: false },
              ratecards: { checked: false },
              ratelimit: { checked: false },
              cancelproperties: { checked: false },
              supply: { checked: false },
              routesapi: { checked: false },
              payapi: { checked: false },
              places: { checked: false },
              offerings: { checked: false },
              legacyridetypes: { checked: false },
              push: { checked: false },
              autonomous: { checked: false },
              enterprisepayapi: { checked: false },
              deliveries: { checked: false },
              matchingexecution: { checked: false },
              compliance: { checked: true },
              users: { checked: false },
              custodian: { checked: false },
              etaproxy: { checked: false },
              ledger: { checked: false },
              tcs: { checked: false },
              routes: { checked: false },
              asynctransitionupdates: { checked: false },
              txnhub: { checked: false },
              ufo: { checked: false },
              dispatch: { checked: false },
              locationssearch: { checked: false },
              matchingstate: { checked: false },
              pricing: { checked: false },
              regions: { checked: false },
              walletapi: { checked: false },
              sharedpayments: { checked: false },
              tripsapi: { checked: false },
              ratings: { checked: false },
              rides: { checked: false },
              coupons: { checked: false },
              venues: { checked: false },
              usergroups: { checked: false },
              pickupoptimizer: { checked: false },
              switchboard: { checked: false },
              passengerqueues: { checked: false },
              ridelocations: { checked: false },
            },
            "2": {
              deployapi: { checked: false },
              sloportal: { checked: false },
              loadtestbot: { checked: false },
              omnibot: { checked: false },
              displaycomponents: { checked: false },
              passengerqueues: { checked: false },
              wwwdriverxpfe: { checked: false },
              autonomousinride: { checked: false },
              cancels: { checked: false },
              drivermode: { checked: false },
              riderdelight: { checked: false },
              walletapi: { checked: false },
              incentiveguarantees: { checked: false },
              graphql: { checked: false },
              venues: { checked: false },
              messaging: { checked: false },
              rateandpay: { checked: false },
              migrationworkers: { checked: false },
              legacyridescheduler: { checked: false },
              payapi: { checked: false },
              www2: { checked: false },
              incentiveshistory: { checked: false },
              lyftentertainment: { checked: false },
              karma: { checked: false },
              routesapi: { checked: false },
              deliveries: { checked: false },
              drivermonitor: { checked: false },
              ledger: { checked: false },
              resourcemgmt: { checked: false },
              sandbox: { checked: false },
              autonomouspartners: { checked: false },
              ridechat: { checked: false },
              campaignscheduler: { checked: false },
              ridestate: { checked: false },
              pusher: { checked: false },
              banners: { checked: false },
              ratings: { checked: false },
              tripsapi: { checked: false },
              autonomousvsm: { checked: false },
              sharelocation: { checked: false },
              custodian: { checked: false },
              enterprisedispatch: { checked: false },
              matchingexecution: { checked: false },
              robocop: { checked: false },
              green: { checked: false },
              ridescheduler: { checked: false },
              routestate: { checked: false },
              autonomous: { checked: false },
              tripintent: { checked: false },
              wwwenterprisefe: { checked: false },
              ampdevicemanager: { checked: false },
            },
          },
        },
      },
    },
  };

  React.useEffect(() => {
    const matches = pathname.match(/(\/sl\/)(.*)/i);
    if (matches && matches[2] && returnedData[matches[2]]) {
      const data = returnedData[matches[2]];
      // console.log("HYDRATING", data.data);
      hydrate(data.data);
      navigate(data.route);
    }
  }, [pathname]);

  return null;
};

const ClutchApp: React.FC<ClutchAppProps> = ({
  availableWorkflows,
  configuration: userConfiguration,
}) => {
  const [workflows, setWorkflows] = React.useState<Workflow[]>([]);
  const [isLoading, setIsLoading] = React.useState<boolean>(true);
  const [hydrateStore, setHydrateStore] = React.useState<HydratedData>({});
  const [tempHydrateStore, setTempHydrateStore] = React.useState<HydratedData>({});

  const loadWorkflows = () => {
    registeredWorkflows(availableWorkflows, userConfiguration, [featureFlagFilter]).then(w => {
      setWorkflows(w);
      setIsLoading(false);
    });
  };

  const localStore = (key: string, data: any) => {
    // console.log("Storing Local Data", { key, data });
    if (key) {
      try {
        window.localStorage.setItem(key, JSON.stringify(data));
      } catch (e) {
        // eslint-disable-next-line no-console
        console.error(e);
      }
    }
  };

  const store = (componentName: string, key: string, data: any, local = true) => {
    // console.log("Storing Data", { componentName, key, data, local });
    // we're clearing our temporary data
    if (!componentName && !key) {
      setTempHydrateStore({});
    } else {
      if (componentName) {
        const newTempHydrate = { ...tempHydrateStore };

        if (!newTempHydrate[componentName]) {
          newTempHydrate[componentName] = {};
        }

        if (key) {
          // if we have a specific key, set it directly to the data
          newTempHydrate[componentName][key] = data;
        } else {
          // if we dont have a specific key, we'll just extend the data onto the component
          newTempHydrate[componentName] = { ...newTempHydrate[componentName], ...data };
        }
        setTempHydrateStore(newTempHydrate);
      }
      if (key && local) {
        localStore(key, data);
      }
    }
  };

  const retrieve = (componentName: string, key?: string, defaultData?: any): any => {
    // console.group("Retrieval");
    // console.log({ componentName, key, defaultData });
    if (hydrateStore && hydrateStore[componentName]) {
      // console.log("Verified componentName", hydrateStore[componentName]);
      if (key && hydrateStore[componentName][key]) {
        // console.log("Have a key", key);
        // console.log("returning", hydrateStore[componentName][key]);
        // console.groupEnd();
        return hydrateStore[componentName][key];
      }
      if (!key) {
        // console.log("returning", hydrateStore[componentName]);
        // console.groupEnd();
        return hydrateStore[componentName];
      }
    }

    if (key) {
      const localData = window.localStorage.getItem(key);
      if (localData) {
        try {
          // console.log("returning local data", localData);
          // console.groupEnd();
          return JSON.parse(localData);
        } catch (e) {
          // eslint-disable-next-line no-console
          console.error(e);
        }
      }
    }

    // console.log("Returning default", defaultData);
    // console.groupEnd();

    return defaultData;
  };

  const remove = (componentName: string, key?: string, local = true) => {
    // console.log("Removing Data", { componentName, key, local });

    const newTempHydrate = { ...tempHydrateStore };
    if (componentName && key) {
      delete newTempHydrate[componentName][key];
    } else if (componentName) {
      delete newTempHydrate[componentName];
    }
    setTempHydrateStore(newTempHydrate);

    if (key && local) {
      window.localStorage.removeItem(key);
    }
  };

  React.useEffect(() => {
    loadWorkflows();
    const interval = setInterval(loadWorkflows, FEATURE_FLAG_POLL_RATE);
    return () => clearInterval(interval);
  }, []);

  const [discoverableWorkflows, setDiscoverableWorkflows] = React.useState([]);
  React.useEffect(() => {
    /** Filter out all of the workflows that are configured to be `hideNav: true`.
     * This prevents the workflows from being discoverable by the user from the UI,
     * both search and drawer navigation.
     *
     * The routes for all configured workflows will still be reachable
     * by manually providing the full path in the URI.
     */
    const pw = _.cloneDeep(workflows).filter(workflow => {
      const publicRoutes = workflow.routes.filter(route => {
        return !(route?.hideNav !== undefined ? route.hideNav : false);
      });
      workflow.routes = publicRoutes; /* eslint-disable-line no-param-reassign */
      return publicRoutes.length !== 0;
    });
    setDiscoverableWorkflows(pw);
  }, [workflows]);

  return (
    <Router>
      <Theme variant="light">
        <div id="App">
          <ApplicationContext.Provider value={{ workflows: discoverableWorkflows }}>
            <StorageContext.Provider
              value={{
                hydrateStore,
                tempHydrateStore,
                data: { retrieve, store, localStore, remove },
              }}
            >
              <Routes>
                <Route path="/" element={<AppLayout isLoading={isLoading} />}>
                  <Route key="landing" path="" element={<Landing />} />
                  {workflows.map((workflow: Workflow) => {
                    const workflowPath = workflow.path.replace(/^\/+/, "").replace(/\/+$/, "");
                    const workflowKey = workflow.path.split("/")[0];
                    return (
                      <Route
                        path={`${workflowPath}/`}
                        key={workflowKey}
                        element={
                          <ErrorBoundary workflow={workflow}>
                            <Outlet />
                          </ErrorBoundary>
                        }
                      >
                        {workflow.routes.map(route => {
                          const heading = route.displayName
                            ? `${workflow.displayName}: ${route.displayName}`
                            : workflow.displayName;
                          return (
                            <Route
                              key={workflow.path}
                              path={`${route.path.replace(/^\/+/, "").replace(/\/+$/, "")}`}
                              element={React.cloneElement(<route.component />, {
                                ...route.componentProps,
                                heading,
                              })}
                            />
                          );
                        })}
                        <Route key={`${workflow.path}/notFound`} path="*" element={<NotFound />} />
                      </Route>
                    );
                  })}
                  <Route
                    key="short-links"
                    path="/sl/*"
                    element={<ShortLinkHydrator hydrate={setHydrateStore} />}
                  />
                  <Route key="notFound" path="*" element={<NotFound />} />
                </Route>
              </Routes>
            </StorageContext.Provider>
          </ApplicationContext.Provider>
        </div>
      </Theme>
    </Router>
  );
};

export default ClutchApp;
