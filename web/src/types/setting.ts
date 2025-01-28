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
});


export type Settings = z.infer<typeof settingSchema>;

export const defaultSettings: Settings = {
    id: "1", // ID is default to "1" because it is used as the default record ID
    barcode_page_width: 50,
    barcode_page_height: 20,
    barcode_width: 1.5,
    barcode_height: 50,
    barcode_orientation: "landscape",
    company_name: "PT Elgatama",
    company_address: "JI. Kyai Caringin No.18A-20, RT.11/RW.4, Cideng, Kecamatan Gambir Daerah Khusus Ibukota Jakarta 10150 - Indonesia",
    company_contact_email: "help@elgatama.com",
    company_contact_phone: "+6281938123234",
};