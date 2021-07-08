import React from 'react';
import clsx from 'clsx';
import { MDXProvider } from '@mdx-js/react';
import Seo from '@theme/Seo';
import Link from '@docusaurus/Link';
import MDXComponents from '@theme/MDXComponents';
import Share from '@site/src/components/Share';
import styles from './styles.module.css';
const MONTHS = [
  'January',
  'February',
  'March',
  'April',
  'May',
  'June',
  'July',
  'August',
  'September',
  'October',
  'November',
  'December',
];

function BlogPostItem(props) {
  const {
    children,
    frontMatter,
    metadata,
    truncated,
    isBlogPostPage = false,
  } = props;
  const { date, permalink, tags, readingTime } = metadata;
  const { authors, title, image, keywords } = frontMatter;

  const renderPostHeader = () => {
    const TitleHeading = isBlogPostPage ? 'h1' : 'h2';
    const match = date.substring(0, 10).split('-');
    const year = match[0];
    const month = MONTHS[parseInt(match[1], 10) - 1];
    const day = parseInt(match[2], 10);
    return (
      <header>
        <TitleHeading
          className={clsx('margin-bottom--sm', styles.blogPostTitle)}>
          {isBlogPostPage ? title : <Link to={permalink}>{title}</Link>}
        </TitleHeading>
        <div className="margin-vert--md" style={{display: "flex", alignItems: "center"}}>
          <time dateTime={date} className={styles.blogPostDate}>
            {month} {day}, {year}{' '}
            {readingTime && <> · {Math.ceil(readingTime)} min read</>}
          </time>
          {!truncated && <>&nbsp;·&nbsp;<Share title={title} authors={authors} style={{margin: "0 7px"}} /></>}
        </div>


        <div className={"avatar margin-vert--md"}>
          {/* Crazy avatar stack code */}
          <div style={{position: 'relative', height: "45px", width: 8 + 45 + ((authors.length - 1) * 20) + "px"}}>
          {authors.map(({ name, avatar }, idx) =>
            <img key={name} className={styles.blogPostAvatar} style={{ zIndex: 1000 - idx, marginLeft: idx * 20 + "px" }} src={avatar} alt={name} />
          )}
          </div>

          <div className="avatar__intro">
            <h4 className={clsx(styles.blogPostAuthor, "avatar__name")}>
              {authors.map(({ name, url }, idx) =>
                <React.Fragment key={name}>
                  <a href={url} target="_blank" rel="noreferrer noopener">
                    {name}
                  </a>
                  {idx != (authors.length - 1) && <span className={clsx(styles.blogPostAuthorSeparator)}>,&nbsp;</span>}
                </React.Fragment>
              )}
            </h4>
          </div>
        </div>
      </header>
    );
  };

  return (
    <>

      <Seo {...{keywords, image}} />

      <article className={clsx(!isBlogPostPage && 'margin-bottom--lg', !isBlogPostPage && styles.blogPostPreview)}>
        {renderPostHeader()}
        <section className="markdown">
          <MDXProvider components={MDXComponents}>{children}</MDXProvider>
        </section>
        {(tags.length > 0 || truncated) && (
          <footer className="row margin-vert--md">
            {tags.length > 0 && (
              <div className="col">
                {tags.map(({ label, permalink: tagPermalink }) => (
                  <Link
                    key={tagPermalink}
                    className={clsx(styles.blogPostTag)}
                    to={tagPermalink}>
                    {label}
                  </Link>
                ))}
              </div>
            )}
            {truncated && (
              <div className="col text--right">
                <Link
                  to={metadata.permalink}
                  aria-label={`Read more about ${title}`}>
                  <strong>Read More</strong>
                </Link>
              </div>
            )}
          </footer>
        )}
      </article>
    </>
  );
}

export default BlogPostItem;
