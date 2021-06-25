import React from "react";
import BrowserOnly from "@docusaurus/BrowserOnly";
import { ArrowContainer, Popover } from "react-tiny-popover"

import "./styles.css";

const trackShare = (source) => {
  // @ts-ignore
  if (window.ga) {
    // @ts-ignore
    window.ga('send', {
      hitType: 'event',
      eventCategory: 'Blog',
      eventAction: `share via ${source}`,
      eventLabel: window.location.pathname,
    });
  }
};

interface ShareLinkProps {
  className: string;
  label: string;
  color: string;
  link?: string;
  onClick?: () => void;
}

const ShareLink = ({ className, label, color, link, onClick }: ShareLinkProps) => (
  <a
    className="share-link"
    href={link}
    target="_blank"
    rel="noopener noreferrer"
    onClick={onClick}
  >
    <span className={`fe ${className}`} style={{ margin: "7px 10px", color }} />
    {label}
  </a>
);

interface ShareProps {
  title: string;
  authors: {
    name: string;
    url: string;
    twitter_username?: string;
  }[];
}

const Share = ({title, authors}: ShareProps) => {
  const [open, setOpen] = React.useState(false);

  return (
      <BrowserOnly>
        {() => {
          const twitter_authors = authors.map(a => a.twitter_username ? `@${a.twitter_username}` : a.name);
          const tweet = encodeURI(
            `https://twitter.com/intent/tweet?text=${title} by ${twitter_authors.join(", ")} ${window.location.href}`
          );
          return (
            <Popover
              isOpen={open}
              positions={["bottom"]}
              padding={10}
              align={"center"}
              onClickOutside={() => setOpen(false)}
              content={({ position, childRect, popoverRect }) => (
                <ArrowContainer
                  position={position}
                  childRect={childRect}
                  popoverRect={popoverRect}
                  arrowColor="var(--ifm-color-content-secondary)"
                  arrowSize={10}
                  arrowStyle={{ opacity: ".1" }}
                  className="popover-arrow-container"
                  arrowClassName="popover-arrow"
                >
                  <div className="share-popover-container">
                    <div className="share-link-list">
                      <ShareLink
                        className="fe-twitter"
                        label="Twitter" color="#1DA1F2"
                        link={tweet}
                        onClick={() => {
                          setOpen(false);
                          trackShare("twitter");
                        }}
                      />
                      <ShareLink
                        className="fe-linkedin"
                        label="LinkedIn"
                        color="#0072b1"
                        link={`https://www.linkedin.com/sharing/share-offsite/?url=${window.location.href}`}
                        onClick={() => {
                          setOpen(false);
                          trackShare("linkedin");
                        }}
                      />
                      <ShareLink
                        className="fe-link"
                        label="Copy Link"
                        color="var(--ifm-color-content-secondary)"
                        onClick={() => {
                          var tmp = document.createElement("input"),
                          href = window.location.href;
                          document.body.appendChild(tmp);
                          tmp.value = href;
                          tmp.select();
                          document.execCommand("copy");
                          document.body.removeChild(tmp);
                          setOpen(false);
                        }}
                      />
                    </div>
                  </div>
                </ArrowContainer>
              )}
            >
              <span className={"fe fe-share"} onClick={() => setOpen(o => !o)}/>
            </Popover>
          )
        }}
      </BrowserOnly>
  );
};

export default Share;