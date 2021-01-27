import React from 'react';
import clsx from 'clsx';
import { MDXProvider } from '@mdx-js/react';
import Head from '@docusaurus/Head';
import Link from '@docusaurus/Link';
import MDXComponents from '@theme/MDXComponents';
import useBaseUrl from '@docusaurus/useBaseUrl';
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
  const imageUrl = useBaseUrl(image, {
    absolute: true,
  });

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
        <div className="margin-vert--md">
          <time dateTime={date} className={styles.blogPostDate}>
            {month} {day}, {year}{' '}
            {readingTime && <> · {Math.ceil(readingTime)} min read</>}
          </time>
        </div>


        <div className={"avatar margin-vert--md"}>
          {/* Crazy avatar stack code */}
          <div style={{position: 'relative', height: "45px", width: 8 + 45 + ((authors.length - 1) * 20) + "px"}}>
          {authors.map(({ name, avatar }, idx) =>
            <img className={styles.blogPostAvatar} style={{ zIndex: 1000 - idx, marginLeft: idx * 20 + "px" }} src={avatar} alt={name} />
          )}
          </div>

          <div className="avatar__intro">
            <h4 className={clsx(styles.blogPostAuthor, "avatar__name")}>
              {authors.map(({ name, url }, idx) =>
                <>
                  <a href={url} target="_blank" rel="noreferrer noopener">
                    {name}
                  </a>
                  {idx != (authors.length - 1) && <span className={clsx(styles.blogPostAuthorSeparator)}>,&nbsp;</span>}
                </>
              )}
            </h4>
          </div>
        </div>
      </header>
    );
  };

  return (
    <>
      <Head>
        {keywords && keywords.length && (
          <meta name="keywords" content={keywords.join(',')} />
        )}
        {image && <meta property="og:image" content={imageUrl} />}
        {image && <meta property="twitter:image" content={imageUrl} />}
        {image && (
          <meta name="twitter:image:alt" content={`Image for ${title}`} />
        )}
      </Head>

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
