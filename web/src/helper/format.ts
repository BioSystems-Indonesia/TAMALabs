const dateFormatRegex = /^\d{4}-\d{2}-\d{2}$/;
const dateParseRegex = /(\d{4})-(\d{2})-(\d{2})/;

const convertDateToString = (value: string | Date) => {
    // value is a `Date` object
    if (!(value instanceof Date) || isNaN(value.getDate())) return '';
    const pad = '00';
    const yyyy = value.getFullYear().toString();
    const MM = (value.getMonth() + 1).toString();
    const dd = value.getDate().toString();
    return `${yyyy}-${(pad + MM).slice(-2)}-${(pad + dd).slice(-2)}`;
};

export const dateFormatter = (value: string | Date) => {
    // null, undefined and empty string values should not go through dateFormatter
    // otherwise, it returns undefined and will make the input an uncontrolled one.
    if (value == null || value === '') return '';
    if (value instanceof Date) return convertDateToString(value);
    // Valid dates should not be converted
    if (dateFormatRegex.test(value)) return value;

    return convertDateToString(new Date(value));
};

export const dateParser = (value: any) => {
    //value is a string of "YYYY-MM-DD" format
    const match = dateParseRegex.exec(value);
    if (match === null || match.length === 0) return;
    const d = new Date(parseInt(match[1]), parseInt(match[2], 10) - 1, parseInt(match[3]));
    if (isNaN(d.getDate())) return;
    return d;
};


export const requiredAstrix = (required?: boolean) => {
    if (required) {
        return "*";
    }

    return "";
};

export function arrayToRecord<T extends { id: string | number }>(array: T[]): Record<string, T> {
    return array.reduce((record, item) => {
        record[String(item.id)] = item;
        return record;
    }, {} as Record<string, T>);
}
