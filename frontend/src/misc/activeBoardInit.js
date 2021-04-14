import columnOrder from '../components/Home/Board/Column/columnOrder';

const activeBoardInit = {
  id: null,
  columns: [
    {
      id: null,
      order: columnOrder.INBOX,
      tasks: [
        {
          id: null,
          order: null,
          title: '',
          description: '',
          subtasks: [
            { id: null, order: null, title: '' },
          ],
        },
      ],
    },
    {
      id: null,
      order: columnOrder.READY,
      tasks: [
        {
          id: null,
          order: null,
          title: '',
          description: '',
          subtasks: [
            { id: null, order: null, title: '' },
          ],
        },
      ],
    },
    {
      id: null,
      order: columnOrder.GO,
      tasks: [
        {
          id: null,
          order: null,
          title: '',
          description: '',
          subtasks: [
            { id: null, order: null, title: '' },
          ],
        },
      ],
    },
    {
      id: null,
      order: columnOrder.DONE,
      tasks: [
        {
          id: null,
          order: null,
          title: '',
          description: '',
          subtasks: [
            { id: null, order: null, title: '' },
          ],
        },
      ],
    },
  ],
};

export default activeBoardInit;
