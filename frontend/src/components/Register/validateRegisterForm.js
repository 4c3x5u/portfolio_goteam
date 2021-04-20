import ValidateUser from '../../validation/ValidateUser';

const validateRegisterForm = (username, password, passwordConfirmation) => {
  let errors = { username: '', password: '', passwordConfirmation: '' };

  const usernameError = ValidateUser.username(username);
  if (usernameError) {
    errors = { ...errors, username: usernameError };
  }

  const passwordError = ValidateUser.password(password);
  if (passwordError) {
    errors = { ...errors, password: passwordError };
  }

  const passwordConfirmationError = ValidateUser.passwordConfirmation(
    passwordConfirmation, password,
  );
  if (passwordConfirmationError) {
    errors = { ...errors, passwordConfirmation: passwordConfirmationError };
  }

  return errors;
};

export default validateRegisterForm;
