export const User = {
  validateUsername: (username) => {
    if (!username) {
      return 'Username cannot be empty.';
    } if (username.length < 5) {
      return 'Username cannot be shorter than 5 characters.';
    } if (username.length > 35) {
      return 'Username cannot be longer than 35 characters.';
    } return '';
  },
  validatePassword: (password) => {
    if (!password) {
      return 'Password cannot be empty.';
    } if (password.length < 8) {
      return 'Password cannot be shorter than 8 characters.';
    } if (password.length > 255) {
      return 'Password cannot be longer than 255 characters.';
    } return '';
  },
  validatePasswordConfirmation: (passwordConfirmation, password) => {
    if (!passwordConfirmation) {
      return 'Confirmation cannot be empty.';
    } if (passwordConfirmation.length < 8) {
      return 'Confirmation cannot be shorter than 8 characters.';
    } if (passwordConfirmation.length > 255) {
      return 'Confirmation cannot be longer than 255 characters.';
    } if (password && password !== passwordConfirmation) {
      return 'Confirmation must match the password.';
    } return '';
  },
};

export default { User };
