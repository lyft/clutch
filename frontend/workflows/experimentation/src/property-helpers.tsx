import type { clutch as IClutch } from "@clutch-sh/api";

// Returns a string representation of a given property.
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
    const date = new Date(property.dateValue as string);
    return date.toLocaleString();
  }
  if (property.stringValue) {
    return property.stringValue;
  }
  return "Unknown";
};

export default propertyToString;
