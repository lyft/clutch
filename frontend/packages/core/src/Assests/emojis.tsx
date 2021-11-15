import * as React from "react";
import styled from "@emotion/styled";

import { IconSizeVariant, STYLE_MAP } from "./global";

const StyledEmoji = styled.svg<{ size?: IconSizeVariant }>(props => ({
  width: `${STYLE_MAP[props.size]?.width || STYLE_MAP.small.width}px`,
  height: `${STYLE_MAP[props.size]?.height || STYLE_MAP.small.height}px`,
}));

export interface EmojiProps {
  size?: IconSizeVariant;
}

const HappyIcon = ({ size }: EmojiProps) => (
  <StyledEmoji size={size} viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
    <circle cx="24" cy="24" r="24" fill="#F59E0B" />
    <circle cx="21.913" cy="21.913" r="21.913" fill="#FBBF24" />
    <path
      d="M19.8261 35.4783C27.3179 35.4783 33.3913 29.8721 33.3913 22.9565H6.26086C6.26086 29.8721 12.3342 35.4783 19.8261 35.4783Z"
      fill="#0D1030"
    />
    <mask
      id="mask0_76_9148"
      style={{ maskType: "alpha" }}
      maskUnits="userSpaceOnUse"
      x="6"
      y="22"
      width="28"
      height="14"
    >
      <path
        d="M19.8261 35.4783C27.3179 35.4783 33.3913 29.8721 33.3913 22.9565H6.26086C6.26086 29.8721 12.3342 35.4783 19.8261 35.4783Z"
        fill="#0D1030"
      />
    </mask>
    <g mask="url(#mask0_76_9148)">
      <circle cx="19.8261" cy="35.4782" r="11.3043" fill="#DB3615" />
    </g>
    <ellipse
      rx="4.17391"
      ry="4.17391"
      transform="matrix(1 0 0 -1 10.4348 16.6957)"
      fill="#0D1030"
    />
    <circle r="4.17391" transform="matrix(1 0 0 -1 10.4348 18.7826)" fill="#FBBF24" />
    <ellipse
      rx="4.17391"
      ry="4.17391"
      transform="matrix(1 0 0 -1 29.2174 16.6957)"
      fill="#0D1030"
    />
    <circle r="4.17391" transform="matrix(1 0 0 -1 29.2174 18.7826)" fill="#FBBF24" />
  </StyledEmoji>
);

const NeutralIcon = ({ size }: EmojiProps) => (
  <StyledEmoji size={size} viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
    <circle cx="24" cy="24" r="24" fill="#F59E0B" />
    <circle cx="24" cy="21.913" r="21.913" fill="#FBBF24" />
    <circle r="4.17391" transform="matrix(1 0 0 -1 14.6087 18.7826)" fill="#0D1030" />
    <ellipse
      rx="4.17391"
      ry="4.17391"
      transform="matrix(1 0 0 -1 33.3913 18.7826)"
      fill="#0D1030"
    />
    <rect x="8.34778" y="27.1304" width="33.3913" height="2.08696" rx="1.04348" fill="#0D1030" />
  </StyledEmoji>
);

const SadIcon = ({ size }: EmojiProps) => (
  <StyledEmoji size={size} viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
    <circle cx="24" cy="24" r="24" fill="#F59E0B" />
    <circle cx="26.087" cy="21.913" r="21.913" fill="#FBBF24" />
    <path
      d="M16.6957 31.3043C16.6957 31.3043 22.0811 25.1727 27.7781 25.0455C33.7131 24.9129 39.6522 31.3043 39.6522 31.3043"
      stroke="#0D1030"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <path
      d="M14.6087 25.2041C14.6087 26.268 13.6744 27.1305 12.5218 27.1305C11.3692 27.1305 10.4348 26.268 10.4348 25.2041C10.4348 24.1401 12.5218 20.8696 12.5218 20.8696C12.5218 20.8696 14.6087 24.1401 14.6087 25.2041Z"
      fill="white"
    />
    <path
      d="M37.5652 12.3718C35.1772 12.3718 33.2413 14.3077 33.2413 16.6957C33.2413 19.0837 35.1772 21.0196 37.5652 21.0196C39.9533 21.0196 41.8891 19.0837 41.8891 16.6957C41.8891 14.3077 39.9533 12.3718 37.5652 12.3718Z"
      fill="#0D1030"
      stroke="white"
      strokeWidth="0.3"
    />
    <circle r="1.04348" transform="matrix(1 0 0 -1 40.6957 17.7392)" fill="white" />
    <circle r="2.08696" transform="matrix(-1 0 0 1 37.5653 16.6958)" fill="white" />
    <path
      d="M18.7826 12.3718C16.3946 12.3718 14.4587 14.3077 14.4587 16.6957C14.4587 19.0837 16.3946 21.0196 18.7826 21.0196C21.1707 21.0196 23.1065 19.0837 23.1065 16.6957C23.1065 14.3077 21.1707 12.3718 18.7826 12.3718Z"
      fill="#0D1030"
      stroke="white"
      strokeWidth="0.3"
    />
    <circle r="2.08696" transform="matrix(-1 0 0 1 18.7827 16.6958)" fill="white" />
    <ellipse rx="1.04348" ry="1.04348" transform="matrix(1 0 0 -1 21.9131 17.7392)" fill="white" />
  </StyledEmoji>
);

export { HappyIcon, NeutralIcon, SadIcon };
