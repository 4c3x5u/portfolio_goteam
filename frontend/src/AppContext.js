import { createContext } from 'react';

const AppContext = createContext({
  user: {},
  boards: [],
  activeBoard: {},
  loadBoard: () => {},
  isLoading: false,
  setIsLoading: () => {},
});

export default AppContext;
