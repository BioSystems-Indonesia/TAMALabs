import { z } from "zod";

export const settingsStoreKey = "settings";

export const orientationChoices = ["portrait", "landscape"] as const;

export const settingSchema = z.object({
    id: z.string(),
    barcode_page_width: z.number(),
    barcode_page_height: z.number(),
    barcode_height: z.number(),
    barcode_width: z.number(),
    barcode_orientation: z.enum(["portrait", "landscape"]),
});


export type Settings = z.infer<typeof settingSchema>;

export const defaultSettings: Settings = {
    id: "1", // ID is default to "1" because it is used as the default record ID
    barcode_page_width: 60,
    barcode_page_height: 45,
    barcode_width: 1.5,
    barcode_height: 30,
    barcode_orientation: "portrait",
};