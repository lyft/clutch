import React, { useEffect } from "react";

import type { AppBanners } from "../Types";

import HeaderNotification from "./HeaderNotification";
import LayoutWithNotifications from "./LayoutWithNotifications";
import compareAppNotificationsData from "./useCompareAppNotificationsData";

interface AppNotificationProps {
  type: "header" | "layout";
  banners: AppBanners;
  workflow?: string;
  children?: React.ReactNode;
}

const AppNotification = ({ type, banners, children, workflow }: AppNotificationProps) => {
  const { shouldUpdate, bannersData, dispatch } = compareAppNotificationsData(banners);

  useEffect(() => {
    if (shouldUpdate) {
      dispatch({
        type: "SetPref",
        payload: {
          key: "banners",
          value: bannersData,
        },
      });
    }
  }, [shouldUpdate]);

  const onDismissAlert = (updatedData: AppBanners) => {
    dispatch({
      type: "SetPref",
      payload: {
        key: "banners",
        value: updatedData,
      },
    });
  };

  return type === "header" ? (
    <HeaderNotification bannersData={bannersData as AppBanners} onDismissAlert={onDismissAlert} />
  ) : (
    <LayoutWithNotifications
      workflow={workflow}
      bannersData={bannersData as AppBanners}
      onDismissAlert={onDismissAlert}
    >
      {children}
    </LayoutWithNotifications>
  );
};

export default AppNotification;
