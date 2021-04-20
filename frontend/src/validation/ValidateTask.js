import Validate from './Validate';

const ValidateTask = {
  title: Validate.requiredString('task title', 50),
};

export default ValidateTask;
