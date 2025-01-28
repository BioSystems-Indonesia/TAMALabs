export function groupBy<T, K extends keyof T>(array: T[], key: K): Record<T[K] & (string | number | symbol), T[]> {
  return array.reduce((acc, item) => {
    const keyValue = item[key] as unknown as T[K] & (string | number | symbol);
    if (!acc[keyValue]) {
      acc[keyValue] = [];
    }
    acc[keyValue].push(item);
    return acc;
  }, {} as Record<T[K] & (string | number | symbol), T[]>);
}