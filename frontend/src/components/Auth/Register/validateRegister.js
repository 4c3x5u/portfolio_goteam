const validateUsername = (username) => {
  let error = '';

  if (!username) {
    error = 'Username cannot be empty.';
  } else if (username.length < 5) {
    error = 'Username cannot be shorter than 5 characters.';
  } else if (username.length > 35) {
    error = 'Username cannot be longer than 35 characters.';
  }

  return error;
};

const validatePassword = (password) => {
  let error = '';

  if (!password) {
    error = 'Password cannot be empty.';
  } else if (password.length < 8) {
    error = 'Password cannot be shorter than 8 characters.';
  } else if (password.length > 255) {
    error = 'Password cannot be longer than 255 characters.';
  }

  return error;
};

const validatePasswordConfirmation = (
  passwordConfirmation, password,
) => {
  let error = '';

  if (!passwordConfirmation) {
    error = 'Confirmation cannot be empty.';
  } else if (passwordConfirmation.length < 8) {
    error = 'Confirmation cannot be shorter than 8 characters.';
  } else if (passwordConfirmation.length > 255) {
    error = 'Confirmation cannot be longer than 255 characters.';
  } else if (password && password !== passwordConfirmation) {
    error = 'Confirmation must match the password.';
  }

  return error;
};

const validateRegister = (username, password, passwordConfirmation) => {
  let errors = {
    username: '',
    password: '',
    passwordConfirmation: '',
  };

  const usernameError = validateUsername(username);
  if (usernameError) {
    errors = { ...errors, username: usernameError };
  }

  const passwordError = validatePassword(password);
  if (passwordError) {
    errors = { ...errors, password: passwordError };
  }

  const passwordConfirmationError = validatePasswordConfirmation(
    passwordConfirmation, password,
  );
  if (passwordConfirmationError) {
    errors = { ...errors, passwordConfirmation: passwordConfirmationError };
  }

  return errors;
};

export default validateRegister;
