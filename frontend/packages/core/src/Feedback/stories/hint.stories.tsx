import React from "react";
import type { Meta } from "@storybook/react";

import Hint from "../hint";

export default {
  title: "Core/Feedback/Hint",
  component: Hint,
} as Meta;

const Template = (props: { url: string; text: string }) => {
  const { text, url } = props;
  return (
    <Hint>
      {url && <img alt="demo url" src={url} />}
      <div style={{ padding: "10px" }}>{text}</div>
    </Hint>
  );
};

export const Primary = Template.bind({});
Primary.args = {
  text: "Some helpful text",
};

export const EmbeddedImage = Template.bind({});
EmbeddedImage.args = {
  url: "https://www.clutch.sh/img/favicon.ico",
  text: "This is the Clutch logo!",
};
