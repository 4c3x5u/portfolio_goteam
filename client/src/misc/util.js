export const getAuthHeaders = () => ({
  headers: {
    'auth-user': sessionStorage.getItem('username'),
    'auth-token': sessionStorage.getItem('auth-token'),
  },
});

export const capFirstLetterOf = (text) => (
  text.charAt(0).toUpperCase() + text.slice(1)
);
