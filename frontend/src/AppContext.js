import { createContext } from 'react';

const AppContext = createContext({
  user: {},
  team: {},
  members: [],
  boards: [],
  activeBoard: {},
  loadBoard: () => {},
  isLoading: false,
  setIsLoading: () => {},
  notify: () => {},
});

export default AppContext;
