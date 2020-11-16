import type { clutch as IClutch } from "@clutch-sh/api";

// Returns a string representation of a given property.
const propertyToString = (property: IClutch.chaos.experimentation.v1.IProperty): string => {
  if (property === undefined) {
    return "Unknown";
  }

  if (property.displayValue !== undefined && property.displayValue != null) {
    return property.displayValue.value;
  }
  if (property.urlValue !== undefined && property.urlValue != null) {
    return property.urlValue.toString();
  }
  if (property.intValue !== undefined && property.intValue != null) {
    return property.intValue.toString();
  }
  if (property.dateValue !== undefined && property.dateValue != null) {
    const date = new Date(property.dateValue as string);
    return date.toLocaleString();
  }
  if (property.stringValue) {
    return property.stringValue;
  }
  return "Unknown";
};

// Compares the values of two properties. It returns:
// * 1 if first property is greater than the second one
// * -1 if the first property is lesser than the second one
// * 0 otherwise
const compareProperties = (
  a: IClutch.chaos.experimentation.v1.IProperty,
  b: IClutch.chaos.experimentation.v1.IProperty
): number => {
  if (a === undefined || b === undefined) {
    return 0;
  }

  if (a.id !== b.id) {
    if (propertyToString(a) > propertyToString(b)) {
      return 1;
    }
    if (propertyToString(a) < propertyToString(b)) {
      return -1;
    }
    return 0;
  }

  const aValue = a.stringValue ?? a.intValue ?? a.dateValue;
  const bValue = b.stringValue ?? b.intValue ?? b.dateValue;

  if (aValue !== undefined && aValue != null) {
    if (aValue > bValue) {
      return 1;
    }
    if (aValue < bValue) {
      return -1;
    }
  }
  return 0;
};

export { propertyToString, compareProperties };
