import Validate from './Validate';

const ValidateUser = {
  username: Validate.requiredString('username', 35, 5),
  password: Validate.requiredString('password', 255, 8),
  passwordConfirmation: (passwordConfirmation, password) => {
    const validationError = (
      Validate.requiredString('confirmation', 255, 8)(passwordConfirmation)
    );

    if (validationError) {
      return validationError;
    } if (password && password !== passwordConfirmation) {
      return 'Confirmation must match the password.';
    } return '';
  },
};

export default ValidateUser;
