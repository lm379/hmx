import { describe, it, expect } from 'vitest'
import { formatDate, formatDateTime, formatDuration } from './dateUtils'

describe('formatDuration', () => {
  it('should return empty string for undefined', () => {
    expect(formatDuration(undefined)).toBe('')
  })

  it('should return the duration string as is', () => {
    expect(formatDuration('01:30:45')).toBe('01:30:45')
    expect(formatDuration('00:05:30')).toBe('00:05:30')
  })

  it('should handle empty string', () => {
    expect(formatDuration('')).toBe('')
  })
})

describe('formatDate', () => {
  it('should return empty string for empty input', () => {
    expect(formatDate('')).toBe('')
  })

  it('should return "今天" for today', () => {
    const today = new Date().toISOString()
    expect(formatDate(today)).toBe('今天')
  })

  it('should return "昨天" for yesterday', () => {
    const yesterday = new Date()
    yesterday.setDate(yesterday.getDate() - 1)
    expect(formatDate(yesterday.toISOString())).toBe('昨天')
  })

  it('should return "X天前" for days within a week', () => {
    const threeDaysAgo = new Date()
    threeDaysAgo.setDate(threeDaysAgo.getDate() - 3)
    expect(formatDate(threeDaysAgo.toISOString())).toBe('3天前')
  })

  it('should format date with month-day for current year', () => {
    const date = new Date()
    date.setMonth(date.getMonth() - 1)
    const result = formatDate(date.toISOString())
    expect(result).toMatch(/^\d{1,2}-\d{1,2}$/)
  })

  it('should format date with year-month-day for past years', () => {
    const lastYear = new Date()
    lastYear.setFullYear(lastYear.getFullYear() - 1)
    const result = formatDate(lastYear.toISOString())
    expect(result).toMatch(/^\d{4}-\d{1,2}-\d{1,2}$/)
  })
})

describe('formatDateTime', () => {
  it('should return empty string for empty input', () => {
    expect(formatDateTime('')).toBe('')
  })

  it('should format datetime correctly', () => {
    const dateStr = '2024-01-15T10:30:45'
    const result = formatDateTime(dateStr)
    expect(result).toBe('2024-01-15 10:30:45')
  })

  it('should pad single digits with zeros', () => {
    const dateStr = '2024-01-05T09:05:03'
    const result = formatDateTime(dateStr)
    expect(result).toBe('2024-01-05 09:05:03')
  })

  it('should handle different date formats', () => {
    const dateStr = '2024-12-31T23:59:59'
    const result = formatDateTime(dateStr)
    expect(result).toBe('2024-12-31 23:59:59')
  })
})
