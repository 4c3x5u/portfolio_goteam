import { createContext } from 'react';

const AppContext = createContext({
  user: {},
  members: [],
  boards: [],
  activeBoard: {},
  loadBoard: () => {},
  isLoading: false,
  setIsLoading: () => {},
});

export default AppContext;
