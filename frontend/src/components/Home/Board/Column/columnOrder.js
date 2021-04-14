const columnOrder = {
  INBOX: 'INBOX',
  READY: 'READY',
  GO: 'GO',
  DONE: 'DONE',
};

export const orderToInt = (order) => {
  switch (order) {
    case columnOrder.INBOX: return 0;
    case columnOrder.READY: return 1;
    case columnOrder.GO: return 2;
    case columnOrder.DONE: return 3;
    default: return Error('Number must be between 0 and 4');
  }
};

export default columnOrder;
