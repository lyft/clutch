import React from "react";
import { Home as HomeIcon, OpenInNew as OpenInNewIcon } from "@mui/icons-material";
import { Meta, Story } from "@storybook/react";

import QuickLinksCard, { QuickLinksProps } from "../quick-links";

export default {
  title: "Core/QuickLinksCard",
  component: QuickLinksCard,
} as Meta;

const Template: Story<QuickLinksProps> = args => <QuickLinksCard {...args} />;

const imagePath =
  "https://user-images.githubusercontent.com/66325812/164936289-6c855c52-1713-4c21-9157-537bf835307e.svg";

const linkGroups = [
  {
    name: "Group 1",
    imagePath,
    links: [
      {
        name: "Link 1",
        url: "http://example.com/1",
        trackingId: "track1",
      },
      {
        name: "Link 2",
        url: "http://example.com/2",
        trackingId: "track2",
      },
    ],
  },
  {
    name: "Group 2",
    imagePath,
    links: [
      {
        name: "Link 3",
        url: "http://example.com/3",
        trackingId: "track3",
      },
    ],
  },
];

const manyGroups = [
  ...linkGroups,
  {
    name: "Group 3",
    imagePath,
    links: [
      {
        name: "Link 4",
        url: "http://example.com/4",
        trackingId: "track4",
      },
      {
        name: "Link 5",
        url: "http://example.com/5",
        trackingId: "track5",
      },
      {
        name: "Link 6",
        url: "http://example.com/6",
        trackingId: "track6",
      },
    ],
  },
  {
    name: "Group 4",
    imagePath,
    links: [
      {
        name: "Link 7",
        url: "http://example.com/7",
        trackingId: "track7",
      },
      {
        name: "Link 8",
        url: "http://example.com/8",
        trackingId: "track8",
      },
    ],
  },
  {
    name: "Group 5",
    imagePath,
    links: [
      {
        name: "Link 9",
        url: "http://example.com/9",
        trackingId: "track9",
      },
    ],
  },
  {
    name: "Group 6",
    imagePath,
    links: [
      {
        name: "Link 10",
        url: "http://example.com/10",
        trackingId: "track10",
      },
    ],
  },
  {
    name: "Group 7",
    imagePath,
    links: [
      {
        name: "Link 11",
        url: "http://example.com/11",
        trackingId: "track11",
      },
      {
        name: "Link 12",
        url: "http://example.com/12",
        trackingId: "track12",
      },
    ],
  },
];

const withBadges = [
  {
    ...linkGroups[0],
    badge: {
      color: "warning",
      content: "1",
    },
  },
  {
    ...linkGroups[1],
    badge: {
      color: "error",
      content: " ",
    },
  },
];

const withIcons = [
  {
    links: [
      {
        ...linkGroups[0].links[0],
        icon: HomeIcon,
      },
      {
        ...linkGroups[0].links[1],
        icon: OpenInNewIcon,
      },
    ],
    name: "Group 1",
    imagePath,
  },
];

export const Default = Template.bind({});
Default.args = {
  linkGroups,
};

export const TooManyGroups = Template.bind({});
TooManyGroups.args = {
  linkGroups: manyGroups,
};

export const WithBadges = Template.bind({});
WithBadges.args = {
  linkGroups: withBadges,
};

export const WithIcons = Template.bind({});
WithIcons.args = {
  linkGroups: withIcons,
};
