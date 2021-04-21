import { createContext } from 'react';

const AppContext = createContext({
  user: {},
  setUser: () => {},
  team: {},
  setTeam: () => {},
  members: [],
  setMembers: () => {},
  boards: [],
  setBoards: () => {},
  activeBoard: {},
  setActiveBoard: () => {},
  loadBoard: () => {},
  isLoading: false,
  setIsLoading: () => {},
  notify: () => {},
});

export default AppContext;
