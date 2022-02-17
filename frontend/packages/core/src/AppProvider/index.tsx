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
import type { HydrateData, HydratedData } from "../Contexts/short-link-context";
import { ShortLinkContext } from "../Contexts/short-link-context";
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
        dash: {
          state: {
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
          splitEvents: true,
        },
      },
    },
    "2456": {
      route: "/dash",
      data: {
        dash: {
          state: {
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
          splitEvents: false,
        },
      },
    },
  };

  React.useEffect(() => {
    const matches = pathname.match(/(\/sl\/)(.*)/i);
    if (matches && matches[2] && returnedData[matches[2]]) {
      const data = returnedData[matches[2]];
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
  const [hydration, setHydration] = React.useState<HydratedData>(undefined);
  const [tempHydrateStore, setTempHydrateStore] = React.useState<HydrateData>(undefined);
  // const location = useLocation();

  const loadWorkflows = () => {
    registeredWorkflows(availableWorkflows, userConfiguration, [featureFlagFilter]).then(w => {
      setWorkflows(w);
      setIsLoading(false);
    });
  };

  // React.useEffect(() => {
  //   console.log("ROUTE CHANGE", location);
  //   // change temp hydrate data
  // }, [location]);

  const storeHydration = (data: any) => {
    console.log("STORING SOME DATA", data);
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
            <ShortLinkContext.Provider value={{ hydration, tempHydrateStore, storeHydration }}>
              <Routes>
                <Route path="/*" element={<AppLayout isLoading={isLoading} />}>
                  <Route key="landing" path="/" element={<Landing />} />
                  {workflows.map((workflow: Workflow) => (
                    <ErrorBoundary workflow={workflow} key={workflow.path.split("/")[0]}>
                      <Route path={`${workflow.path}/*`} element={<Outlet />}>
                        {workflow.routes.map(route => {
                          const heading = route.displayName
                            ? `${workflow.displayName}: ${route.displayName}`
                            : workflow.displayName;
                          return (
                            <Route
                              key={workflow.path}
                              path={`${route.path}`}
                              element={React.cloneElement(<route.component />, {
                                ...route.componentProps,
                                heading,
                              })}
                            />
                          );
                        })}
                        <Route key={`${workflow.path}/notFound`} path="*" element={<NotFound />} />
                      </Route>
                    </ErrorBoundary>
                  ))}
                  <Route
                    key="short-links"
                    path="/sl/*"
                    element={<ShortLinkHydrator hydrate={setHydration} />}
                  />
                  <Route key="notFound" path="*" element={<NotFound />} />
                </Route>
              </Routes>
            </ShortLinkContext.Provider>
          </ApplicationContext.Provider>
        </div>
      </Theme>
    </Router>
  );
};

export default ClutchApp;
