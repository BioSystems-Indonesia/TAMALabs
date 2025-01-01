import { useStore, type useStoreResult } from "react-admin";
import { Settings, defaultSettings, settingsStoreKey } from "../types/setting";

export default function useSettings(): useStoreResult<Settings> {
  return useStore<Settings>(settingsStoreKey, defaultSettings);
}
