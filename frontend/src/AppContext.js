import { createContext } from 'react';
import activeBoardInit from './misc/activeBoardInit';

const AppContext = createContext({
  user: {},
  boards: [],
  activeBoard: activeBoardInit,
  loadBoard: () => {},
  isLoading: false,
});

export default AppContext;
