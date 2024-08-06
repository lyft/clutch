/**
 * Context: https://stackoverflow.com/questions/65793723/i-cant-find-the-matching-typescript-type-for-displayname-as-a-prop-of-type
 * @param Component the component to fetch the displayName from
 * @returns the displayName of the component
 */
const getDisplayName = (element: React.ReactElement<any>) => {
  const node = element as React.ReactElement<React.ComponentType<any>>;
  const { type } = (node as unknown) as React.ReactElement<React.FunctionComponent>;
  const displayName =
    typeof type === "function"
      ? (type as React.FunctionComponent).displayName ||
        (type as React.FunctionComponent).name ||
        "Unknown"
      : type;
  return displayName;
};

export default getDisplayName;
