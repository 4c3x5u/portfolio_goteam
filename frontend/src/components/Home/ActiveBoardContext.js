import { createContext } from 'react';

const ActiveBoardContext = createContext({
  activeBoardId: null,
  setActiveBoardId: () => {},
});

export default ActiveBoardContext;
