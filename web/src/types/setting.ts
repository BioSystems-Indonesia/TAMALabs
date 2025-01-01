import { z } from "zod";

export const settingsStoreKey = "settings";

export const settingSchema = z.object({
    id: z.string(),
    barcode_size_width: z.number().default(60),
    barcode_size_height: z.number().default(45),
});

export type Settings = z.infer<typeof settingSchema>;

export const defaultSettings: Settings = {
    id: "1", // ID is default to "1" because it is used as the default record ID
    barcode_size_width: 60,
    barcode_size_height: 45,
};