
export interface ErrorPayload {
    path : string,
    status_code : number,
    error : string,
    extra_info : Record<string, unknown>
}
