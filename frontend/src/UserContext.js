import React from 'react';

const UserContext = React.createContext({
  currentUser: {
    username: '',
    teamId: null,
    isAdmin: false,
    isAuthenticated: false,
  },
  setCurrentUser: () => {},
});

export default UserContext;
