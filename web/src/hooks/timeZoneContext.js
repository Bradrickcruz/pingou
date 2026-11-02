import { createContext, useContext } from "react";
import { DEFAULT_TIME_ZONE } from "../lib/timezone";

export const TimeZoneContext = createContext({
  timeZone: DEFAULT_TIME_ZONE,
  setTimeZone: () => {},
});

export function useTimeZone() {
  return useContext(TimeZoneContext);
}
