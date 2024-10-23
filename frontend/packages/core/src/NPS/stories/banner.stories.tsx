import React from "react";
import { MemoryRouter } from "react-router-dom";
import type { Meta } from "@storybook/react";

import { Header } from "../../AppLayout/header";
import { HappyEmoji, NeutralEmoji, SadEmoji } from "../../Assets/emojis";
import { ApplicationContext, UserPreferencesProvider } from "../../Contexts";
import type { HeaderItem, TriggeredHeaderData } from "../../Contexts/app-context";
import type { FeedbackBannerProps } from "../banner";
import { Banner } from "../banner";

export default {
  title: "Core/NPS/Banner",
  component: Banner,
  argTypes: {
    icon: {
      options: ["None", "Happy", "Neutral", "Sad"],
      mapping: { None: null, Happy: <HappyEmoji />, Neutral: <NeutralEmoji />, Sad: <SadEmoji /> },
    },
  },
} as Meta;

const Template = ({ ...props }: FeedbackBannerProps) => {
  const [triggeredHeaderData, setTriggeredHeaderData] = React.useState<TriggeredHeaderData>();
  return (
    <ApplicationContext.Provider
      // eslint-disable-next-line react/jsx-no-constructed-context-values
      value={{
        workflows: [
          {
            developer: { name: "Lyft", contactUrl: "mailto:hello@clutch.sh" },
            displayName: "EC2",
            group: "AWS",
            path: "ec2",
            routes: [
              {
                // eslint-disable-next-line react/no-unstable-nested-components
                component: () => <div>Terminate Instance</div>,
                componentProps: { resolverType: "clutch.aws.ec2.v1.Instance" },
                description: "Terminate an EC2 instance.",
                displayName: "Terminate Instance",
                path: "instance/terminate",
                requiredConfigProps: ["resolverType"],
                trending: true,
              },
              {
                // eslint-disable-next-line react/no-unstable-nested-components
                component: () => <div>Resize ASG</div>,
                componentProps: { resolverType: "clutch.aws.ec2.v1.AutoscalingGroup" },
                description: "Resize an autoscaling group.",
                displayName: "Resize Autoscaling Group",
                path: "asg/resize",
                requiredConfigProps: ["resolverType"],
              },
            ],
          },
        ],
        triggerHeaderItem: (item: HeaderItem, data: unknown) =>
          setTriggeredHeaderData({
            ...triggeredHeaderData,
            [item]: {
              ...(data as any),
            },
          }),
        triggeredHeaderData,
      }}
    >
      <UserPreferencesProvider>
        <MemoryRouter>
          <Header enableNPS />
          <Banner {...props} />
        </MemoryRouter>
      </UserPreferencesProvider>
    </ApplicationContext.Provider>
  );
};

export const Primary = Template.bind({});
