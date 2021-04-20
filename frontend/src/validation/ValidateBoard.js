const ValidateBoard = {
  name: (name) => {
    if (!name) {
      return 'Board name cannot be empty.';
    } if (name.length > 35) {
      return 'Board name cannot be longer than 35 characters.';
    } return '';
  },
};

export default ValidateBoard;
