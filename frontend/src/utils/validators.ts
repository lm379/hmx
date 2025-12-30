// Regular expressions for validation
export const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
export const phoneRegex = /^1[3-9]\d{9}$/;

/**
 * Validates an email address.
 * @param email The email string to validate.
 * @returns true if valid, false otherwise.
 */
export const validateEmail = (email: string): boolean => {
  return emailRegex.test(email);
};

/**
 * Validates a phone number (China Mainland).
 * @param phone The phone string to validate.
 * @returns true if valid, false otherwise.
 */
export const validatePhone = (phone: string): boolean => {
  return phoneRegex.test(phone);
};
