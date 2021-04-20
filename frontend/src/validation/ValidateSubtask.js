import Validate from './Validate';

const ValidateSubtask = {
  title: Validate.requiredString('subtask title', 50),
};

export default ValidateSubtask;
