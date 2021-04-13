import { createContext } from 'react';

const AppContext = createContext({
  currentUser: {
    username: '',
    teamId: null,
    isAdmin: false,
    isAuthenticated: false,
  },
  setCurrentUser: () => {},
  boards: [{ id: null, name: '' }],
  setBoards: () => {},
  activeBoardId: null,
  setActiveBoardId: () => {},
});

export default AppContext;
