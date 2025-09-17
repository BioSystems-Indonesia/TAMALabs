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
  company_name: z.string(),
  company_address: z.string().optional(),
  company_contact_email: z.string().optional(),
  company_contact_phone: z.string().optional(),
  company_contact_hp: z.string().optional(),
});

export type Settings = z.infer<typeof settingSchema>;

export const defaultSettings: Settings = {
  id: "1", // ID is default to "1" because it is used as the default record ID
  barcode_page_width: 50,
  barcode_page_height: 20,
  barcode_width: 1,
  barcode_height: 50,
  barcode_orientation: "landscape",
  company_name: "RUMAH SAKIT ALINDA HUSADA",
  company_address: "Jl. Raya Tanjung Lesung KM.01, Kec. Panimbang, Kab. Pandeglang, Banten 42281",
  company_contact_email: "rumahsakit.alindahusada@gmail.com",
  company_contact_phone: "02535806781",
  company_contact_hp: "087772714887 | 081399366841"
};
