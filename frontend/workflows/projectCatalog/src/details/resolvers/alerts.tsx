import type { clutch as IClutch } from "@clutch-sh/api";
import { client } from "@clutch-sh/core";
import { faSlack } from "@fortawesome/free-brands-svg-icons";
import { get } from "lodash";

import type { AlertMassageOptions, ProjectAlerts, User } from "../alerts/types";

const massageProjectAlerts = (incidents: any[], options?: AlertMassageOptions): ProjectAlerts => {
  let openCount = 0;
  let triggeredCount = 0;
  let ackCount = 0;
  const onCallUsers: User[] = [];

  incidents.forEach(incident => {
    incident.assignments.forEach(assignment => {
      if (assignment.summary !== "NOOP") {
        onCallUsers.push({
          name: assignment.summary,
          url: assignment.html_url,
        });
      }
    });

    switch (incident.status.toLowerCase()) {
      case "open":
        openCount += 1;
        break;
      case "triggered":
        triggeredCount += 1;
        break;
      case "acknowledged":
        ackCount += 1;
        break;
      default:
        break;
    }
  });
  return {
    title: options?.title ?? "Alerts",
    lastAlert: new Date(get(incidents, ["0", "created_at"])).getTime(),
    summary: {
      open: {
        count: openCount,
      },
      triggered: {
        count: triggeredCount,
      },
      acknowledged: {
        count: ackCount,
      },
    },
    onCall: {
      text: options?.text ?? "Slack to Page Oncall",
      icon: options?.icon ?? faSlack,
      url: options?.url ?? undefined,
      users: onCallUsers,
    },
    create: {
      text: "File an Incident",
      url: get(incidents, ["0", "service", "html_url"]),
    },
  } as ProjectAlerts;
};

const fetchAlerts = async (
  serviceIds: string[],
  options?: AlertMassageOptions
): Promise<ProjectAlerts | any[]> => {
  if (!serviceIds.length) {
    return Promise.resolve([]);
  }

  const { statuses = ["open", "triggered", "acknowledged"], offset = 0 } = options;

  const res = await client.post("/v1/proxy/request", {
    service: "pagerduty",
    httpMethod: "GET",
    path: `incidents?limit=100&offset=${offset}`,
    request: { statuses, service_ids: serviceIds },
  } as IClutch.proxy.v1.IRequestProxyRequest);

  let results = res?.data?.response?.incidents || [];

  if (res?.data?.response?.more === true) {
    const adjustedOffset = offset + res?.data?.response?.incidents?.length;
    results = results.concat(await fetchAlerts(serviceIds, { ...options, offset: adjustedOffset }));
    return results;
  }

  return massageProjectAlerts(results, options);
};

export default fetchAlerts;
