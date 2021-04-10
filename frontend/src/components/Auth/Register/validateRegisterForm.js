import { User } from '../../../misc/validators';

const validateRegisterForm = (username, password, passwordConfirmation) => {
  let errors = {
    username: '',
    password: '',
    passwordConfirmation: '',
  };

  const usernameError = User.validateUsername(username);
  if (usernameError) {
    errors = { ...errors, validateUsername: usernameError };
  }

  const passwordError = User.validatePassword(password);
  if (passwordError) {
    errors = { ...errors, password: passwordError };
  }

  const passwordConfirmationError = User.validatePasswordConfirmation(
    passwordConfirmation, password,
  );
  if (passwordConfirmationError) {
    errors = { ...errors, passwordConfirmation: passwordConfirmationError };
  }

  return errors;
};

export default validateRegisterForm;
