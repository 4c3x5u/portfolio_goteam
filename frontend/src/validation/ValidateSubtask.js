const ValidateSubtask = {
  title: (title) => {
    if (!title) {
      return 'Subtask title cannot be empty.';
    } if (title.length > 50) {
      return 'Subtask title cannot be longer than 50 characters.';
    } return '';
  },
};

export default ValidateSubtask;
