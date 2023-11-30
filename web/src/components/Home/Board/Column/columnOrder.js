const columnOrder = {
  INBOX: 'inbox',
  READY: 'ready',
  GO: 'go!',
  DONE: 'done',
  parseInt: (order) => {
    switch (order) {
      case 1:
        return columnOrder.INBOX;
      case 2:
        return columnOrder.READY;
      case 3:
        return columnOrder.GO;
      case 4:
        return columnOrder.DONE;
      default:
        return Error('Number must be between 1 and 4');
    }
  },
};

export default columnOrder;
