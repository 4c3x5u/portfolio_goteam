import React from 'react';

const UserContext = React.createContext({
  currentUser: {
    username: '',
    teamId: null,
    isAdmin: false,
    isAuthenticated: false,
  },
  setCurrentUser: () => {},
  boards: [{ id: null, name: '' }],
  setBoards: () => {},
});

export default UserContext;
