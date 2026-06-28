export const TIME_RANGES = ["Last 30m", "2h", "24h", "7d"] as const
export type TimeRange = (typeof TIME_RANGES)[number]
