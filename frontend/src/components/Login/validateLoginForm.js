import { User } from '../../misc/validators';

const validateLoginForm = (username, password) => {
  let errors = { username: '', password: '' };

  const usernameError = User.validateUsername(username);
  if (usernameError) {
    errors = { ...errors, username: usernameError };
  }

  const passwordError = User.validatePassword(password);
  if (passwordError) {
    errors = { ...errors, password: passwordError };
  }

  return errors;
};

export default validateLoginForm;
