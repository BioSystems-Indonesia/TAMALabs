export type Device = {
    id: number;
    name: string;
    type: string;
    ip_address: string;
    port: number;
}

export interface DeviceTypeFeatureList {
    id:              string;
    name:            string;
    additional_info: AdditionalInfo;
}

export interface AdditionalInfo {
    can_send:            boolean;
    can_receive:         boolean;
    have_authentication: boolean;
    have_path:           boolean;
}


export const DeviceType = {
    "A15": "A15",
    "BA400": "BA400",
    "BA200": "BA200",
    "Other": "Other"
} as const

export type DeviceTypeValue = typeof DeviceType[keyof typeof DeviceType];