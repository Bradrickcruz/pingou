import { createContext, useContext } from "react";
import { DEFAULT_PREFERENCE } from "../lib/theme";

export const ThemeContext = createContext({
  preference: DEFAULT_PREFERENCE,
  resolved: "light",
  setPreference: () => {},
});

export function useTheme() {
  return useContext(ThemeContext);
}
