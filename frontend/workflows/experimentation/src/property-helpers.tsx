import type { clutch as IClutch } from "@clutch-sh/api";

const propertyToString = (property: IClutch.chaos.experimentation.v1.IProperty): string => {
  if (property === undefined) {
    return "Unknown";
  }

  if (property.displayValue !== undefined && property.displayValue != null) {
    return property.displayValue.value;
  }
  if (property.intValue !== undefined && property.intValue != null) {
    return property.intValue.toString();
  }
  if (property.dateValue !== undefined && property.dateValue != null) {
    const date = new Date(property.dateValue);
    return date.toLocaleString();
  }
  if (property.stringValue) {
    return property.stringValue;
  }
  return "Unknown";
};

const compareProperties = (
  a: IClutch.chaos.experimentation.v1.IProperty,
  b: IClutch.chaos.experimentation.v1.IProperty
): number => {
  if (a === undefined || b === undefined) {
    return 0;
  }

  if (a.identifier !== b.identifier) {
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
