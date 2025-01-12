import { useStore, type useStoreResult } from "react-admin";
import { Settings, defaultSettings, settingsStoreKey } from "../types/setting";

export default function useSettings(): useStoreResult<Settings> {
  const [settings, setSettings] =  useStore<Settings>(settingsStoreKey, defaultSettings);
  console.debug("store settings", settings)
  const mergeSettings = {
    ...defaultSettings,
    ...settings,
  };
  return [mergeSettings, setSettings];
}
