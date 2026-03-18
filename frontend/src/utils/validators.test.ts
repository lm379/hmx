import { describe, it, expect } from 'vitest'
import { validateEmail, validatePhone } from './validators'

describe('validateEmail', () => {
  it('should return true for valid email addresses', () => {
    expect(validateEmail('test@example.com')).toBe(true)
    expect(validateEmail('user123@mail.example.com')).toBe(true)
    expect(validateEmail('user+tag@example.com')).toBe(true)
    expect(validateEmail('first.last@example.com')).toBe(true)
  })

  it('should return false for invalid email addresses', () => {
    expect(validateEmail('invalid')).toBe(false)
    expect(validateEmail('invalid@')).toBe(false)
    expect(validateEmail('@example.com')).toBe(false)
    expect(validateEmail('user@.com')).toBe(false)
    expect(validateEmail('')).toBe(false)
  })

  it('should handle edge cases', () => {
    expect(validateEmail('a@b.cd')).toBe(true)
    expect(validateEmail('user@example')).toBe(false)
    expect(validateEmail('user@example.c')).toBe(false)
  })
})

describe('validatePhone', () => {
  it('should return true for valid Chinese phone numbers', () => {
    expect(validatePhone('13123456789')).toBe(true)
    expect(validatePhone('15123456789')).toBe(true)
    expect(validatePhone('18123456789')).toBe(true)
    expect(validatePhone('19999999999')).toBe(true)
  })

  it('should return false for invalid phone numbers', () => {
    expect(validatePhone('12345678901')).toBe(false) // starts with 2
    expect(validatePhone('10123456789')).toBe(false) // second digit is 0
    expect(validatePhone('12123456789')).toBe(false) // second digit is 2
    expect(validatePhone('1234567890')).toBe(false) // too short
    expect(validatePhone('123456789012')).toBe(false) // too long
    expect(validatePhone('')).toBe(false)
  })

  it('should only accept digits', () => {
    expect(validatePhone('abc12345678')).toBe(false)
    expect(validatePhone('13a2345678')).toBe(false)
  })
})
