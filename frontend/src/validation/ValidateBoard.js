import Validate from './Validate';

const ValidateBoard = {
  name: Validate.requiredString('board name', 35),
};

export default ValidateBoard;
