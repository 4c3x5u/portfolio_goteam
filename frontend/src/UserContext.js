import { createContext } from 'react';

const UserContext = createContext({
  user: {
    username: '',
    teamId: null,
    isAdmin: false,
    isAuthenticated: false,
  },
  setUser: () => {},
});

export default UserContext;
