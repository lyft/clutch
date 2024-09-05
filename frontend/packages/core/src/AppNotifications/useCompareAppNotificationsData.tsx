import get from "lodash/get";
import isEmpty from "lodash/isEmpty";
import isEqual from "lodash/isEqual";

import { useUserPreferences } from "../Contexts/preferences-context";
import type { AppBanners } from "../Types";

const useCompareAppNotificationsData = (banners: AppBanners) => {
  const { preferences, dispatch } = useUserPreferences();
  const bannersPreferences: AppBanners = get(preferences, "banners");

  const bannersData = {
    header: {},
    multiWorkflow: {},
    perWorkflow: {},
  };
  let shouldUpdate = false;

  if (!isEmpty(banners?.header)) {
    const headerPreferences = {
      message: bannersPreferences?.header?.message,
      linkText: bannersPreferences?.header.linkText,
      link: bannersPreferences?.header.link,
      severity: bannersPreferences?.header.severity,
    };

    if (!isEqual(banners?.header, headerPreferences)) {
      bannersData.header = { ...banners?.header, dismissed: false };
      shouldUpdate = true;
    } else {
      bannersData.header = { ...bannersPreferences?.header };
    }
  }

  if (!isEmpty(banners?.multiWorkflow)) {
    const multiWorkflowPreferences = {
      title: bannersPreferences?.multiWorkflow?.title,
      message: bannersPreferences?.multiWorkflow?.message,
      workflows: bannersPreferences?.multiWorkflow.workflows,
      link: bannersPreferences?.multiWorkflow.link,
      linkText: bannersPreferences?.multiWorkflow.linkText,
      severity: bannersPreferences?.multiWorkflow.severity,
    };

    if (!isEqual(banners?.multiWorkflow, multiWorkflowPreferences)) {
      bannersData.multiWorkflow = { ...banners?.multiWorkflow, dismissed: false };
      shouldUpdate = true;
    } else {
      bannersData.multiWorkflow = { ...bannersPreferences?.multiWorkflow };
    }
  }

  if (!isEmpty(banners?.perWorkflow)) {
    Object.keys(banners?.perWorkflow).forEach(key => {
      if (bannersPreferences?.perWorkflow?.[key]) {
        const perWorkflowPreferences = {
          title: bannersPreferences?.perWorkflow?.[key]?.title,
          message: bannersPreferences?.perWorkflow?.[key]?.message,
          linkText: bannersPreferences?.perWorkflow?.[key].linkText,
          link: bannersPreferences?.perWorkflow?.[key].link,
          severity: bannersPreferences?.perWorkflow?.[key].severity,
        };

        if (!isEqual(banners?.perWorkflow?.[key], perWorkflowPreferences)) {
          bannersData.perWorkflow[key] = { ...banners?.perWorkflow?.[key], dismissed: false };
          shouldUpdate = true;
        } else {
          bannersData.perWorkflow[key] = bannersPreferences?.perWorkflow?.[key];
        }
      } else {
        bannersData.perWorkflow[key] = { ...banners?.perWorkflow?.[key], dismissed: false };
        shouldUpdate = true;
      }
    });
  }

  return { shouldUpdate, bannersData, dispatch };
};

export default useCompareAppNotificationsData;
