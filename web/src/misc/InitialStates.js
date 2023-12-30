const InitialStates = {
  user: {
    username: '',
    teamId: null,
    isAdmin: false,
    isAuthenticated: false,
  },

  team: {
    id: null,
    inviteToken: '',
  },

  members: [{
    username: '',
    isActive: false,
    isAdmin: false,
  }],

  boards: [{
    id: null,
    name: '',
  }],

  activeBoard: {
    id: null,
    columns: [
      { id: null, order: 0, tasks: [] },
      { id: null, order: 1, tasks: [] },
      { id: null, order: 2, tasks: [] },
      { id: null, order: 3, tasks: [] },
    ],
  },
};

export default InitialStates;
