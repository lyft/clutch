import React from "react";
import { action } from "@storybook/addon-actions";
import type { Meta } from "@storybook/react";

import type { ErrorProps } from "../index";
import Error from "../index";

export default {
  title: "Core/Error",
  component: Error,
} as Meta;

const Template = (props: ErrorProps) => <Error {...props} />;

export const Client = Template.bind({});
Client.args = {
  subject: {
    status: {
      code: 500,
      text: "Client Error",
    },
    message: "Request exceeded timed out of 1 ms",
  },
};

export const Network = Template.bind({});
Network.args = {
  subject: {
    status: {
      code: 404,
      text: "Not Found",
    },
    message: "Not Found",
    data: {},
  },
};

export const Primary = Template.bind({});
Primary.args = {
  subject: {
    code: 5,
    message: "Resource with name `foobar` not found",
    status: {
      code: 404,
      text: "Not Found",
    },
  },
};

export const WithRetry = Template.bind({});
WithRetry.args = {
  ...Primary.args,
  onRetry: action("retry-click"),
};

export const WithLinks = Template.bind({});
WithLinks.args = {
  subject: {
    code: 5,
    message: "Resource with name `foobar` not found",
    status: {
      code: 404,
      text: "Not Found",
    },
    details: [
      {
        _type: "types.googleapis.com/google.rpc.Help",
        links: [
          {
            description: "Please file a ticket here for more help.",
            link: "https://www.example.com",
          },
        ],
      },
    ],
  },
};

export const WithUnknownDetails = Template.bind({});
WithUnknownDetails.args = {
  subject: {
    code: 5,
    message: "Resource with name `foobar` not found",
    status: {
      code: 404,
      text: "Not Found",
    },
    details: [
      {
        _type: "types.googleapis.com/google.rpc.Help",
        links: [
          {
            description: "Please file a ticket here for more help.",
            link: "https://www.example.com",
          },
        ],
      },
      {
        _type: "foobar",
        something: [
          {
            key: "value",
          },
        ],
      },
    ],
  },
};

export const WithWrappedDetails = Template.bind({});
WithWrappedDetails.args = {
  subject: {
    code: 9,
    message: "an error occurred on one or more clusters",
    status: {
      code: 404,
      text: "Not Found",
    },
    details: [
      {
        _type: "types.googleapis.com/google.rpc.Help",
        links: [
          {
            description: "Please file a ticket here for more help.",
            link: "https://www.example.com",
          },
        ],
      },
      {
        _type: "type.googleapis.com/clutch.api.v1.ErrorDetails",
        wrapped: [
          {
            code: 2,
            message: "core-staging-0: yikes",
          },
          {
            code: 16,
            message: "core-staging-1: nono",
            details: [
              {
                type: "type.googleapis.com/clutch.k8s.v1.Status",
                status: "Failure",
                message: "nono",
                reason: "Unauthorized",
                code: 401,
              },
            ],
          },
        ],
      },
    ],
  },
  onRetry: action("retry-click"),
};
