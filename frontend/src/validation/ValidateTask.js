const ValidateTask = {
  title: (title) => {
    if (!title) {
      return 'Task title cannot be empty.';
    } if (title.length > 50) {
      return 'Task title cannot be longer than 50 characters.';
    } return '';
  },
};

export default ValidateTask;
