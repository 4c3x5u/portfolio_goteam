import { createContext } from 'react';
import activeBoardInit from './misc/activeBoardInit';

const AppContext = createContext({
  user: {
    username: '',
    teamId: null,
    isAdmin: false,
    isAuthenticated: false,
  },
  setUser: () => {},
  activeBoard: activeBoardInit,
  setActiveBoard: () => {},
  loadBoard: () => {},
  isLoading: false,
  setIsLoading: () => {},
});

export default AppContext;
