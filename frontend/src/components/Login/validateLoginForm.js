import ValidateUser from '../../validation/ValidateUser';

const validateLoginForm = (username, password) => {
  let errors = { username: '', password: '' };

  const usernameError = ValidateUser.username(username);
  if (usernameError) {
    errors = { ...errors, username: usernameError };
  }

  const passwordError = ValidateUser.password(password);
  if (passwordError) {
    errors = { ...errors, password: passwordError };
  }

  return errors;
};

export default validateLoginForm;
